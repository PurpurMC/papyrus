package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Commit struct {
	Author      string `bson:"author,omitempty" json:"author"`
	Email       string `bson:"email,omitempty" json:"email"`
	Summary     string `bson:"summary,omitempty" json:"summary"`
	Description string `bson:"description,omitempty" json:"description"`
	Hash        string `bson:"hash,omitempty" json:"hash"`
	Timestamp   int64  `bson:"timestamp,omitempty" json:"timestamp"`
}

type File struct {
	Id           primitive.ObjectID `bson:"id,omitempty" json:"-"`
	InternalName string             `bson:"internal_name,omitempty" json:"-"`
	ContentType  string             `bson:"contentType,omitempty" json:"-"`
	Name         string             `bson:"name,omitempty" json:"name"`
	SHA512       string             `bson:"sha512,omitempty" json:"sha512"`
}
