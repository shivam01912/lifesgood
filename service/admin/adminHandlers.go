package admin

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"lifesgood/db/mongo"
	"lifesgood/model"
	"lifesgood/service/util"
	"log"
	"net/http"
	"strings"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := getBasePageVars()

	t, _ := template.ParseFiles("./data/templates/admin_login_template.gohtml")

	t.ExecuteTemplate(w, "Login", vars)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func fetchAdminCredentials(username string) model.Credentials {
	var filter, option interface{}
	filter = bson.D{{"username", username}}
	option = bson.D{{"_id", 0}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	result := mongo.FindOne(client, ctx, "lifesgood", "credentials", filter, option)

	data, err := bson.Marshal(result)
	if err != nil {
		log.Println("Unable to marshal Credentials")
	}

	var credentials model.Credentials
	err = bson.Unmarshal(data, &credentials)
	if err != nil {
		log.Println("Unable to Unmarshal Credentials")
		log.Println(err)
	}

	return credentials
}

func ProcessLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	credentials := fetchAdminCredentials(r.PostForm.Get("username"))

	userPass := r.PostForm.Get("password")

	match := CheckPasswordHash(userPass, credentials.Password)

	if !match {
		w.Write([]byte("You are not an Admin"))
		return
	}

	vars := getBasePageVars()

	t, _ := template.ParseFiles("./data/templates/admin_page_template.gohtml")

	err = t.ExecuteTemplate(w, "Admin", vars)
	if err != nil {
		log.Println("Error in executing admin template : ", err)
		return
	}
}

func AddBlogHandler(w http.ResponseWriter, r *http.Request) {
	vars := getBasePageVars()

	t, _ := template.ParseFiles("./data/templates/add_blog_template.gohtml")

	err := t.ExecuteTemplate(w, "AddBlog", vars)
	if err != nil {
		log.Println("Error in executing home template : ", err)
		return
	}
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

	post := model.Blog{
		Title:     r.PostForm.Get("title"),
		Brief:     r.PostForm.Get("brief"),
		Tags:      tags,
		Content:   fileBytes,
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

	_, err = w.Write([]byte("Blog Added Successfully"))
	if err != nil {
		return
	}
}

func getBasePageVars() map[string]interface{} {
	result := map[string]interface{}{}
	util.AddHeader(result)
	util.AddFooter(result)
	return result
}
