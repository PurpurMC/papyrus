package types

type Commit struct {
	Author      string `bson:"author,omitempty" json:"author"`
	Email       string `bson:"email,omitempty" json:"email"`
	Summary     string `bson:"summary,omitempty" json:"summary"`
	Description string `bson:"description,omitempty" json:"description"`
	Hash        string `bson:"hash,omitempty" json:"hash"`
	Timestamp   int64  `bson:"timestamp,omitempty" json:"timestamp"`
}

type File struct {
	Id          string `bson:"id,omitempty" json:"-"`
	ContentType string `bson:"contentType,omitempty" json:"-"`
	Name        string `bson:"name,omitempty" json:"name"`
	SHA512      string `bson:"sha512,omitempty" json:"sha512"`
}
