package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateBlog(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}, doc interface{}) (*mongo.UpdateResult, error) {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.UpdateOne(ctx, query, doc)

	return result, err
}
