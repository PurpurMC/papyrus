package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"github.com/purpurmc/papyrus/utils"
)

func getVersion(c *gin.Context) {
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

	if c.DefaultQuery("detailed", "false") == "true" {
		c.JSON(200, db.VersionToResponseDetailed(database, *version))
	} else {
		c.JSON(200, db.VersionToResponse(database, *version))
	}
}
