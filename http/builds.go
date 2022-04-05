package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"github.com/purpurmc/papyrus/utils"
)

func getBuild(c *gin.Context) {
	database, _ := db.NewMongo()
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

	builds := db.GetBuildsFromVersion(database, version.Id)
	var build *types.Build
	if c.Param("build") == "latest" {
		versionRespone := db.VersionToResponse(database, *version)
		for _, b := range builds {
			if b.Name == versionRespone.Builds.Latest {
				build = &b
				break
			}
		}
	} else {
		for _, b := range builds {
			if b.Name == c.Param("build") {
				build = &b
				break
			}
		}
	}

	if build == nil {
		utils.Return404(c)
		return
	}

	c.JSON(200, db.BuildToResponse(database, *build))
}

func downloadBuild(c *gin.Context) {
	database, bucket := db.NewMongo()
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

	builds := db.GetBuildsFromVersion(database, version.Id)
	var build *types.Build
	if c.Param("build") == "latest" {
		versionRespone := db.VersionToResponse(database, *version)
		for _, b := range builds {
			if b.Name == versionRespone.Builds.Latest {
				build = &b
				break
			}
		}
	} else {
		for _, b := range builds {
			if b.Name == c.Param("build") {
				build = &b
				break
			}
		}
	}

	if build == nil {
		utils.Return404(c)
		return
	}

	var file *types.File
	for _, f := range build.Files {
		if f.Name == c.Param("file") {
			file = &f
			break
		}
	}

	if file == nil {
		utils.Return404(c)
		return
	}

	var data = db.DownloadFile(bucket, file.Id)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	c.Data(200, file.ContentType, data)
}
