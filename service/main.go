package main

import (
	"lifesgood/service/admin"
	"lifesgood/service/requestHandlers"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("../data/templates"))
	http.Handle("/css/", fs)

	http.HandleFunc("/home", requestHandlers.HomePageHandler)
	http.HandleFunc("/blog", requestHandlers.BlogHandler)

	http.HandleFunc("/admin/login", admin.LoginHandler)
	http.HandleFunc("/admin", admin.ProcessLogin)

	http.HandleFunc("/admin/addblog", admin.AddBlogHandler)
	http.HandleFunc("/admin/addblog/publish", admin.ProcessPublishBlog)

	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("../data"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
