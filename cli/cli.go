package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/purpurmc/papyrus/shared"
)

func Run(config shared.Config, projectName string, versionName string, buildNumber int, filePath string) {
	project := createProjectIfNotExists(projectName)
	version := createVersionIfNotExists(project, versionName)

	if doesBuildExist(version, buildNumber) {
		color.Red("Build %d already exists for the version %s", buildNumber, versionName)
		return
	}

	jenkins := getJenkinsData(config.CLIConfig.JenkinsURL, projectName, buildNumber)

	md5 := ""
	if jenkins.Result == "SUCCESS" {
		path := fmt.Sprintf("%s-%s-%d", projectName, versionName, buildNumber)
		shared.DownloadFile(fmt.Sprintf("%s/job/%s/ws/%s", config.CLIConfig.JenkinsURL, projectName, filePath), path) // todo: custom path
		md5 = shared.GetMD5(path)
	}

	build := shared.Build{
		Project: project.Name,
		Version: version.Name,
		Build: buildNumber,
		Result: jenkins.Result,
		Duration: jenkins.Duration,
		Commits: nil, // todo
		Timestamp: jenkins.Timestamp,
		MD5: md5,
	}

	addBuild(project, version, build)
}
