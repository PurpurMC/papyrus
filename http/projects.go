package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"github.com/purpurmc/papyrus/utils"
)

func listProjects(c *gin.Context) {
	database, _ := db.NewMongo()
	defer database.Client().Disconnect(context.TODO())

	projects := db.GetProjects(database, nil)

	var projectNames []string
	if len(projects) > 0 {
		projectNames = make([]string, len(projects)-1)
		for _, project := range projects {
			projectNames = append(projectNames, project.Name)
		}
	} else {
		projectNames = make([]string, 0)
	}

	c.JSON(200, types.ProjectsResponse{Projects: projectNames})
}

func getProject(c *gin.Context) {
	database, _ := db.NewMongo()
	defer database.Client().Disconnect(context.TODO())

	project := db.GetProject(database, &types.Project{Name: c.Param("project")})
	if project == nil {
		utils.Return404(c)
		return
	}

	c.JSON(200, db.ProjectToResponse(database, *project))
}
