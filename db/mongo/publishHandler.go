package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"lifesgood/model"
)

func InsertBlog(client *mongo.Client, ctx context.Context, dataBase, col string, doc model.Blog) (*mongo.InsertOneResult, error) {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.InsertOne(ctx, doc)

	return result, err
}
