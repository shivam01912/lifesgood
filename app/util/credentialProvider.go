package util

import (
	"go.mongodb.org/mongo-driver/bson"
	"lifesgood/db/mongo"
	"lifesgood/model"
	"log"
)

func FetchAdminCredentials(username string) model.Credentials {
	var filter, option interface{}
	filter = bson.D{{"username", username}}
	option = bson.D{{"_id", 0}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	result := mongo.FindOne(client, ctx, mongo.DBName, mongo.CredentialCollection, filter, option)

	data, err := bson.Marshal(result)
	if err != nil {
		log.Println("Unable to marshal Credentials")
	}

	var credentials model.Credentials
	err = bson.Unmarshal(data, &credentials)
	if err != nil {
		log.Println("Unable to Unmarshal Credentials")
		log.Println(err)
	}

	return credentials
}
