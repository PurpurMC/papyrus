package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"github.com/purpurmc/papyrus/utils"
)

func GetBuild(c *gin.Context) {
	database, _ := db.NewMongo()
	defer database.Client().Disconnect(context.TODO())

	project := db.GetProject(database, &types.Project{Name: c.Param("project")})
	if project == nil {
		utils.Return404(c)
		return
	}

	version := db.GetVersion(database, &types.Version{
		ProjectId: project.Id,
		Name:      c.Param("version"),
	})

	if version == nil {
		utils.Return404(c)
		return
	}

	var build *types.Build
	if c.Param("build") == "latest" {
		build = db.GetBuild(database, &types.Build{
			VersionId: version.Id,
			Name:      db.VersionToResponse(database, *version).Builds.Latest,
		})
	} else {
		build = db.GetBuild(database, &types.Build{
			VersionId: version.Id,
			Name:      c.Param("build"),
		})
	}

	if build == nil {
		utils.Return404(c)
		return
	}

	commits := make([]gin.H, 0)
	for _, commit := range build.Commits {
		commits = append(commits, gin.H{
			"author":      commit.Author,
			"description": commit.Description,
			"hash":        commit.Hash,
			"email":       commit.Email,
			"timestamp":   commit.Timestamp,
		})
	}

	v1 := database.Collection("v1")
	var legacyData LegacyBuildData
	err := v1.FindOne(context.TODO(), &LegacyBuildData{BuildId: build.Id}).Decode(&legacyData)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"project":   project.Name,
		"version":   version.Name,
		"build":     build.Name,
		"result":    build.Result,
		"duration":  0,
		"commits":   commits,
		"timestamp": build.CreatedAt,
		"md5": legacyData.MD5,
	})
}

func DownloadBuild(c *gin.Context) {
	database, bucket := db.NewMongo()
	defer database.Client().Disconnect(context.TODO())

	project := db.GetProject(database, &types.Project{Name: c.Param("project")})
	if project == nil {
		utils.Return404(c)
		return
	}

	version := db.GetVersion(database, &types.Version{
		ProjectId: project.Id,
		Name:      c.Param("version"),
	})

	if version == nil {
		utils.Return404(c)
		return
	}

	var build *types.Build
	if c.Param("build") == "latest" {
		build = db.GetBuild(database, &types.Build{
			VersionId: version.Id,
			Name:      db.VersionToResponse(database, *version).Builds.Latest,
		})
	} else {
		build = db.GetBuild(database, &types.Build{
			VersionId: version.Id,
			Name:      c.Param("build"),
		})
	}

	if build == nil {
		utils.Return404(c)
		return
	}

	data := db.DownloadFile(bucket, build.Files[0].Id)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", build.Files[0].Name))
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	c.Data(200, build.Files[0].ContentType, data)
}
