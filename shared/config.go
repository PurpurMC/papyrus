package shared

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	StoragePath string `json:"storage_path"`
	CLIConfig CLIConfig `json:"cli"`
	WebConfig WebConfig `json:"web"`
}

type CLIConfig struct {
	JenkinsURL string `json:"jenkins_url"`
	JenkinsFilePath string `json:"jenkins_file_path"`
	PostbuildScript string `json:"postbuild_script"`
}

type WebConfig struct {
	IP string `json:"ip"`
	Dev bool `json:"dev"`
}

func GetConfig() Config {
	file, err := ioutil.ReadFile("/etc/papyrus.json")
	config := Config{}
	err = json.Unmarshal(file, &config)

	if err != nil {
		panic(err)
	}
	return config
}

func SaveConfig(config Config)  {
	jsonBytes, err := json.MarshalIndent(config, "", "\t")
	err = ioutil.WriteFile("/etc/papyrus.json", jsonBytes, os.ModePerm)

	if err != nil {
		panic(err)
	}
}
