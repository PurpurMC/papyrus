package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Data struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	Name string `json:"name"`
	Versions []Version `json:"versions"`
}

type Version struct {
	Name string `json:"name"`
	Builds []Build `json:"builds"`
}

type Build struct {
	Project string `json:"project"`
	Version string `json:"version"`
	Build int `json:"build"`
	Result string `json:"result"`
	Duration int `json:"duration"`
	Commits []Commit `json:"commits"`
	Timestamp int `json:"timestamp"`
	MD5 string `json:"md5"`
}

type Commit struct {
	Author string `json:"author"`
	Description string `json:"description"`
	Hash string `json:"hash"`
	Email string `json:"email"`
	Timestamp int `json:"timestamp"`
}

func saveDataFile(data Data) {
	jsonBytes, err := json.Marshal(data)
	checkError(err)

	checkError(ioutil.WriteFile("data.json", jsonBytes, os.ModePerm)) // todo: location
}

func getDataFile() Data {
	file, err := ioutil.ReadFile("data.json")
	checkError(err)

	data := Data{}
	checkError(json.Unmarshal(file, &data))
	return data
}

func createProjectIfNotExists(projectName string) Project {
	data := getDataFile()
	for _, project := range data.Projects {
		if project.Name == projectName {
			return project
		}
	}

	project := Project{Name: projectName}
	data.Projects = append(data.Projects, project)
	saveDataFile(data)
	return project
}

func saveProject(project Project) {
	data := getDataFile()
	for i, p := range data.Projects {
		if p.Name == project.Name {
			data.Projects[i] = project
			saveDataFile(data)
			return
		}
	}

	data.Projects = append(data.Projects, project)
	saveDataFile(data)
}

func createVersionIfNotExists(project Project, versionName string) Version {
	for _, version := range project.Versions {
		if version.Name == versionName {
			return version
		}
	}

	version := Version{Name: versionName}
	project.Versions = append(project.Versions, version)
	saveProject(project)
	return version
}

func saveVersion(project Project, version Version) {
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

func addBuild(project Project, version Version, build Build) {
	version.Builds = append(version.Builds, build)
	saveVersion(project, version)
}

func doesBuildExist(version Version, buildNumber int) bool {
	for _, build := range version.Builds {
		if build.Build == buildNumber {
			return true
		}
	}

	return false
}
