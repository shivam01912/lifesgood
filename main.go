package main

import (
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./data/templates"))
	router.PathPrefix("/css/").Handler(fs)

	//direct REST API handler
	router.HandleFunc("/", requestHandler.HomePageHandler)
	router.HandleFunc("/home", requestHandler.HomePageHandler)
	router.HandleFunc("/blog", requestHandler.BlogHandler)

	//admin login handlers
	router.HandleFunc("/admin/login", admin.LoginHandler)
	router.HandleFunc("/admin/home", admin.AdminHome)
	router.HandleFunc("/admin", admin.ProcessLogin)

	//create blog handlers
	router.HandleFunc("/admin/addblog", adminBlogHandler.AddBlogHandler)
	router.HandleFunc("/admin/addblog/publish", adminBlogHandler.ProcessPublishBlog)

	//update blog handlers
	router.HandleFunc("/blog/likes", adminBlogHandler.LikesIncrement)

	//delete blog handlers
	router.HandleFunc("/admin/deleteblog", adminBlogHandler.DeletePageHandler)
	router.HandleFunc("/blog/delete", adminBlogHandler.DeleteBlogHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
