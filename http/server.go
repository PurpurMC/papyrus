package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

	router.RedirectTrailingSlash = true

	fmt.Println(fmt.Sprintf("Listening on %s:%d", host, port))
	if err := router.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		panic(err)
	}
}
