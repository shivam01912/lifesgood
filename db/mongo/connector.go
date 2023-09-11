package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lifesgood/model"
	"log"
	"time"
)

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect() (*mongo.Client, context.Context, context.CancelFunc) {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://lifesgood:lifesgood@lifesgood.j9e1mix.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	return client, ctx, cancel
}

func InsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc model.Blog) (*mongo.InsertOneResult, error) {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.InsertOne(ctx, doc)

	return result, err
}

func FindAll(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) *mongo.Cursor {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.Find(ctx, query, options.Find().SetProjection(field).SetSort(bson.D{{"createdat", -1}}))
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func FindOne(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) bson.M {

	collection := client.Database(dataBase).Collection(col)

	var result bson.M
	collection.FindOne(ctx, query, options.FindOne().SetProjection(field)).Decode(&result)

	return result
}
