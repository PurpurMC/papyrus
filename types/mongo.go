package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	CreatedAt int64              `bson:"created_at,omitempty"`
}

type Version struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ProjectId primitive.ObjectID `bson:"project_id,omitempty"`
	CreatedAt int64              `bson:"created_at,omitempty"`
	Name      string             `bson:"name,omitempty"`
}

type Build struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	VersionId primitive.ObjectID `bson:"version_id,omitempty"`
	CreatedAt int64              `bson:"created_at,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Result    string             `bson:"result,omitempty"`
	Commits   []Commit           `bson:"commits,omitempty"`
	Files     []File             `bson:"files,omitempty"`
}

type FileMetadata struct {
	ContentType string `bson:"content_type,omitempty"`
}
