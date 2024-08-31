package blog

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"io"
	"lifesgood/app/util"
	"lifesgood/db/mongo"
	"lifesgood/model"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func UpdateBlogPageHandler(w http.ResponseWriter, r *http.Request) {
	if !util.ValidateCookie(w, r) {
		return
	}

	//fetch blog likes
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

	blog := fetchById(filter)

	tags := strings.Join(blog.Tags, ",")

	updateBlogVars := map[string]interface{}{
		"Link":  "/blog/update?id=" + objectId.Hex(),
		"Title": blog.Title,
		"Brief": blog.Brief,
		"Tags":  tags,
	}

	util.PopulateBasePageVars(updateBlogVars)

	t, _ := template.ParseFiles("./data/templates/update_blog_template.gohtml")

	err = t.ExecuteTemplate(w, "UpdateBlog", updateBlogVars)
	if err != nil {
		log.Println("Error in executing update blog template : ", err)
		return
	}
}

func ProcessUpdateBlog(w http.ResponseWriter, r *http.Request) {
	if !util.ValidateCookie(w, r) {
		return
	}

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

	err = r.ParseMultipartForm(10 << 10)
	if err != nil {
		log.Println(err)
	}

	file, _, err := r.FormFile("content")
	if err != nil {
		log.Println("Error Retrieving the File", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if errors.Is(err, http.ErrMissingFile) {
		log.Println(err)
	} else if err != nil {
		log.Println("Error reading the file : ", err)
		return
	}

	tags := strings.Split(r.PostForm.Get("tags"), ",")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}

	post := bson.D{{"$set", model.Blog{
		Title:        r.PostForm.Get("title"),
		Brief:        r.PostForm.Get("brief"),
		Tags:         tags,
		Content:      fileBytes,
		ModifiedDate: time.Now().Unix(),
	}}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	updateOneResult := mongo.UpdateBlog(client, ctx, mongo.DBName, mongo.BlogCollection, filter, post)

	if err != nil {
		w.Write([]byte("Unable to Publish blog"))
		w.Write([]byte(err.Error()))
		return
	}

	log.Println(updateOneResult)

	_, err = w.Write([]byte("Blog Updated Successfully"))
	if err != nil {
		return
	}
}

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

func fetchById(filter interface{}) *model.Blog {
	var option interface{}
	option = bson.D{{"_id", 0}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	result := mongo.FindOne(client, ctx, mongo.DBName, mongo.BlogCollection, filter, option)

	data, err := bson.Marshal(result)
	if err != nil {
		log.Println("Unable to marshal Blog")
	}

	var blog model.Blog
	err = bson.Unmarshal(data, &blog)
	if err != nil {
		log.Println("Unable to Unmarshal Blog")
	}

	return &blog
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
