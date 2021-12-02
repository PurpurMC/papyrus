package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/shared"
	"os"
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

					println(build.Result)
					if build.Result != "SUCCESS" {
						c.JSON(404, gin.H{
							"error": "build failed, nothing to download",
						})
						return
					}

					config := shared.GetConfig()
					fileName := fmt.Sprintf("%s-%s-%d", projectName, versionName, build.Build)
					file, _ := os.Open(fmt.Sprintf("%s/%s", config.StoragePath, fileName))

					c.Header("Content-Type", "application/jar")
					c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", fileName, build.Extension))
					c.Header("Content-Length", strconv.FormatInt(getFileSize(file), 10))
					c.File(fmt.Sprintf("%s/%s", config.StoragePath, fileName))
				}
			}
		}
	}
}
