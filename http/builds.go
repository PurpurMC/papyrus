package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"github.com/purpurmc/papyrus/utils"
)

func getBuild(c *gin.Context) {
	database := db.NewMongo()
	defer database.Client().Disconnect(context.TODO())

	project := db.GetProject(database, &types.Project{Name: c.Param("project")})
	if project == nil {
		utils.Return404(c)
		return
	}

	version := db.GetVersion(database, &types.Version{
		ProjectId: project.Id,
		Name: c.Param("version"),
	})

	if version == nil {
		utils.Return404(c)
		return
	}

	build := db.GetBuild(database, &types.Build{
		VersionId: version.Id,
		Name: c.Param("build"),
	})

	if build == nil {
		utils.Return404(c)
		return
	}

	c.JSON(200, db.BuildToResponse(database, *build))
}

func downloadBuild(c *gin.Context) {
	project := c.Param("project")
	version := c.Param("version")
	build := c.Param("build")
	file := c.Param("file")

	c.JSON(200, gin.H{
		"project": project,
		"version": version,
		"build":   build,
		"file":    file,
	})
}
