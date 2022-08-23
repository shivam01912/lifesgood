package main

import (
	"lifesgood/service/requestHandlers"
	"lifesgood/service/admin"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", requestHandlers.HomePageHandler)
	http.HandleFunc("/blog", requestHandlers.BlogHandler)

	http.HandleFunc("/admin/login", admin.LoginHandler)
	http.HandleFunc("/admin", admin.ProcessLogin)

	http.HandleFunc("/admin/addblog", admin.AddBlogHandler)
	http.HandleFunc("/admin/addblog/publish", admin.ProcessPublishBlog)

	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("../data"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
