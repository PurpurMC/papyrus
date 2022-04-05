package cmd

import (
	"context"
	"fmt"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete something from the database",
	Long:  `Ths command will delete a project, version or build from the database.`,
}

var deleteProjectCommand = &cobra.Command{
	Use:   "project [project]",
	Short: "Delete a project from the database",
	Long:  `Ths command will delete a project from the database.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		database, bucket := db.NewMongo()
		defer database.Client().Disconnect(context.TODO())
		projectName := args[0]

		fmt.Println("Are you sure you want to delete the project " + projectName + "? (y/n)")
		if !confirmation() {
			return
		}

		fmt.Println("Deleting project " + projectName + "...")
		project := db.GetProject(database, &types.Project{Name: projectName})
		if project == nil {
			fmt.Println("Project not found!")
			return
		}

		versions := db.GetVersions(database, &types.Version{ProjectId: project.Id})
		for _, version := range versions {
			builds := db.GetBuildsFromVersion(database, version.Id)
			for _, build := range builds {
				deleteFiles(database, bucket, &build, false, true)
				db.DeleteBuilds(database, &types.Build{Id: build.Id})
			}

			db.DeleteVersions(database, &types.Version{Id: version.Id})
		}

		db.DeleteProjects(database, &types.Project{Name: projectName})
		fmt.Println("Project deleted!")
	},
}

var deleteVersionCommand = &cobra.Command{
	Use:   "version [project] [version]",
	Short: "Delete a version from the database",
	Long:  `Ths command will delete a version from the database.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		database, bucket := db.NewMongo()
		defer database.Client().Disconnect(context.TODO())
		projectName := args[0]
		versionName := args[1]

		fmt.Println("Are you sure you want to delete the version " + versionName + " of " + projectName + "? (y/n)")
		if !confirmation() {
			return
		}

		project := db.GetProject(database, &types.Project{Name: projectName})
		if project == nil {
			fmt.Println("Project not found!")
			return
		}

		version := db.GetVersion(database, &types.Version{ProjectId: project.Id, Name: versionName})
		if version == nil {
			fmt.Println("Version not found!")
			return
		}

		builds := db.GetBuildsFromVersion(database, version.Id)
		for _, build := range builds {
			deleteFiles(database, bucket, &build, true, false)
			db.DeleteBuilds(database, &types.Build{Id: build.Id})
		}

		db.DeleteVersions(database, &types.Version{ProjectId: project.Id, Name: versionName})
		fmt.Println("Version deleted!")
	},
}

var deleteBuildCommand = &cobra.Command{
	Use:   "build [project] [version] [build]",
	Short: "Delete a build from the database",
	Long:  `Ths command will delete a build from the database.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		database, bucket := db.NewMongo()
		defer database.Client().Disconnect(context.TODO())
		projectName := args[0]
		versionName := args[1]
		buildName := args[2]

		fmt.Println("Are you sure you want to delete the build " + buildName + "? (y/n)")
		if !confirmation() {
			return
		}

		project := db.GetProject(database, &types.Project{Name: projectName})
		if project == nil {
			fmt.Println("Project not found!")
			return
		}

		version := db.GetVersion(database, &types.Version{ProjectId: project.Id, Name: versionName})
		if version == nil {
			fmt.Println("Version not found!")
			return
		}

		var build *types.Build
		builds := db.GetBuildsFromVersion(database, version.Id)
		for _, b := range builds {
			if b.Name == buildName {
				build = &b
				break
			}
		}

		if build == nil {
			fmt.Println("Build not found!")
			return
		}

		deleteFiles(database, bucket, build, false, false)
		db.DeleteBuilds(database, &types.Build{Id: build.Id, Name: buildName})
		fmt.Println("Build deleted!")
	},
}

func init() {
	rootCommand.AddCommand(deleteCommand)
	deleteCommand.DisableFlagsInUseLine = true

	deleteCommand.AddCommand(deleteProjectCommand)
	deleteCommand.AddCommand(deleteVersionCommand)
	deleteCommand.AddCommand(deleteBuildCommand)
}

func confirmation() bool {
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		fmt.Println("Aborting...")
		return false
	}
	return true
}

func deleteFiles(database *mongo.Database, bucket *gridfs.Bucket, build *types.Build, checkVersions bool, checkProjects bool) {
	builds := db.GetBuilds(database, nil)

	for _, file := range build.Files {
		existsInOtherBuilds := false
		for _, otherBuild := range builds {
			if otherBuild.Id == build.Id || (checkVersions && otherBuild.VersionIds[0] != build.VersionIds[0]) {
				continue
			}

			if checkProjects {
				version := db.GetVersion(database, &types.Version{Id: build.VersionIds[0]})
				project := db.GetProject(database, &types.Project{Id: version.ProjectId})

				otherVersion := db.GetVersion(database, &types.Version{Id: otherBuild.VersionIds[0]})
				otherProject := db.GetProject(database, &types.Project{Id: otherVersion.ProjectId})

				if project.Id == otherProject.Id {
					continue
				}
			}

			for _, otherFile := range otherBuild.Files {
				if otherFile.Name == file.Name {
					existsInOtherBuilds = true
					break
				}
			}
		}

		if !existsInOtherBuilds {
			db.DeleteFile(bucket, file.Id)
		}
	}
}
