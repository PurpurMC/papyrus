package cli

import (
	"encoding/json"
	"fmt"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/webhook"
	"github.com/purpurmc/papyrus/shared"
	"io/ioutil"
	"net/http"
	"strings"
)

type JenkinsData struct {
	Result string `json:"result"`
	Duration int `json:"duration"`
	Timestamp int `json:"timestamp"`
	ChangeSet ChangeSet `json:"changeSet"`
}

type ChangeSet struct {
	Commit []Commit `json:"items"`
}

type Commit struct {
	Author Author `json:"author"`
	Title string `json:"msg"`
	Comment string `json:"comment"`
	Hash string `json:"commitId"`
	Email string `json:"authorEmail"`
	Timestamp int `json:"timestamp"`
}

type Author struct {
	Name string `json:"fullName"`
}

func getJenkinsData(url string, project string, build string) JenkinsData {
	var responseObject JenkinsData
	response, err := http.Get(url + "/job/" + project + "/" + build + "/api/json")
	responseData, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(responseData, &responseObject)

	if err != nil {
		panic(err)
	}
	return responseObject
}

func getCommits(data JenkinsData) []shared.Commit {
	var commits []shared.Commit
	for _, commit := range data.ChangeSet.Commit {
		commits = append(commits, shared.Commit{
			Author: commit.Author.Name,
			Title: commit.Title,
			Comment: commit.Comment,
			Hash: commit.Hash,
			Email: commit.Email,
			Timestamp: commit.Timestamp,
		})
	}
	return commits
}

func createProjectIfNotExists(projectName string) shared.Project {
	data := shared.GetData()
	for _, project := range data.Projects {
		if project.Name == projectName {
			return project
		}
	}

	project := shared.Project{Name: projectName}
	data.Projects = append(data.Projects, project)
	shared.SaveData(data)
	return project
}

func createVersionIfNotExists(project shared.Project, versionName string) shared.Version {
	for _, version := range project.Versions {
		if version.Name == versionName {
			return version
		}
	}

	version := shared.Version{Name: versionName}
	project.Versions = append(project.Versions, version)
	saveProject(project)
	return version
}

func addBuild(project shared.Project, version shared.Version, build shared.Build) {
	version.Builds = append(version.Builds, build)
	if build.Result == "SUCCESS" {
		version.Latest = build
	}
	saveVersion(project, version)
}

func doesBuildExist(version shared.Version, buildNumber string) bool {
	for _, build := range version.Builds {
		if build.Build == buildNumber {
			return true
		}
	}

	return false
}

func saveProject(project shared.Project) {
	data := shared.GetData()
	for i, p := range data.Projects {
		if p.Name == project.Name {
			data.Projects[i] = project
			shared.SaveData(data)
			return
		}
	}

	data.Projects = append(data.Projects, project)
	shared.SaveData(data)
}

func saveVersion(project shared.Project, version shared.Version) {
	for i, v := range project.Versions {
		if v.Name == version.Name {
			project.Versions[i] = version
			saveProject(project)
			return
		}
	}

	project.Versions = append(project.Versions, version)
	saveProject(project)
}

func deleteProject(projectName string) {
	data := shared.GetData()
	for i, project := range data.Projects {
		if project.Name == projectName {
			data.Projects = append(data.Projects[:i], data.Projects[i+1:]...)
			shared.SaveData(data)
			return
		}
	}
}

func deleteVersion(projectName string, versionName string) {
	data := shared.GetData()
	for _, project := range data.Projects {
		if project.Name == projectName {
			for j, version := range project.Versions {
				if version.Name == versionName {
					project.Versions = append(project.Versions[:j], project.Versions[j+1:]...)
					saveProject(project)
					return
				}
			}
		}
	}
}

func deleteBuild(projectName string, versionName string, buildNumber string) {
	data := shared.GetData()
	for _, project := range data.Projects {
		if project.Name == projectName {
			for _, version := range project.Versions {
				if version.Name == versionName {
					for i, build := range version.Builds {
						if build.Build == buildNumber {
							version.Builds = append(version.Builds[:i], version.Builds[i+1:]...)
							saveVersion(project, version)
							return
						}
					}
				}
			}
		}
	}
}

func replaceFilePathVariables(template string, config shared.Config, project shared.Project, build string, file string) string {
	replaced := strings.ReplaceAll(template, "{url}", config.CLIConfig.JenkinsURL)
	replaced = strings.ReplaceAll(replaced, "{project}", project.Name)
	replaced = strings.ReplaceAll(replaced, "{build}", build)
	replaced = strings.ReplaceAll(replaced, "{file}", file)
	return replaced
}

func replaceVariables(template string, changesTemplate string, build shared.Build) string {
	replaced := strings.ReplaceAll(template, "{project}", strings.Title(build.Project))
	replaced = strings.ReplaceAll(replaced, "{version}", build.Version)
	replaced = strings.ReplaceAll(replaced, "{build}", build.Build)
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
	replaced = strings.ReplaceAll(replaced, "{title}", commit.Title)
	replaced = strings.ReplaceAll(replaced, "{description}", commit.Comment)
	replaced = strings.ReplaceAll(replaced, "{timestamp}", fmt.Sprintf("%d", commit.Timestamp))
	replaced = strings.ReplaceAll(replaced, "{hash}", commit.Hash)
	replaced = strings.ReplaceAll(replaced, "{short_hash}", shared.First(commit.Hash, 7))
	replaced = strings.ReplaceAll(replaced, "{email}", commit.Email)
	replaced = strings.ReplaceAll(replaced, "{timestamp}", fmt.Sprintf("%d", commit.Timestamp))
	return replaced
}

func sendWebhook(build shared.Build) {
	config := shared.GetConfig()
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
