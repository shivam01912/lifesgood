package main

import (
	"lifesgood/service/requestHandlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", requestHandlers.HomePageHandler)
	http.HandleFunc("/sample", requestHandlers.BlogHandler)

	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("../data"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
