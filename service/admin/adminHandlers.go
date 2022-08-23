package admin

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"time"
	"log"
	"strings"
	"lifesgood/service/util"
	"lifesgood/db/mongo"
	"lifesgood/model"
)

func AddBlogHandler(w http.ResponseWriter, r *http.Request) {
	vars := map[string]interface{}{}

	util.AddHeader(vars)
	util.AddFooter(vars)

	t, _ := template.ParseFiles("../data/templates/add_blog_template.html")

	t.ExecuteTemplate(w, "AddBlog", vars)
}

func ProcessAddBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 10)
	if err != nil {
		log.Fatal(err)
	}

	file, _, err := r.FormFile("content")
    if err != nil {
        log.Fatal("Error Retrieving the File", err)
        return
    }
    defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatal(err)
    }
    
	tags := strings.Split(r.PostForm.Get("tags"), ",")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}

	post := model.Blog {
		Title: r.PostForm.Get("title"),
		Brief: r.PostForm.Get("brief"),
		Tags: tags,
		Content: fileBytes,
		CreatedAt: time.Now().Unix(),
	}

	client, ctx, cancel := mongo.Connect()
    defer mongo.Close(client, ctx, cancel)

	insertOneResult := mongo.InsertOne(client, ctx, "lifesgood", "blogs", post)
 
	log.Println(insertOneResult)

	// fmt.Println(post)
	// fmt.Println(time.Unix(time.Now().Unix(), 0))
	
}