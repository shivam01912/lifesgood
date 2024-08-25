package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteBlog(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}) (*mongo.DeleteResult, error) {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.DeleteOne(ctx, query)

	return result, err
}
