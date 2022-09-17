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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := getBasePageVars()

	t, _ := template.ParseFiles("../data/templates/admin_login_template.html")
	
	t.ExecuteTemplate(w, "Login", vars)
}

func ProcessLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	if r.PostForm.Get("username") != "admin" || r.PostForm.Get("password") != "admin" {
		w.Write([]byte("You are not an Admin"))
		return
	}

	vars := getBasePageVars()

	t, _ := template.ParseFiles("../data/templates/admin_page_template.html")

	t.ExecuteTemplate(w, "Admin", vars)
}

func AddBlogHandler(w http.ResponseWriter, r *http.Request) {
	vars := getBasePageVars()

	t, _ := template.ParseFiles("../data/templates/add_blog_template.html")

	t.ExecuteTemplate(w, "AddBlog", vars)
}

func ProcessPublishBlog(w http.ResponseWriter, r *http.Request) {
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

	insertOneResult, err := mongo.InsertOne(client, ctx, "lifesgood", "blogs", post)
 
	if err != nil {
		w.Write([]byte("Unable to Publish blog"))
		w.Write([]byte(err.Error()))
		return
	}

	log.Println(insertOneResult)

	w.Write([]byte("Blog Added Successfully"))
}

func getBasePageVars() map[string]interface{} {
	result := map[string]interface{}{}
	util.AddHeader(result)
	util.AddFooter(result)
	return result
}