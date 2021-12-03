package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Data struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	Name string        `json:"name"`
	Versions []Version `json:"versions"`
}

type Version struct {
	Name string `json:"name"`
	Latest Build `json:"latest"`
	Builds []Build `json:"builds"`
}

type Build struct {
	Project string `json:"project"`
	Version string `json:"version"`
	Build string `json:"build"`
	Result string `json:"result"`
	Duration int     `json:"duration"`
	Commits []Commit `json:"commits"`
	Timestamp int    `json:"timestamp"`
	MD5 string `json:"md5"`
	Extension string `json:"extension"`
}

type Commit struct {
	Author string `json:"author"`
	Title string `json:"title"`
	Comment string `json:"comment"`
	Hash string `json:"hash"`
	Email string `json:"email"`
	Timestamp int `json:"timestamp"`
}

func GetData() Data {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/data.json", GetConfig().StoragePath))
	data := Data{}
	err = json.Unmarshal(file, &data)

	if err != nil {
		panic(err)
	}
	return data
}

func SaveData(data Data) {
	jsonBytes, err := json.Marshal(data)
	err = ioutil.WriteFile(fmt.Sprintf("%s/data.json", GetConfig().StoragePath), jsonBytes, os.ModePerm)

	if err != nil {
		panic(err)
	}
}
