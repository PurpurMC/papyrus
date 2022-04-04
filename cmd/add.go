package cmd

import (
	"context"
	"fmt"
	"github.com/purpurmc/papyrus/data"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"path/filepath"
	"time"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add something new to the database",
	Long:  `Ths command will add a new project, version or build to the database.`,
}

var addProjectCommand = &cobra.Command{
	Use:   "project [project]",
	Short: "Add a new project to the database",
	Long:  `This command will add a new project to the database.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		database, _ := db.NewMongo()
		defer database.Client().Disconnect(context.TODO())
		projectName := args[0]

		fmt.Println("Adding project", projectName)

		var createdAt int64
		if cmd.Flags().Changed("createdAt") {
			createdAt, _ = cmd.Flags().GetInt64("createdAt")
		} else {
			createdAt = time.Now().Unix()
		}

		db.InsertProject(database, types.Project{
			Name:      projectName,
			CreatedAt: createdAt,
		})

		fmt.Println("Project added!")
	},
}

var addVersionCommand = &cobra.Command{
	Use:   "version [project] [version]",
	Short: "Add a new version to the database",
	Long:  `This command will add a new version to the database.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		database, _ := db.NewMongo()
		defer database.Client().Disconnect(context.TODO())
		projectName := args[0]
		versionName := args[1]

		fmt.Println("Adding version", versionName, "to project", projectName)

		var createdAt int64
		if cmd.Flags().Changed("createdAt") {
			createdAt, _ = cmd.Flags().GetInt64("createdAt")
		} else {
			createdAt = time.Now().Unix()
		}

		project := db.GetProject(database, &types.Project{Name: projectName})
		var projectId primitive.ObjectID
		if project == nil {
			projectId = db.InsertProject(database, types.Project{
				Name:      projectName,
				CreatedAt: time.Now().Unix(),
			})
		} else {
			projectId = project.Id
		}

		db.InsertVersion(database, types.Version{
			ProjectId: projectId,
			CreatedAt: createdAt,
			Name:      versionName,
		})

		fmt.Println("Version added!")
	},
}

