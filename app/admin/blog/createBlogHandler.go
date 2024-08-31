package blog

import (
	"html/template"
	"io"
	"lifesgood/app/config"
	"lifesgood/app/util"
	"lifesgood/db/mongo"
	"lifesgood/model"
	"log"
	"net/http"
	"strings"
	"time"
)

func AddBlogHandler(w http.ResponseWriter, r *http.Request) {
	if !util.ValidateCookie(w, r) {
		return
	}

	vars := util.ConstructBasePageVars()

	t, _ := template.ParseFiles("./data/templates/add_blog_template.gohtml")

	err := t.ExecuteTemplate(w, "AddBlog", vars)
	if err != nil {
		log.Println("Error in executing home template : ", err)
		return
	}
}

func ProcessPublishBlog(w http.ResponseWriter, r *http.Request) {
	if !util.ValidateCookie(w, r) {
		return
	}

	err := r.ParseMultipartForm(10 << 10)
	if err != nil {
		log.Println(err)
	}

	file, _, err := r.FormFile("content")
	if err != nil {
		log.Println("Error Retrieving the File", err)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	tags := strings.Split(r.PostForm.Get("tags"), ",")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}

	post := model.Blog{
		Title:        r.PostForm.Get("title"),
		Brief:        r.PostForm.Get("brief"),
		Tags:         tags,
		Content:      fileBytes,
		CreatedAt:    time.Now().Unix(),
		ModifiedDate: time.Now().Unix(),
	}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	insertOneResult, err := mongo.InsertBlog(client, ctx, mongo.DBName, mongo.BlogCollection, post)

	if err != nil {
		w.Write([]byte("Unable to Publish blog"))
		w.Write([]byte(err.Error()))
		return
	}

	log.Println(insertOneResult)

	_, err = w.Write([]byte("Blog Added Successfully"))
	if err != nil {
		return
	}
}

func isPreviewFlow(flow config.Flow) bool {
	if flow == config.CREATE {
		return true
	}
	return false
}
