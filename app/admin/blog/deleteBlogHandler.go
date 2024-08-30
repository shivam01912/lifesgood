package blog

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"lifesgood/app/config"
	"lifesgood/app/provider"
	"lifesgood/app/util"
	"lifesgood/db/mongo"
	"log"
	"net/http"
)

func DeletePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	homeVars := map[string]interface{}{}

	util.PopulateBasePageVars(homeVars)
	provider.AddCards(homeVars, config.DELETE)

	t, err := template.ParseFiles("./data/templates/home_template.gohtml")

	if err != nil {
		log.Println("Error parsing template : ", err)
	}

	err = t.ExecuteTemplate(w, "Home", homeVars)
	if err != nil {
		log.Println("Error in executing delete main page template : ", err)
		return
	}
}

func DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	//fetch blog likes
	id, ok := r.URL.Query()["id"]

	if !ok || len(id) < 1 {
		log.Fatal("Url Param 'id' is missing")
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id[0])
	if err != nil {
		log.Fatal("Invalid id")
		return
	}

	var filter interface{}
	filter = bson.D{{"_id", objectId}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	deleteOneResult, err := mongo.DeleteBlog(client, ctx, mongo.DBName, mongo.BlogCollection, filter)
	if err != nil {
		log.Println("Failed to update the likes for the blog with id : ", objectId, err)
		w.Write([]byte("Unable to delete blog."))
	}

	log.Println("Delete successful for :", deleteOneResult)

	w.Write([]byte("Blog deleted successfully."))
}
