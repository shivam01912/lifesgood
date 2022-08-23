package mongo

import (
	"context"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lifesgood/model"
)

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc){
	 
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	return client, ctx, cancel
}

func InsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc model.Blog) *mongo.InsertOneResult {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
 