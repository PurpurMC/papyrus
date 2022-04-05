package v1

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/purpurmc/papyrus/db"
	"github.com/purpurmc/papyrus/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func MigrateV1(url string, defaultFilename string) {
	database, bucket := db.NewMongo()

	response, err := http.Get(url + "/v2")
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	var projectsResponse types.ProjectsResponse
	err = json.NewDecoder(response.Body).Decode(&projectsResponse)
	if err != nil {
		panic(err)
	}

	for _, project := range projectsResponse.Projects {
		projectId := db.InsertProject(database, types.Project{
			Name:      project,
			CreatedAt: time.Now().Unix(),
		})

		response, err := http.Get(url + "/v2/" + project)
		if err != nil {
			panic(err)
		}

		var projectResponse types.ProjectResponse
		err = json.NewDecoder(response.Body).Decode(&projectResponse)
		if err != nil {
			panic(err)
		}

		for _, version := range projectResponse.Versions {
			versionId := db.InsertVersion(database, types.Version{
				ProjectId: projectId,
				CreatedAt: time.Now().Unix(),
				Name:      version,
			})

			response, err := http.Get(url + "/v2/" + project + "/" + version)

			var versionResponse types.VersionResponse
			err = json.NewDecoder(response.Body).Decode(&versionResponse)
			if err != nil {
				panic(err)
			}

			for _, build := range versionResponse.Builds.All {
				response, err := http.Get(url + "/v2/" + project + "/" + version + "/" + build)

				var buildResponse LegacyBuildResponse
				err = json.NewDecoder(response.Body).Decode(&buildResponse)
				if err != nil {
					panic(err)
				}

				fileResponse, err := http.Get(url + "/v2/" + project + "/" + version + "/" + build + "/download")
				if err != nil {
					panic(err)
				}

				// parse fileResponse as a byte array
				var file []byte
				file, err = ioutil.ReadAll(fileResponse.Body)
				if err != nil {
					panic(err)
				}

				fileId, fileName, hash, contentType := db.UploadFile(bucket, file)

				commits := make([]types.Commit, 0)
				for _, commit := range buildResponse.Commits {
					splits := strings.Split(commit.Description, "\n")
					var summary string
					if len(splits) > 0 {
						summary = splits[0]
					} else {
						summary = commit.Description
					}

					commits = append(commits, types.Commit{
						Author:      commit.Author,
						Email:       commit.Email,
						Summary:     summary,
						Description: commit.Description,
						Hash:        commit.Hash,
						Timestamp:   commit.Timestamp,
					})
				}

				buildId := db.InsertBuild(database, types.Build{
					VersionId: versionId,
					CreatedAt: buildResponse.Timestamp,
					Name:      build,
					Result:    buildResponse.Result,
					Flags:     make([]string, 0),
					Commits:   commits,
					Files: []types.File{
						{
							Id:           fileId,
							InternalName: fileName,
							ContentType:  contentType,
							Name:         defaultFilename,
							SHA512:       hash,
						},
					},
				})

				md5hash := md5.Sum(file)
				v1 := database.Collection("v1")
				_, err = v1.InsertOne(context.TODO(), LegacyBuildData{
					BuildId: buildId,
					MD5:     hex.EncodeToString(md5hash[:]),
				})
				if err != nil {
					panic(err)
				}

				fmt.Println("Migrated " + project + "/" + version + "/" + build)

				response.Body.Close()
			}

			response.Body.Close()
		}

		response.Body.Close()
	}
}

type LegacyBuildData struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	BuildId primitive.ObjectID `bson:"build_id,omitempty"`
	MD5     string             `bson:"md5,omitempty"`
}

type LegacyBuildResponse struct {
	Project string `json:"project"`
	Version string `json:"version"`
	Name    string `json:"build"`
	Result  string `json:"result"`
	Commits []struct {
		Author      string `json:"author"`
		Description string `json:"description"`
		Hash        string `json:"hash"`
		Email       string `json:"email"`
		Timestamp   int64  `json:"timestamp"`
	} `json:"commits"`
	Timestamp int64  `json:"timestamp"`
	MD5       string `json:"md5"`
}
