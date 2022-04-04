package db

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/lucsky/cuid"
	"github.com/purpurmc/papyrus/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"os"
)

func CreateCollection(database *mongo.Database, collectionName string) {
	err := database.CreateCollection(context.TODO(), collectionName, nil)
	if err != nil {
		fmt.Println("Error creating collection: ", err)
		os.Exit(1)
	}
}

func InsertProject(database *mongo.Database, project types.Project) primitive.ObjectID {
	return Insert(database, "projects", project)
}

func InsertVersion(database *mongo.Database, version types.Version) primitive.ObjectID {
	return Insert(database, "versions", version)
}

func InsertBuild(database *mongo.Database, build types.Build) primitive.ObjectID {
	return Insert(database, "builds", build)
}

func Insert[T any](database *mongo.Database, collectionName string, object T) primitive.ObjectID {
	collection := database.Collection(collectionName)
	result, err := collection.InsertOne(context.TODO(), object)
	if err != nil {
		fmt.Println("Error inserting object: ", err)
		os.Exit(1)
	}
	return result.InsertedID.(primitive.ObjectID)
}

func GetProject(database *mongo.Database, filter *types.Project) *types.Project {
	return GetSingle(database, "projects", filter)
}

func GetVersion(database *mongo.Database, filter *types.Version) *types.Version {
	return GetSingle(database, "versions", filter)
}

func GetBuild(database *mongo.Database, filter *types.Build) *types.Build {
	return GetSingle(database, "builds", filter)
}

func GetSingle[T any](database *mongo.Database, collectionName string, filter *T) *T {
	collection := database.Collection(collectionName)
	var result *T
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		fmt.Println("Error getting project: ", err)
		os.Exit(1)
	}
	return result
}

func GetProjects(database *mongo.Database, filter *types.Project) []types.Project {
	return GetMultiple(database, "projects", filter)
}

func GetVersions(database *mongo.Database, filter *types.Version) []types.Version {
	return GetMultiple(database, "versions", filter)
}

func GetBuilds(database *mongo.Database, filter *types.Build) []types.Build {
	return GetMultiple(database, "builds", filter)
}

func GetMultiple[T any](database *mongo.Database, collectionName string, filter *T) []T {
	collection := database.Collection(collectionName)
	var cursor *mongo.Cursor
	var err error
	if filter == nil {
		cursor, err = collection.Find(context.TODO(), bson.D{})
	} else {
		cursor, err = collection.Find(context.TODO(), filter)
	}

	if err != nil {
		fmt.Println("Error getting objects: ", err)
		os.Exit(1)
	}

	var objects []T
	for cursor.Next(context.TODO()) {
		var object T
		if err := cursor.Decode(&object); err != nil {
			fmt.Println("Error decoding object: ", err)
			os.Exit(1)
		}
		objects = append(objects, object)
	}
	return objects
}

func UploadFile(bucket *gridfs.Bucket, data []byte) (string, string, string) {
	fileName := cuid.New()
	hash := sha512.Sum512(data)
	contentType := mimetype.Detect(data)

	_, err := bucket.UploadFromStream(fileName, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return fileName, hex.EncodeToString(hash[:]), contentType.String()
}

func DownloadFile(bucket *gridfs.Bucket, fileName string) []byte {
	var buffer bytes.Buffer
	_, err := bucket.DownloadToStreamByName(fileName, &buffer)
	if err != nil {
		fmt.Println("Error downloading file: ", err)
		os.Exit(1)
	}

	return buffer.Bytes()
}
