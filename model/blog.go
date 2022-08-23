package model

type Blog struct {
	// Title string `bson:"title,omitempty"`
	// Brief string `bson:"brief,omitempty"`
	// Tags []string `bson:"tags,omitempty"`
	// Content []byte `bson:"content,omitempty"`
	// Views int `bson:"views,omitempty"`
	// Likes int `bson:"likes,omitempty"`
	// Comments []string `bson:"comments,omitempty"`
	Title string
	Brief string
	Tags []string
	Content []byte
	Views int `default:0`
	Likes int `default:0`
	Comments []string
	CreatedAt int64
}