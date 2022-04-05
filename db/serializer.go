package db

import (
	"github.com/purpurmc/papyrus/types"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProjectToResponse(database *mongo.Database, project types.Project) types.ProjectResponse {
	versions := GetVersions(database, &types.Version{ProjectId: project.Id})

	versionNames := make([]string, 0)
	for _, version := range versions {
		versionNames = append(versionNames, version.Name)
	}

	return types.ProjectResponse{
		Project:   project.Name,
		CreatedAt: project.CreatedAt,
		Versions:  versionNames,
	}
}

func VersionToResponse(database *mongo.Database, version types.Version) types.VersionResponse {
	project := GetProject(database, &types.Project{Id: version.ProjectId})
	builds := GetBuildsFromVersion(database, version.Id)

	var newestBuild types.Build
	if len(builds) > 0 {
		newestBuild = builds[0]
		for _, build := range builds {
			if build.CreatedAt > newestBuild.CreatedAt {
				newestBuild = build
			}
		}
	} else {
		newestBuild = types.Build{Name: ""}
	}

	buildNames := make([]string, 0)
	for _, build := range builds {
		buildNames = append(buildNames, build.Name)
	}

	return types.VersionResponse{
		Project:   project.Name,
		Version:   version.Name,
		CreatedAt: version.CreatedAt,
		Builds: struct {
			Latest string   `json:"latest"`
			All    []string `json:"all"`
		}{
			Latest: newestBuild.Name,
			All:    buildNames,
		},
	}
}

func VersionToResponseDetailed(database *mongo.Database, version types.Version) types.VersionResponseDetailed {
	project := GetProject(database, &types.Project{Id: version.ProjectId})
	builds := GetBuildsFromVersion(database, version.Id)

	var newestBuild types.Build
	if len(builds) > 0 {
		newestBuild = builds[0]
		for _, build := range builds {
			if build.CreatedAt > newestBuild.CreatedAt {
				newestBuild = build
			}
		}
	} else {
		newestBuild = types.Build{Name: ""}
	}

	responseBuilds := make([]types.BuildResponse, 0)
	for _, build := range builds {
		responseBuilds = append(responseBuilds, BuildToResponse(database, build))
	}

	return types.VersionResponseDetailed{
		Project:   project.Name,
		Version:   version.Name,
		CreatedAt: version.CreatedAt,
		Builds: struct {
			Latest string                `json:"latest"`
			All    []types.BuildResponse `json:"all"`
		}{
			Latest: newestBuild.Name,
			All:    responseBuilds,
		},
	}
}

func BuildToResponse(database *mongo.Database, build types.Build) types.BuildResponse {
	versions := GetVersionsFromBuild(database, build.Id)
	project := GetProject(database, &types.Project{Id: versions[0].ProjectId})

	versionNames := make([]string, 0)
	for _, version := range versions {
		versionNames = append(versionNames, version.Name)
	}

	flags := make([]string, 0)
	if build.Flags != nil {
		flags = build.Flags
	}

	commits := make([]types.Commit, 0)
	if build.Commits != nil {
		commits = build.Commits
	}

	files := make([]types.File, 0)
	if build.Files != nil {
		files = build.Files
	}

	return types.BuildResponse{
		Project:   project.Name,
		Versions:  versionNames,
		Build:     build.Name,
		CreatedAt: build.CreatedAt,
		Result:    build.Result,
		Flags:     flags,
		Commits:   commits,
		Files:     files,
	}
}
