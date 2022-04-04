package db

import (
	"github.com/purpurmc/papyrus/types"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProjectToResponse(database *mongo.Database, project types.Project) types.ProjectResponse {
	versions := GetVersions(database, &types.Version{ProjectId: project.Id})

	var versionNames []string
	if len(versions) > 0 {
		versionNames = make([]string, len(versions)-1)
		for _, version := range versions {
			versionNames = append(versionNames, version.Name)
		}
	} else {
		versionNames = make([]string, 0)
	}

	return types.ProjectResponse{
		Project:  project.Name,
		CreatedAt: project.CreatedAt,
		Versions: versionNames,
	}
}

func VersionToResponse(database *mongo.Database, version types.Version) types.VersionResponse {
	project := GetProject(database, &types.Project{Id: version.ProjectId})
	builds := GetBuilds(database, &types.Build{VersionId: version.Id})

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

	var buildNames []string
	if len(builds) > 0 {
		buildNames = make([]string, len(builds)-1)
		for _, build := range builds {
			buildNames = append(buildNames, build.Name)
		}
	} else {
		buildNames = make([]string, 0)
	}

	return types.VersionResponse{
		Project: project.Name,
		Version: version.Name,
		CreatedAt: version.CreatedAt,
		Latest:  newestBuild.Name,
		Builds:  buildNames,
	}
}

func VersionToResponseDetailed(database *mongo.Database, version types.Version) types.VersionResponseDetailed {
	project := GetProject(database, &types.Project{Id: version.ProjectId})
	builds := GetBuilds(database, &types.Build{VersionId: version.Id})

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

	var responseBuilds []types.BuildResponse
	if len(builds) > 0 {
		responseBuilds = make([]types.BuildResponse, len(builds)-1)
		for _, build := range builds {
			responseBuilds = append(responseBuilds, BuildToResponse(database, build))
		}
	} else {
		responseBuilds = make([]types.BuildResponse, 0)
	}

	return types.VersionResponseDetailed{
		Project: project.Name,
		Version: version.Name,
		CreatedAt: version.CreatedAt,
		Latest:  newestBuild.Name,
		Builds:  responseBuilds,
	}
}

func BuildToResponse(database *mongo.Database, build types.Build) types.BuildResponse {
	version := GetVersion(database, &types.Version{Id: build.VersionId})
	project := GetProject(database, &types.Project{Id: version.ProjectId})

	var flags []string
	if build.Flags == nil {
		flags = make([]string, 0)
	} else {
		flags = build.Flags
	}

	var commits []types.Commit
	if build.Commits == nil {
		commits = make([]types.Commit, 0)
	} else {
		commits = build.Commits
	}

	var files []types.File
	if build.Files == nil {
		files = make([]types.File, 0)
	} else {
		files = build.Files
	}

	return types.BuildResponse{
		Project:   project.Name,
		Version:   version.Name,
		Build:     build.Name,
		CreatedAt: build.CreatedAt,
		Result:    build.Result,
		Flags: flags,
		Commits:   commits,
		Files:     files,
	}
}
