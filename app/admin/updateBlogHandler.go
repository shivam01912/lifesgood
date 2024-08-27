package admin

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lifesgood/db/mongo"
	"lifesgood/model"
	"log"
	"net/http"
	"strconv"
)

func LikesIncrement(w http.ResponseWriter, r *http.Request) {

	//fetch blog likes
	id, ok := r.URL.Query()["id"]

	if !ok || len(id) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id[0])
	if err != nil {
		log.Println("Invalid id")
	}

	var filter, option interface{}
	filter = bson.D{{"_id", objectId}}
	option = bson.D{{"_id", 0}, {"likes", 1}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	result := mongo.FindOne(client, ctx, mongo.DBName, mongo.BlogCollection, filter, option)

	data, err := bson.Marshal(result)
	if err != nil {
		log.Println("Unable to marshal Blog")
	}

	var blog model.Blog
	err = bson.Unmarshal(data, &blog)
	if err != nil {
		log.Println("Unable to Unmarshal Blog")
	}

	//update the blog likes

	inc, ok := r.URL.Query()["inc"]

	doInc, _ := strconv.ParseBool(inc[0])

	newLikes := 0

	if doInc {
		newLikes = blog.Likes + 1
	} else {
		newLikes = blog.Likes - 1
	}

	post := bson.D{{"$set", bson.M{
		"likes": newLikes,
	}}}

	updateOneResult, err := mongo.UpdateBlog(client, ctx, mongo.DBName, mongo.BlogCollection, filter, post)
	if err != nil {
		log.Println("Failed to update the likes for the blog with id : ", objectId, err)
	}

	data, err = bson.Marshal(updateOneResult)
	if err != nil {
		log.Println("Unable to marshal Blog")
	}

	err = bson.Unmarshal(data, &blog)
	if err != nil {
		log.Println("Unable to Unmarshal Blog")
	}

	//return new likes count
	_, err = fmt.Fprintf(w, strconv.Itoa(newLikes))
	if err != nil {
		return
	}
}
