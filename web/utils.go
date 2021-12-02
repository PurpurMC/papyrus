package web

import (
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/shared"
	"os"
)

func getVersions(project shared.Project) []string {
	var versions []string
	for _, version := range project.Versions {
		versions = append(versions, version.Name)
	}
	return versions
}

func getBuilds(version shared.Version) []int {
	var builds []int
	for _, build := range version.Builds {
		builds = append(builds, build.Build)
	}
	return builds
}

// get commits from build
func getCommits(build shared.Build) []gin.H {
	var commits []gin.H
	for _, commit := range build.Commits {
		commits = append(commits, gin.H{
			"author": commit.Author,
			"description": commit.Comment,
			"hash": commit.Hash,
			"email": commit.Email,
			"timestamp": commit.Timestamp,
		})
	}

	if commits == nil {
		return []gin.H{}
	}

	return commits
}

func getFileSize(file *os.File) int64 {
	stat, err := file.Stat()
	if err != nil {
		return 0
	}

	return stat.Size()
}
