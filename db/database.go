package db

import (
	"context"
	"fmt"
	"github.com/purpurmc/papyrus/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	collection := database.Collection("projects")
	result, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		fmt.Println("Error inserting project: ", err)
		os.Exit(1)
	}
	return result.InsertedID.(primitive.ObjectID)
}

func InsertVersion(database *mongo.Database, version types.Version) primitive.ObjectID {
	collection := database.Collection("versions")
	result, err := collection.InsertOne(context.TODO(), version)
	if err != nil {
		fmt.Println("Error inserting version: ", err)
		os.Exit(1)
	}
	return result.InsertedID.(primitive.ObjectID)
}

func InsertBuild(database *mongo.Database, build types.Build) primitive.ObjectID {
	collection := database.Collection("builds")
	result, err := collection.InsertOne(context.TODO(), build)
	if err != nil {
		fmt.Println("Error inserting build: ", err)
		os.Exit(1)
	}
	return result.InsertedID.(primitive.ObjectID)
}

func GetProjects(database *mongo.Database, filter *types.Project) []types.Project {
	collection := database.Collection("projects")
	var cursor *mongo.Cursor
	var err error
	if filter == nil {
		cursor, err = collection.Find(context.TODO(), bson.D{})
	} else {
		cursor, err = collection.Find(context.TODO(), filter)
	}
	if err != nil {
		fmt.Println("Error getting projects: ", err)
		os.Exit(1)
	}
	var projects []types.Project
	for cursor.Next(context.TODO()) {
		var project types.Project
		err := cursor.Decode(&project)
		if err != nil {
			fmt.Println("Error decoding project: ", err)
			os.Exit(1)
		}
		projects = append(projects, project)
	}
	return projects
}

func GetVersions(database *mongo.Database, filter *types.Version) []types.Version {
	collection := database.Collection("versions")
	var cursor *mongo.Cursor
	var err error
	if filter == nil {
		cursor, err = collection.Find(context.TODO(), bson.D{})
	} else {
		cursor, err = collection.Find(context.TODO(), filter)
	}
	if err != nil {
		fmt.Println("Error getting versions: ", err)
		os.Exit(1)
	}
	var versions []types.Version
	for cursor.Next(context.TODO()) {
		var version types.Version
		err := cursor.Decode(&version)
		if err != nil {
			fmt.Println("Error decoding version: ", err)
			os.Exit(1)
		}
		versions = append(versions, version)
	}
	return versions
}

func GetBuilds(database *mongo.Database, filter *types.Build) []types.Build {
	collection := database.Collection("builds")
	var cursor *mongo.Cursor
	var err error
	if filter == nil {
		cursor, err = collection.Find(context.TODO(), bson.D{})
	} else {
		cursor, err = collection.Find(context.TODO(), filter)
	}
	if err != nil {
		fmt.Println("Error getting builds: ", err)
		os.Exit(1)
	}
	var builds []types.Build
	for cursor.Next(context.TODO()) {
		var build types.Build
		err := cursor.Decode(&build)
		if err != nil {
			fmt.Println("Error decoding build: ", err)
			os.Exit(1)
		}
		builds = append(builds, build)
	}
	return builds
}

func GetProject(database *mongo.Database, filter *types.Project) *types.Project {
	collection := database.Collection("projects")
	var project *types.Project
	err := collection.FindOne(context.TODO(), filter).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		fmt.Println("Error getting project: ", err)
		os.Exit(1)
	}
	return project
}

func GetVersion(database *mongo.Database, filter *types.Version) *types.Version {
	collection := database.Collection("versions")
	var version *types.Version
	err := collection.FindOne(context.TODO(), filter).Decode(&version)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		fmt.Println("Error getting version: ", err)
		os.Exit(1)
	}
	return version
}

func GetBuild(database *mongo.Database, filter *types.Build) *types.Build {
	collection := database.Collection("builds")
	var build *types.Build
	err := collection.FindOne(context.TODO(), filter).Decode(&build)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		fmt.Println("Error getting build: ", err)
		os.Exit(1)
	}
	return build
}
