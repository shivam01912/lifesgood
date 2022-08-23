package requestHandlers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"lifesgood/service/util"
	"github.com/gomarkdown/markdown"
)

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	md, _ := os.ReadFile("../data/sample.md")
	html := markdown.ToHTML(md, nil, nil)

	blogVars := map[string]interface{}{
		"Title":   "My First Blog",
		"Content": template.HTML(string(html)),
	}

	util.AddHeader(blogVars)
	util.AddFooter(blogVars)

	t, _ := template.ParseFiles("../data/templates/blog_template.html")

	t.ExecuteTemplate(w, "Blog", blogVars)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	homeVars := map[string]interface{}{}

	util.AddHeader(homeVars)
	util.AddFooter(homeVars)
	addCards(homeVars)

	t, _ := template.ParseFiles("../data/templates/home_template.html")
	
	t.ExecuteTemplate(w, "Home", homeVars)
}

func addCards(vars map[string]interface{}) {
	var buf bytes.Buffer
	for i := 0; i < 5; i++ {
		card, _ := template.ParseFiles("../data/templates/card_template.html")
		card.ExecuteTemplate(&buf, "Card", nil)
	}

	vars["Content"] = template.HTML(buf.String())
}
