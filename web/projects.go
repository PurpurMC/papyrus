package web

import (
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/shared"
)

func getProjects(c *gin.Context) {
	data := shared.GetData()
	var projects []string

	for _, project := range data.Projects {
		projects = append(projects, project.Name)
	}

	c.JSON(200, gin.H{
		"projects": projects,
	})
}

func getProject(c *gin.Context) {
	data := shared.GetData()
	projectName := c.Param("project")

	for _, project := range data.Projects {
		if project.Name == projectName {
			c.JSON(200, gin.H{
				"project": project.Name,
				"versions": getVersions(project),
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"error": "project not found",
	})
}

