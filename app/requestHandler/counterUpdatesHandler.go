package requestHandler

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

func ViewsIncrement(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["id"]

	if !ok || len(id) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id[0])
	if err != nil {
		log.Println("Invalid id")
	}

	var filter interface{}
	filter = bson.D{{"_id", objectId}}

	post := bson.D{{"$inc", bson.M{
		"views": 1,
	}}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	var blog model.Blog
	mongo.UpdateBlog(client, ctx, mongo.DBName, mongo.BlogCollection, filter, post).Decode(&blog)
	if err != nil {
		log.Println("Failed to update the likes for the blog with id : ", objectId, err)
	}

	//return new views count
	_, err = fmt.Fprintf(w, strconv.Itoa(blog.Views))
	if err != nil {
		return
	}
}

func LikesIncrement(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["id"]

	if !ok || len(id) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id[0])
	if err != nil {
		log.Println("Invalid id")
	}

	var filter interface{}
	filter = bson.D{{"_id", objectId}}

	//update the blog likes
	newLikes := updateBlogLikes(r, filter)

	//return new likes count
	_, err = fmt.Fprintf(w, strconv.Itoa(newLikes))
	if err != nil {
		return
	}
}

func updateBlogLikes(r *http.Request, filter interface{}) int {
	inc, ok := r.URL.Query()["inc"]
	if !ok {
		log.Println("Url Param 'inc' is missing")
		return 0
	}

	doInc, _ := strconv.ParseBool(inc[0])

	delta := 0

	if doInc {
		delta = 1
	} else {
		delta = -1
	}

	post := bson.D{{"$inc", bson.M{
		"likes": delta,
	}}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	var blog model.Blog
	mongo.UpdateBlog(client, ctx, mongo.DBName, mongo.BlogCollection, filter, post).Decode(&blog)

	return blog.Likes
}
