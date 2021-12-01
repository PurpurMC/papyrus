package web

import (
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/shared"
)

func getVersion(c *gin.Context) {
	data := shared.GetData()
	projectName := c.Param("project")
	versionName := c.Param("version")

	for _, project := range data.Projects {
		if project.Name == projectName {
			for _, version := range project.Versions {
				if version.Name == versionName {
					c.JSON(200, gin.H{
						"project": project.Name,
						"version": version.Name,
						"builds": gin.H{
							"latest": version.Latest.Build,
							"all": getBuilds(version),
						},
					})
					return
				}
			}
		}
	}

	c.JSON(404, gin.H{
		"error": "version not found",
	})
}
