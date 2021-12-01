package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) < 3 || len(args) > 4 {
		if len(args) == 1 {
			if args[0] == "setup" {
				setup()
				return
			}

			if args[0] == "debug" {
				debug()
				return
			}

			if args[0] == "reset" {
				reset()
				return
			}
		}

		printUsage()
		return
	}

	projectName := args[0]
	versionName := args[1]
	buildNumber, err := strconv.Atoi(args[2])
	filePath := args[3]

	if err != nil {
		color.Red("%s is not a valid number", args[2])
		return
	}

	if projectName == "" || versionName == "" || buildNumber <= 0 || filePath == "" {
		printUsage()
		return
	}

	project := createProjectIfNotExists(projectName)
	version := createVersionIfNotExists(project, versionName)

	if doesBuildExist(version, buildNumber) {
		color.Red("Build %d already exists for the project %s", buildNumber, projectName)
		return
	}

	config := getConfig()
	jenkins := getJenkinsData(config.JenkinsURL, projectName, buildNumber)

	md5 := ""
	if jenkins.Result == "SUCCESS" {
		path := fmt.Sprintf("%s-%s-%d", projectName, versionName, buildNumber)
		downloadFile(fmt.Sprintf("%s/job/%s/ws/%s", config.JenkinsURL, projectName, filePath), path)
		md5 = getMD5(path)
	}

	println(md5)

	build := Build{
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

	fmt.Printf("%+v", build)
}
