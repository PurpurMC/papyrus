package cli

import (
	"fmt"
	"github.com/purpurmc/papyrus/shared"
	"os"
)

func Add(config shared.Config, projectName string, versionName string, buildNumber string, filePath string) {
	project := createProjectIfNotExists(projectName)
	version := createVersionIfNotExists(project, versionName)

	if doesBuildExist(version, buildNumber) {
		fmt.Printf("Build %s already exists for version %s", buildNumber, versionName)
		return
	}

	jenkins := getJenkinsData(config.CLIConfig.JenkinsURL, projectName, buildNumber)

	md5 := ""
	if jenkins.Result == "SUCCESS" {
		path := fmt.Sprintf("%s/%s-%s-%s", config.StoragePath, projectName, versionName, buildNumber)
		shared.DownloadFile(replaceFilePathVariables(config.CLIConfig.JenkinsFilePath, config, project, buildNumber, filePath), path)
		md5 = shared.GetMD5(path)
	}

	build := shared.Build{
		Project: project.Name,
		Version: version.Name,
		Build: buildNumber,
		Result: jenkins.Result,
		Duration: jenkins.Duration,
		Commits: getCommits(jenkins),
		Timestamp: jenkins.Timestamp,
		MD5: md5,
		Extension: shared.After(filePath, "."),
	}

	addBuild(project, version, build)

	if config.CLIConfig.Webhook {
		sendWebhook(build)
	}
}

func Delete(deletionType string) {
	args := os.Args[3:]

	switch deletionType {
	case "project":
		if len(args) != 1 {
			shared.PrintUsage()
			return
		}

		projectName := args[0]

		fmt.Println("Deleting project " + projectName)
		deleteProject(projectName)
		fmt.Println("Project deleted")
		return
	case "version":
		if len(args) != 2 {
			shared.PrintUsage()
			return
		}

		projectName := args[0]
		versionName := args[1]

		fmt.Println("Deleting version " + versionName + " from project " + projectName)
		deleteVersion(projectName, versionName)
		fmt.Println("Version deleted")
		return
	case "build":
		if len(args) != 3 {
			shared.PrintUsage()
			return
		}

		projectName := args[0]
		versionName := args[1]
		buildNumber := args[2]

		fmt.Println("Deleting build " + buildNumber + " from version " + versionName + " in project " + projectName)
		deleteBuild(projectName, versionName, buildNumber)
		fmt.Println("Build deleted")
		return
	default:
		shared.PrintUsage()
	}
}

func TestWebhook() {
	fmt.Println("Testing webhook")
	var commits []shared.Commit
	commits = append(commits, shared.Commit{
		Author: "ben",
		Title: "test commit uno",
		Comment: "test commit uno\n\nthis is a test commit",
		Hash: "9eabf5b536662000f79978c4d1b6e4eff5c8d785",
		Email: "ben@omega24.dev",
		Timestamp: 125155125,
	})

	commits = append(commits, shared.Commit{
		Author: "frank",
		Title: "test commit dos",
		Comment: "test commit dos\n\nthis is not a test commit",
		Hash: "29932f3915935d773dc8d52c292cadd81c81071d",
		Email: "frank@gmail.com",
		Timestamp: 5919895512,
	})

	success := shared.Build{
		Project: "test",
		Version: "1.0.0",
		Build: "101",
		Result: "SUCCESS",
		Duration: 505,
		Commits: commits,
		Timestamp: 1256981234,
		MD5: "md5",
		Extension: ".jar",
	}

	failure := shared.Build{
		Project: "test",
		Version: "1.0.0",
		Build: "101",
		Result: "FAILURE",
		Duration: 505,
		Commits: commits,
		Timestamp: 1256981234,
		MD5: "md5",
		Extension: ".jar",
	}

	sendWebhook(success)
	sendWebhook(failure)
	fmt.Println("Webhook tested")
}
