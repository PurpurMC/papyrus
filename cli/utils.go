package cli

import (
	"encoding/json"
	"github.com/purpurmc/papyrus/shared"
	"io/ioutil"
	"net/http"
	"os/exec"
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

func runPostbuildScript(build shared.Build, script string) {
	bytes, err := json.Marshal(build)

	if err != nil {
		panic(err)
	}

	command := strings.Split(strings.ReplaceAll(script, "{data}", string(bytes)), " ")
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	println(string(output))
}
