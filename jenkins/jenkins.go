package jenkins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetJenkinsData(url string, project string, build string) JenkinsData {
	response, err := http.Get(fmt.Sprintf("%s/job/%s/%s/api/json", url, project, build))
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	var data JenkinsData
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	return data
}
