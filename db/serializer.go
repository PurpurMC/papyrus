package db

import (
	"github.com/purpurmc/papyrus/types"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProjectToResponse(database *mongo.Database, project types.Project) types.ProjectResponse {
	versions := GetVersions(database, &types.Version{ProjectId: project.Id})

	versionNames := make([]string, len(versions)-1)
	for _, version := range versions {
		versionNames = append(versionNames, version.Name)
	}

	return types.ProjectResponse{
		Project:  project.Name,
		Versions: versionNames,
	}
}

func VersionToResponse(database *mongo.Database, version types.Version) types.VersionResponse {
	project := GetProject(database, &types.Project{Id: version.ProjectId})
	builds := GetBuilds(database, &types.Build{VersionId: version.Id})

	newestBuild := builds[0]
	for _, build := range builds {
		if build.CreatedAt > newestBuild.CreatedAt {
			newestBuild = build
		}
	}

	buildNames := make([]string, len(builds)-1)
	for _, build := range builds {
		buildNames = append(buildNames, build.Name)
	}

	return types.VersionResponse{
		Project: project.Name,
		Version: version.Name,
		Latest:  newestBuild.Name,
		Builds:  buildNames,
	}
}

func VersionToResponseDetailed(database *mongo.Database, version types.Version) types.VersionResponseDetailed {
	project := GetProject(database, &types.Project{Id: version.ProjectId})
	builds := GetBuilds(database, &types.Build{VersionId: version.Id})

	newestBuild := builds[0]
	for _, build := range builds {
		if build.CreatedAt > newestBuild.CreatedAt {
			newestBuild = build
		}
	}

	responseBuilds := make([]types.BuildResponse, len(builds)-1)
	for _, build := range builds {
		responseBuilds = append(responseBuilds, BuildToResponse(database, build))
	}

	return types.VersionResponseDetailed{
		Project: project.Name,
		Version: version.Name,
		Latest:  newestBuild.Name,
		Builds:  responseBuilds,
	}
}

func BuildToResponse(database *mongo.Database, build types.Build) types.BuildResponse {
	version := GetVersion(database, &types.Version{Id: build.VersionId})
	project := GetProject(database, &types.Project{Id: version.ProjectId})

	return types.BuildResponse{
		Project:   project.Name,
		Version:   version.Name,
		Build:     build.Name,
		CreatedAt: build.CreatedAt,
		Result:    build.Result,
		Commits:   build.Commits,
		Files:     build.Files,
	}
}
