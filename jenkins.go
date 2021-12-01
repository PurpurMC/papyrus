package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type JenkinsData struct {
	Result string `json:"result"`
	Duration int `json:"duration"`
	Timestamp int `json:"timestamp"`
}

func getJenkinsData(url string, project string, build int) JenkinsData {
	response, err := http.Get(url + "/job/" + project + "/" + strconv.Itoa(build) + "/api/json")
	checkError(err)

	responseData, err := ioutil.ReadAll(response.Body)
	checkError(err)

	var responseObject JenkinsData
	checkError(json.Unmarshal(responseData, &responseObject))

	return responseObject
}
