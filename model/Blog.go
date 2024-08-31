package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Title        string             `bson:"title,omitempty"`
	Brief        string             `bson:"brief,omitempty"`
	Tags         []string           `bson:"tags,omitempty"`
	Content      []byte             `bson:"content,omitempty"`
	Views        int                `bson:"views,omitempty"`
	Likes        int                `bson:"likes,omitempty"`
	Comments     []string           `bson:"comments,omitempty"`
	CreatedAt    int64              `bson:"createdat,omitempty"`
	ModifiedDate int64              `bson:"modifieddate,omitempty"`
}
