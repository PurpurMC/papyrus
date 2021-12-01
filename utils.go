package main

import (
	"crypto/md5"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
)

const filechunk = 8192

func printUsage() {
	color.Cyan("papyrus usage:")
	color.White("papyrus <project> <version> <build> <file_path>")
}

func setup() {
	println("Checking for config file")
	saveConfig(Config{
		JenkinsURL: "https://jenkins.example.com",
		StoragePath: "/srv/papyrus",
	})

	saveDataFile(Data{
		Projects: nil,
	})
}

func debug() {
	fmt.Printf("%+v", getDataFile())
}

func reset() {
	saveDataFile(Data{
		Projects: nil,
	})
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func downloadFile(url string, path string) {
	resp, err := http.Get(url)
	checkError(err)

	file, err := os.Create(path)
	checkError(err)

	_, err = io.Copy(file, resp.Body)
	checkError(err)

	checkError(resp.Body.Close())
}

func getMD5(path string) string {
	file, err := os.Open(path)
	checkError(err)

	hash := md5.New()
	_, err = io.Copy(hash, file)
	checkError(err)

	return fmt.Sprintf("%x", hash.Sum(nil))
}
