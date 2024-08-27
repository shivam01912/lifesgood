package main

import (
	"github.com/joho/godotenv"
	"lifesgood/app/admin"
	"lifesgood/app/requestHandler"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fs := http.FileServer(http.Dir("./data/templates"))
	http.Handle("/css/", fs)

	http.HandleFunc("/", requestHandler.HomePageHandler)
	http.HandleFunc("/home", requestHandler.HomePageHandler)

	http.HandleFunc("/blog", requestHandler.BlogHandler)

	http.HandleFunc("/blog/likes", admin.LikesIncrement)

	http.HandleFunc("/admin/login", admin.LoginHandler)
	http.HandleFunc("/admin/home", admin.AdminHome)
	http.HandleFunc("/admin", admin.ProcessLogin)

	http.HandleFunc("/admin/addblog", admin.AddBlogHandler)
	http.HandleFunc("/admin/addblog/publish", admin.ProcessPublishBlog)

	//http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("../data"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
