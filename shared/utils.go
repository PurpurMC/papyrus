package shared

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Setup todo: cli & web
func Setup() {
	SaveConfig(Config{
		StoragePath: "/srv/papyrus",
		CLIConfig: CLIConfig{
			JenkinsURL: "https://jenkins.example.com",
		},
	})

	SaveData(Data{
		Projects: nil,
	})
}

func Reset() {
	SaveData(Data{
		Projects: nil,
	})
}

// PrintDebug todo: cli & web
func PrintDebug() {
	fmt.Printf("%+v", GetData())
}

func GetMD5(path string) string {
	file, err := os.Open(path)

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func DownloadFile(url string, path string) {
	resp, err := http.Get(url)
	file, err := os.Create(path)
	_, err = io.Copy(file, resp.Body)
	err = resp.Body.Close()

	if err != nil {
		panic(err)
	}
}
