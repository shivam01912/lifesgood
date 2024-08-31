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
	go router.HandleFunc("/", requestHandler.HomePageHandler)
	go router.HandleFunc("/home", requestHandler.HomePageHandler)
	go router.HandleFunc("/blog", requestHandler.BlogHandler)

	//admin login handlers
	router.HandleFunc("/admin/login", admin.LoginHandler)
	router.HandleFunc("/admin/home", admin.AdminHome)
	router.HandleFunc("/admin", admin.ProcessLogin)

	//blog counter update handlers
	go router.HandleFunc("/blog/likes", requestHandler.LikesIncrement)
	go router.HandleFunc("/blog/views", requestHandler.ViewsIncrement)

	//create blog handlers
	router.HandleFunc("/admin/addblog", adminBlogHandler.AddBlogHandler)
	router.HandleFunc("/admin/blog/preview", adminBlogHandler.PreviewBlog)
	router.HandleFunc("/admin/blog/publish", adminBlogHandler.ProcessPublishBlog)

	//common admin handlers
	router.HandleFunc("/admin/modifyblog", adminBlogHandler.DeletePageHandler)

	//update blog handlers
	router.HandleFunc("/admin/updateblog", adminBlogHandler.UpdateBlogPageHandler)
	router.HandleFunc("/blog/update", adminBlogHandler.ProcessUpdateBlog)

	//delete blog handlers
	router.HandleFunc("/blog/delete", adminBlogHandler.DeleteBlogHandler)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Println("Unknown error : ", err)
		return
	}
}
