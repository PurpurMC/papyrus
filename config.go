package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	JenkinsURL string `json:"jenkins_url"`
	StoragePath string `json:"storage_path"`
}

func saveConfig(config Config)  {
	jsonBytes, err := json.MarshalIndent(config, "", "\t")
	checkError(err)

	checkError(ioutil.WriteFile("config.json", jsonBytes, os.ModePerm)) // todo: location
}

func getConfig() Config {
	file, err := ioutil.ReadFile("config.json")
	checkError(err)

	config := Config{}
	checkError(json.Unmarshal(file, &config))
	return config
}
