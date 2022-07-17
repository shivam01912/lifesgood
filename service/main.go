package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

func handler(w http.ResponseWriter, r *http.Request) {
	navbar, _ := os.ReadFile("../data/navbar.html")
	md, _ := os.ReadFile("../data/sample.md")
	html := markdown.ToHTML(md, nil, nil)
	footer, _ := os.ReadFile("../data/footer.html")

	blogVars := map[string]interface{}{
		"Header":  template.HTML(string(navbar)),
		"Title":   "My First Blog",
		"Content": template.HTML(string(html)),
		"Footer":  template.HTML(string(footer)),
	}

	t, _ := template.ParseFiles("../data/blog_template.html")

	t.ExecuteTemplate(w, "Blog", blogVars)
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("../data"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
