package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/shared"
	"strconv"
)

func getBuild(c *gin.Context) {
	data := shared.GetData()
	projectName := c.Param("project")
	versionName := c.Param("version")
	buildNumber := c.Param("build")

	for _, project := range data.Projects {
		if project.Name == projectName {
			for _, version := range project.Versions {
				if version.Name == versionName {
					var build *shared.Build

					if buildNumber == "latest" {
						build = &version.Latest
					} else {
						for _, versionBuild := range version.Builds {
							if strconv.Itoa(versionBuild.Build) == buildNumber {
								build = &versionBuild
							}
						}

						if build == nil {
							c.JSON(404, gin.H{
								"error": "build not found",
							})
							return
						}
					}

					c.JSON(200, gin.H{
						"project": project.Name,
						"version": version.Name,
						"build": build.Build,
						"result": build.Result,
						"duration": build.Duration,
						"commits": getCommits(*build),
						"timestamp": build.Timestamp,
						"md5": build.MD5,
					})
				}
			}
		}
	}
}

func downloadBuild(c *gin.Context) {
	data := shared.GetData()
	projectName := c.Param("project")
	versionName := c.Param("version")
	buildNumber := c.Param("build")

	for _, project := range data.Projects {
		if project.Name == projectName {
			for _, version := range project.Versions {
				if version.Name == versionName {
					var build *shared.Build

					if buildNumber == "latest" {
						build = &version.Latest
					} else {
						for _, versionBuild := range version.Builds {
							if strconv.Itoa(versionBuild.Build) == buildNumber {
								build = &versionBuild
							}
						}

						if build == nil {
							c.JSON(404, gin.H{
								"error": "build not found",
							})
							return
						}
					}

					path := fmt.Sprintf("%s-%s-%d", projectName, versionName, build.Build)
					c.Header("Content-Type", "application/jar")
					c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.jar", path)) // todo: extension
					// todo: content length
					c.File(path)
				}
			}
		}
	}
}
