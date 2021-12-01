package shared

import (
	"encoding/json"
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
	Build int `json:"build"`
	Result string `json:"result"`
	Duration int     `json:"duration"`
	Commits []Commit `json:"commits"`
	Timestamp int    `json:"timestamp"`
	MD5 string `json:"md5"`
	Extension string `json:"extension"`
}

type Commit struct {
	Author string `json:"author"`
	Description string `json:"description"`
	Hash string `json:"hash"`
	Email string `json:"email"`
	Timestamp int `json:"timestamp"`
}

func GetData() Data {
	file, err := ioutil.ReadFile("data.json")
	data := Data{}
	err = json.Unmarshal(file, &data)

	if err != nil {
		panic(err)
	}
	return data
}

func SaveData(data Data) {
	jsonBytes, err := json.Marshal(data)
	err = ioutil.WriteFile("data.json", jsonBytes, os.ModePerm) // todo: location

	if err != nil {
		panic(err)
	}
}
