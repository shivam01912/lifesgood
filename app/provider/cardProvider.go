package provider

import (
	"bytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"lifesgood/app/config"
	"lifesgood/db/mongo"
	"log"
	"time"
)

func AddCards(vars map[string]interface{}, flow config.Flow) {
	var buf bytes.Buffer
	var filter, option interface{}
	filter = bson.D{}
	option = bson.D{{"_id", 1}, {"title", 1}, {"brief", 1}, {"tags", 1}, {"likes", 1}, {"createdat", 1}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	cursor := mongo.FindAll(client, ctx, mongo.DBName, mongo.BlogCollection, filter, option)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Println(err)
	}
	log.Println("Number of blogs fetched : ", len(results))

	for _, blog := range results {
		d, _ := blog["createdat"].(int64)
		date := time.Unix(d, 0).Format("2 Jan, 2006")
		cardVars := map[string]interface{}{
			"Title":      blog["title"],
			"Brief":      blog["brief"],
			"Tags":       blog["tags"],
			"Date":       date,
			"Likes":      blog["likes"],
			"Link":       "/blog?id=" + blog["_id"].(primitive.ObjectID).Hex(),
			"IsHomeFlow": isHomeFlow(flow),
			"UpdateLink": "/admin/updateblog?id=" + blog["_id"].(primitive.ObjectID).Hex(),
			"DeleteLink": "/blog/delete?id=" + blog["_id"].(primitive.ObjectID).Hex(),
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

func isHomeFlow(flow config.Flow) bool {
	if flow == config.HOME {
		return true
	}
	return false
}
