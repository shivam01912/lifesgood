package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateBlog(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}, doc interface{}) *mongo.SingleResult {

	collection := client.Database(dataBase).Collection(col)

	result := collection.FindOneAndUpdate(ctx, query, doc, options.FindOneAndUpdate().SetReturnDocument(options.After))

	return result
}
