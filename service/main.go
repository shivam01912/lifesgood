package main

import (
	"lifesgood/service/requestHandlers"
	"lifesgood/service/admin"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", requestHandlers.HomePageHandler)
	http.HandleFunc("/sample", requestHandlers.BlogHandler)
	http.HandleFunc("/admin", admin.AddBlogHandler)
	http.HandleFunc("/admin/addblog", admin.ProcessAddBlog)

	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("../data"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
