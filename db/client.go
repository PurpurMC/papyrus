package db

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo() (*mongo.Database, *gridfs.Bucket) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.port"))))
	if err != nil {
		panic(err)
	}

	database := client.Database(viper.GetString("db.db"))
	bucket, err := gridfs.NewBucket(database)
	if err != nil {
		panic(err)
	}

	return database, bucket
}
