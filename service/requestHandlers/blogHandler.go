package requestHandlers

import (
	"html/template"
	"net/http"
	"log"
	"time"
	"lifesgood/service/util"
	"github.com/gomarkdown/markdown"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lifesgood/db/mongo"
	"lifesgood/model"
)

func BlogHandler(w http.ResponseWriter, r *http.Request) {

	id, ok := r.URL.Query()["id"]
    
    if !ok || len(id) < 1 {
        log.Println("Url Param 'id' is missing")
        return
    }

	objectId, err := primitive.ObjectIDFromHex(id[0])
	if err != nil{
		log.Println("Invalid id")
	}

	blog := fetchById(objectId)
	
	html := markdown.ToHTML(blog.Content, nil, nil) 

	date := time.Unix(blog.CreatedAt, 0).Format("2 Jan, 2006")

	blogVars := map[string]interface{}{
		"Title": blog.Title,
		"Content": template.HTML(string(html)),
		"Tags": blog.Tags,
		"Date": date,
	}

	util.AddHeader(blogVars)
	util.AddFooter(blogVars)

	t, _ := template.ParseFiles("../data/templates/blog_template.html")

	t.ExecuteTemplate(w, "Blog", blogVars)
}

func fetchById(objectId primitive.ObjectID) model.Blog{
	var filter, option interface{}
	filter = bson.D{{"_id", objectId}}
	option = bson.D{{"_id", 0}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	result := mongo.FindOne(client, ctx, "lifesgood", "blogs", filter, option)

	data, err := bson.Marshal(result)
    if err != nil {
        log.Println("Unable to marshal Blog")
    }

	var blog model.Blog
    err = bson.Unmarshal(data, &blog)
	if err != nil{
		log.Println("Unable to Unmarshal Blog")
	}

	return blog
}