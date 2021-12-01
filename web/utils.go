package web

import "github.com/purpurmc/papyrus/shared"

func getVersions(project shared.Project) []string {
	var versions []string
	for _, version := range project.Versions {
		versions = append(versions, version.Name)
	}
	return versions
}
