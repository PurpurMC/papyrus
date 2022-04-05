package db

import (
	"github.com/purpurmc/papyrus/types"
	"github.com/purpurmc/papyrus/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetBuildsFromVersion(database *mongo.Database, versionId primitive.ObjectID) []types.Build {
	version := GetVersion(database, &types.Version{Id: versionId})
	allBuilds := GetBuilds(database, nil)
	builds := make([]types.Build, 0)
	for _, build := range allBuilds {
		if utils.ObjectIdInSlice(version.Id, build.VersionIds) {
			builds = append(builds, build)
		}
	}
	return builds
}

func GetVersionsFromBuild(database *mongo.Database, buildId primitive.ObjectID) []types.Version {
	build := GetBuild(database, &types.Build{Id: buildId})
	allVersions := GetVersions(database, nil)
	versions := make([]types.Version, 0)
	for _, version := range allVersions {
		if utils.ObjectIdInSlice(version.Id, build.VersionIds) {
			versions = append(versions, version)
		}
	}
	return versions
}
