package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func FindAll(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}, field interface{}) *mongo.Cursor {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.Find(ctx, query, options.Find().SetProjection(field).SetSort(bson.D{{"createdat", -1}}))
	if err != nil {
		log.Println(err)
	}

	return result
}

func FindOne(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}, field interface{}) bson.M {

	collection := client.Database(dataBase).Collection(col)

	var result bson.M
	collection.FindOne(ctx, query, options.FindOne().SetProjection(field)).Decode(&result)

	return result
}
