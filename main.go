package main

import (
	"github.com/joho/godotenv"
	"lifesgood/app/admin"
	adminBlogHandler "lifesgood/app/admin/blog"
	"lifesgood/app/requestHandler"
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	fs := http.FileServer(http.Dir("./data/templates"))
	http.Handle("/css/", fs)

	//direct REST API handler
	http.HandleFunc("/", requestHandler.HomePageHandler)
	http.HandleFunc("/home", requestHandler.HomePageHandler)
	http.HandleFunc("/blog", requestHandler.BlogHandler)

	//admin login handlers
	http.HandleFunc("/admin/login", admin.LoginHandler)
	http.HandleFunc("/admin/home", admin.AdminHome)
	http.HandleFunc("/admin", admin.ProcessLogin)

	//create blog handlers
	http.HandleFunc("/admin/addblog", adminBlogHandler.AddBlogHandler)
	http.HandleFunc("/admin/addblog/publish", adminBlogHandler.ProcessPublishBlog)

	//update blog handlers
	http.HandleFunc("/blog/likes", adminBlogHandler.LikesIncrement)

	//delete blog handlers
	http.HandleFunc("/admin/deleteblog", adminBlogHandler.DeletePageHandler)
	http.HandleFunc("/blog/delete", adminBlogHandler.DeleteBlogHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
