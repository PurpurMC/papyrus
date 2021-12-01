package cli

import (
	"fmt"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/webhook"
	"github.com/fatih/color"
	"github.com/purpurmc/papyrus/shared"
	"strings"
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

		var embedSettings shared.EmbedConfig
		if build.Result == "SUCCESS" {
			embedSettings = config.CLIConfig.SuccessEmbed
		} else {
			embedSettings = config.CLIConfig.FailureEmbed
		}

		var embed []discord.Embed
		embed = append(embed, discord.NewEmbedBuilder().
			SetTitle(replaceVariables(embedSettings.Title, embedSettings.Changes, build)).
			SetDescription(replaceVariables(embedSettings.Description, embedSettings.Changes, build)).
			SetColor(embedSettings.Color).
			Build())

		_, err := client.CreateEmbeds(embed)
		if err != nil {
			panic(err)
		}
	}
}

func replaceVariables(template string, changesTemplate string, build shared.Build) string {
	replaced := strings.ReplaceAll(template, "{project}", build.Project)
	replaced = strings.ReplaceAll(replaced, "{version}", build.Version)
	replaced = strings.ReplaceAll(replaced, "{build}", fmt.Sprintf("%d", build.Build))
	replaced = strings.ReplaceAll(replaced, "{result}", build.Result)
	replaced = strings.ReplaceAll(replaced, "{duration}", fmt.Sprintf("%d", build.Duration))
	replaced = strings.ReplaceAll(replaced, "{changes}", generateChanges(changesTemplate, build))
	replaced = strings.ReplaceAll(replaced, "{timestamp}", fmt.Sprintf("%d", build.Timestamp))
	replaced = strings.ReplaceAll(replaced, "{md5}", build.MD5)
	return replaced
}

func generateChanges(template string, build shared.Build) string {
	var changes string
	for _, commit := range build.Commits {
		changes += replaceChangesVariables(template, commit)
	}
	return changes
}

func replaceChangesVariables(template string, commit shared.Commit) string {
	replaced := strings.ReplaceAll(template, "{author}", commit.Author)
	replaced = strings.ReplaceAll(replaced, "{title}", shared.Before(commit.Description, string('\n')))
	replaced = strings.ReplaceAll(replaced, "{description}", shared.After(commit.Description, string('\n')))
	replaced = strings.ReplaceAll(replaced, "{timestamp}", fmt.Sprintf("%d", commit.Timestamp))
	replaced = strings.ReplaceAll(replaced, "{hash}", commit.Hash)
	replaced = strings.ReplaceAll(replaced, "{short_hash}", shared.First(commit.Hash, 7))
	replaced = strings.ReplaceAll(replaced, "{email}", commit.Email)
	replaced = strings.ReplaceAll(replaced, "{timestamp}", fmt.Sprintf("%d", commit.Timestamp))
	return replaced
}