var addBuildCommand = &cobra.Command{
	Use:   "build [project] [version] [build]",
	Short: "Add a new build to the database",
	Long:  `This command will add a new build to the database.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		database, bucket := db.NewMongo()
		defer database.Client().Disconnect(context.TODO())
		projectName := args[0]
		versionName := args[1]
		buildName := args[2]

		dataSource, _ := cmd.Flags().GetString("data")
		flags, _ := cmd.Flags().GetStringSlice("flags")

		fmt.Println("Adding build", buildName, "to version", versionName, "of project", projectName)

		project := db.GetProject(database, &types.Project{Name: projectName})
		var projectId primitive.ObjectID
		if project == nil {
			projectId = db.InsertProject(database, types.Project{
				Name:      projectName,
				CreatedAt: time.Now().Unix(),
			})
		} else {
			projectId = project.Id
		}

		version := db.GetVersion(database, &types.Version{ProjectId: projectId, Name: versionName})
		var versionId primitive.ObjectID
		if version == nil {
			versionId = db.InsertVersion(database, types.Version{
				ProjectId: projectId,
				CreatedAt: time.Now().Unix(),
				Name:      versionName,
			})
		} else {
			versionId = version.Id
		}

		switch dataSource {
		case "jenkins":
			if !cmd.Flags().Changed("jenkinsUrl") || !cmd.Flags().Changed("jenkinsJob") || !cmd.Flags().Changed("jenkinsBuild") {
				fmt.Println("You must specify a Jenkins URL, job and build")
				return
			}

			jenkinsUrl, _ := cmd.Flags().GetString("jenkinsUrl")
			jenkinsJob, _ := cmd.Flags().GetString("jenkinsJob")
			jenkinsBuild, _ := cmd.Flags().GetString("jenkinsBuild")
			jenkinsWorkspaceFiles, _ := cmd.Flags().GetStringArray("jenkinsWorkspace")

			fmt.Println("Getting build info from Jenkins")

			jenkinsData := data.GetJenkinsData(jenkinsUrl, jenkinsJob, jenkinsBuild)
			if jenkinsData == nil {
				fmt.Println("Failed to get Jenkins data")
				return
			}

			fmt.Println("Getting workspace files from Jenkins")

			var files []types.File
			if len(jenkinsWorkspaceFiles) > 0 {
				files = make([]types.File, len(jenkinsWorkspaceFiles)-1)

				for _, file := range jenkinsWorkspaceFiles {
					data := data.DownloadJenkinsWorkspaceFile(jenkinsUrl, jenkinsJob, file)
					if data == nil {
						fmt.Println("Failed to download workspace file:", file)
						fmt.Println("Skipping...")
						continue
					}

					fileId, fileName, hash, contentType := db.UploadFile(bucket, data)
					files = append(files, types.File{
						Id:           fileId,
						InternalName: fileName,
						ContentType:  contentType,
						Name:         filepath.Base(file),
						SHA512:       hash,
					})
				}
			} else {
				files = make([]types.File, 0)
			}

			fmt.Println("Creating build from Jenkins data")

			var commits []types.Commit
			if len(jenkinsData.ChangeSet.Items) > 0 {
				commits = make([]types.Commit, len(jenkinsData.ChangeSet.Items)-1)
				for _, item := range jenkinsData.ChangeSet.Items {
					commits = append(commits, types.Commit{
						Author:      item.Author.Name,
						Email:       item.AuthorEmail,
						Summary:     item.Summary,
						Description: item.Description,
						Hash:        item.Hash,
						Timestamp:   item.Timestamp,
					})
				}
			} else {
				commits = make([]types.Commit, 0)
			}

			db.InsertBuild(database, types.Build{
				VersionId: versionId,
				CreatedAt: jenkinsData.Timestamp,
				Name:      buildName,
				Result:    jenkinsData.Result,
				Flags:     flags,
				Commits:   commits,
				Files:     files,
			})
		default:
			fmt.Println("Invalid data source, currently only jenkins is supported")
			return
		}

		fmt.Println("Build added!")
	},
}

func init() {
	rootCommand.AddCommand(addCommand)
	addCommand.DisableFlagsInUseLine = true

	addCommand.AddCommand(addProjectCommand)
	addProjectCommand.Flags().Int64("createdAt", 0, "The date the version was created")

	addCommand.AddCommand(addVersionCommand)
	addVersionCommand.Flags().Int64("createdAt", 0, "The date the version was created")

	addCommand.AddCommand(addBuildCommand)
	addBuildCommand.Flags().String("data", "jenkins", "The way to get data for the build, currently only jenkins is supported")
	addBuildCommand.Flags().StringArray("flags", make([]string, 0), "The flags for the build")

	addBuildCommand.Flags().String("jenkinsUrl", "", "The url to the jenkins server")
	addBuildCommand.Flags().String("jenkinsJob", "", "The job on the jenkins server")
	addBuildCommand.Flags().String("jenkinsBuild", "", "The build on the jenkins server")
	addBuildCommand.Flags().StringArray("jenkinsWorkspace", make([]string, 0), "The files to get from the current jenkins workspace")
	// addBuildCommand.Flags().StringArray("jenkinsArtifacts", make([]string, 0), "The files to get from the builds jenkins artifacts. By default all artifacts are downloaded") // todo: implement this

	/* todo: manual data source
	addBuildCommand.Flags().Int64("createdAt", 0, "The date the version was created")
	addBuildCommand.Flags().String("result", "", "The result of the build")
	addBuildCommand.Flags().StringSlice("commits", make([]string, 0), "The commits for the build")
	addBuildCommand.Flags().StringSlice("files", make([]string, 0), "The files for the build")
	*/
}
