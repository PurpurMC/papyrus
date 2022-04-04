package cmd

import (
	"context"
	"fmt"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"github.com/purpurmc/papyrus/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"os"
	"time"
)

var updateCommand = &cobra.Command{
	Use:   "update",
	Short: "After updating papyrus, this command will update the config and the db",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetDefault("_version", 1)

		viper.SetDefault("http.host", "127.0.0.1")
		viper.SetDefault("http.port", 8080)
		viper.SetDefault("http.debug", false)

		viper.SetDefault("http.routes.prefix", "/v1")
		viper.SetDefault("http.routes.list-projects", "/")
		viper.SetDefault("http.routes.get-project", "/:project")
		viper.SetDefault("http.routes.get-version", "/:project/:version")
		viper.SetDefault("http.routes.get-build", "/:project/:version/:build")
		viper.SetDefault("http.routes.download-build", "/:project/:version/:build/:file")

		viper.SetDefault("db.host", "127.0.0.1")
		viper.SetDefault("db.port", 27017)
		viper.SetDefault("db.username", "username")
		viper.SetDefault("db.password", "password")
		viper.SetDefault("db.db", "papyrus")

		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error writing config file, please make sure you have the correct permissions")
			fmt.Println(err)
			os.Exit(1)
		}

		database, bucket := db.NewMongo()
		defer database.Client().Disconnect(context.TODO())

		collections, err := database.ListCollectionNames(context.TODO(), bson.D{{}})
		if err != nil {
			fmt.Println("Error listing collections")
			fmt.Println(err)
			os.Exit(1)
		}

		if !utils.StringInSlice("projects", collections) {
			db.CreateCollection(database, "projects")
		}

		if !utils.StringInSlice("versions", collections) {
			db.CreateCollection(database, "versions")
		}

		if !utils.StringInSlice("builds", collections) {
			db.CreateCollection(database, "builds")
		}

		projectId := db.InsertProject(database, types.Project{
			Name:      "purpur",
			CreatedAt: time.Now().Unix(),
		})

		versionId := db.InsertVersion(database, types.Version{
			ProjectId: projectId,
			CreatedAt: time.Now().Unix(),
			Name:      "1.18.2",
		})

		data, err := ioutil.ReadFile("/home/ben/downloads/purpur-1.18.2-1583.jar")
		if err != nil {
			fmt.Println("Error reading file")
			fmt.Println(err)
			os.Exit(1)
		}

		fileName, hash, contentType := db.UploadFile(bucket, data)

		db.InsertBuild(database, types.Build{
			VersionId: versionId,
			CreatedAt: time.Now().Unix(),
			Name:      "1583",
			Result:    "SUCCESS",
			Commits: []types.Commit{
				{
					Author:      "BillyGalbreath",
					Email:       "blake.galbreath@gmail.com",
					Summary:     "clean up level/getWorld usage in various patches",
					Description: "clean up level/getWorld usage in various patches\n",
					Hash:        "e9af29e2c6b61a370a3a1bc4e2e26a1c895a555d",
					Timestamp:   1646769225000,
				},
			},
			Files: []types.File{
				{
					Id:          fileName,
					ContentType: contentType,
					Name:        "purpur.jar",
					SHA512:      hash,
				},
			},
		})
	},
}

func init() {
	rootCommand.AddCommand(updateCommand)
}
