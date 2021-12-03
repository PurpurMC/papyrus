package web

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/shared"
)

func Web(config shared.Config) {
	if !config.WebConfig.Dev {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("docs/", false)))
	router.GET("/v2", getProjects)
	router.GET("/v2/:project", getProject)
	router.GET("/v2/:project/:version", getVersion)
	router.GET("/v2/:project/:version/:build", getBuild)
	router.GET("/v2/:project/:version/:build/download", downloadBuild)

	err := router.Run(config.WebConfig.IP)
	if err != nil {
		panic(err)
	}
}
