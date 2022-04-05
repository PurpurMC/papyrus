package cmd

import (
	"context"
	"fmt"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

var updateCommand = &cobra.Command{
	Use:   "update",
	Short: "After updating papyrus, this command will update the config and the db",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating config...")
		viper.SetDefault("_version", 1)

		viper.SetDefault("http.host", "127.0.0.1")
		viper.SetDefault("http.port", 8080)
		viper.SetDefault("http.debug", false)

		viper.SetDefault("http.routes.prefix", "/v2")
		viper.SetDefault("http.routes.list-projects", "/")
		viper.SetDefault("http.routes.get-project", "/:project")
		viper.SetDefault("http.routes.get-version", "/:project/:version")
		viper.SetDefault("http.routes.get-build", "/:project/:version/:build")
		viper.SetDefault("http.routes.download-build", "/:project/:version/:build/:file")

		viper.SetDefault("http.routes.docs.enabled", false)
		viper.SetDefault("http.routes.docs.prefix", "/docs")
		viper.SetDefault("http.routes.docs.directory", "/")

		viper.SetDefault("http.v1-compat.enabled", false)
		viper.SetDefault("http.v1-compat.prefix", "/v1")

		viper.SetDefault("db.host", "127.0.0.1")
		viper.SetDefault("db.port", 27017)
		viper.SetDefault("db.username", "username")
		viper.SetDefault("db.password", "password")
		viper.SetDefault("db.db", "papyrus")

		viper.SetDefault("utils.cloudflare-access-token", "")

		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error writing config file, please make sure you have the correct permissions")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Updating database...")
		database, _ := db.NewMongo()
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

		if !utils.StringInSlice("v1", collections) {
			db.CreateCollection(database, "v1")
		}

		fmt.Println("Done!")
	},
}

func init() {
	rootCommand.AddCommand(updateCommand)
}
