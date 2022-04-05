package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/purpurmc/papyrus/types"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

func GetJenkinsData(url string, project string, build string) *types.JenkinsData {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/job/%s/%s/api/json", strings.TrimSuffix(url, "/"), project, build), nil)
	request.Header.Add("cf-access-token", viper.GetString("utils.cloudflare-access-token"))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if response.StatusCode != 200 {
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	var data types.JenkinsData
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	return &data
}

func DownloadJenkinsWorkspaceFile(url string, project string, path string) []byte {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/job/%s/ws/%s", strings.TrimSuffix(url, "/"), project, path), nil)
	request.Header.Add("cf-access-token", viper.GetString("utils.cloudflare-access-token"))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if response.StatusCode != 200 {
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, response.Body)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}