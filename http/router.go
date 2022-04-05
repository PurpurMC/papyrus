package http

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/v1"
	"github.com/spf13/viper"
	"strings"
)

func Start() {
	host := viper.GetString("http.host")
	port := viper.GetInt("http.port")

	if !viper.GetBool("http.debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	prefix := router.Group(viper.GetString("http.routes.prefix"))
	{
		prefix.GET(strings.TrimSuffix(viper.GetString("http.routes.list-projects"), "/"), listProjects)
		prefix.GET(strings.TrimSuffix(viper.GetString("http.routes.get-project"), "/"), getProject)
		prefix.GET(strings.TrimSuffix(viper.GetString("http.routes.get-version"), "/"), getVersion)
		prefix.GET(strings.TrimSuffix(viper.GetString("http.routes.get-build"), "/"), getBuild)
		prefix.GET(strings.TrimSuffix(viper.GetString("http.routes.download-build"), "/"), downloadBuild)
	}

	if viper.GetBool("http.routes.docs.enabled") {
		router.Use(static.Serve(viper.GetString("http.routes.docs.prefix"), static.LocalFile(viper.GetString("http.routes.docs.directory"), true)))
	}

	if viper.GetBool("http.v1-compat.enabled") {
		compat := router.Group(viper.GetString("http.v1-compat.prefix"))
		{
			compat.GET("", listProjects)
			compat.GET("/:project", getProject)
			compat.GET("/:project/:version", getVersion)
			compat.GET("/:project/:version/:build", v1.GetBuild)
			compat.GET("/:project/:version/:build/download", v1.DownloadBuild)
		}
	}

	router.RedirectTrailingSlash = true

	fmt.Println(fmt.Sprintf("Listening on %s:%d", host, port))
	if err := router.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		panic(err)
	}
}
