package cli

import (
	"fmt"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/webhook"
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
		shared.DownloadFile(fmt.Sprintf("%s/job/%s/%d/artifact/%s", config.CLIConfig.JenkinsURL, projectName, buildNumber, filePath), path) // todo: custom path
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

	if config.CLIConfig.Webhook {
		client := webhook.NewClient(discord.Snowflake(config.CLIConfig.WebhookID), config.CLIConfig.WebhookToken)
		_, err := client.CreateContent(fmt.Sprintf("Build %d for %s %s (%d) has been uploaded", buildNumber, projectName, versionName, build.Duration))
		if err != nil {
			panic(err)
		}
	}
}
