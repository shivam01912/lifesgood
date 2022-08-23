package requestHandlers

import (
	"bytes"
	"html/template"
	"net/http"
	"log"
	"lifesgood/service/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lifesgood/db/mongo"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	homeVars := map[string]interface{}{}

	util.AddHeader(homeVars)
	util.AddFooter(homeVars)
	addCards(homeVars)

	t, _ := template.ParseFiles("../data/templates/home_template.html")
	
	t.ExecuteTemplate(w, "Home", homeVars)
}

func addCards(vars map[string]interface{}) {
	var buf bytes.Buffer
	var filter, option interface{}
	filter = bson.D{}
	option = bson.D{{"_id", 1}, {"title", 1}, {"brief", 1}}
	
	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	cursor := mongo.FindAll(client, ctx, "lifesgood", "blogs", filter, option)
	
	var results []bson.M

	if err := cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	
	for _, blog := range results {
		cardVars := map[string]interface{}{
			"Title": blog["title"],
			"Brief": blog["brief"],
			"Link": "/blog?id="+blog["_id"].(primitive.ObjectID).Hex(),
		}

		card, _ := template.ParseFiles("../data/templates/card_template.html")
		card.ExecuteTemplate(&buf, "Card", cardVars)
	}

	vars["Content"] = template.HTML(buf.String())
}
