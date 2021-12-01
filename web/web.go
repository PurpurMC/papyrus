package web

import (
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/shared"
)

func Web(config shared.Config) {
	if !config.WebConfig.Dev {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.GET("/v1", getProjects)
	router.GET("/v1/:project", getProject)
	router.GET("/v1/:project/:version", getVersion)

	err := router.Run(config.WebConfig.IP)
	if err != nil {
		panic(err)
	}
}
