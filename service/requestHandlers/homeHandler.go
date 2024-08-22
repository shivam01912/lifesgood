package requestHandlers

import (
	"bytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"lifesgood/db/mongo"
	"lifesgood/service/util"
	"log"
	"net/http"
	"time"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	homeVars := map[string]interface{}{}

	util.AddHeader(homeVars)
	util.AddFooter(homeVars)
	addCards(homeVars)

	t, err := template.ParseFiles("./data/templates/home_template.gohtml")

	if err != nil {
		log.Println("Error parsing template : ", err)
	}

	err = t.ExecuteTemplate(w, "Home", homeVars)
	if err != nil {
		log.Println("Error in executing home template : ", err)
		return
	}
}

func addCards(vars map[string]interface{}) {
	var buf bytes.Buffer
	var filter, option interface{}
	filter = bson.D{}
	option = bson.D{{"_id", 1}, {"title", 1}, {"brief", 1}, {"tags", 1}, {"likes", 1}, {"createdat", 1}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	cursor := mongo.FindAll(client, ctx, "lifesgood", "blogs", filter, option)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	log.Println("Number of blogs fetched : ", len(results))

	for _, blog := range results {
		d, _ := blog["createdat"].(int64)
		date := time.Unix(d, 0).Format("2 Jan, 2006")
		cardVars := map[string]interface{}{
			"Title": blog["title"],
			"Brief": blog["brief"],
			"Tags":  blog["tags"],
			"Date":  date,
			"Likes": blog["likes"],
			"Link":  "/blog?id=" + blog["_id"].(primitive.ObjectID).Hex(),
		}

		card, err := template.ParseFiles("./data/templates/card_template.gohtml")

		if err != nil {
			log.Println(err)
		}

		err = card.ExecuteTemplate(&buf, "Card", cardVars)
		if err != nil {
			log.Println("Error in executing card template : ", err)
			return
		}
	}

	vars["Content"] = template.HTML(buf.String())
}
