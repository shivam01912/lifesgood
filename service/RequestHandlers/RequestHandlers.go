package requestHandlers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	md, _ := os.ReadFile("../data/sample.md")
	html := markdown.ToHTML(md, nil, nil)

	blogVars := map[string]interface{}{
		"Title":   "My First Blog",
		"Content": template.HTML(string(html)),
	}

	addHeader(blogVars)
	addFooter(blogVars)

	t, _ := template.ParseFiles("../data/templates/blog_template.html")

	t.ExecuteTemplate(w, "Blog", blogVars)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	blogVars := map[string]interface{}{
		"Title": "My First Blog",
	}

	addHeader(blogVars)
	addFooter(blogVars)
	addCards(blogVars)

	t, _ := template.ParseFiles("../data/templates/home_template.html")

	// w.Write(html)
	t.ExecuteTemplate(w, "Home", blogVars)
}

func addCards(vars map[string]interface{}) {
	var buf bytes.Buffer
	for i := 0; i < 5; i++ {
		card, _ := template.ParseFiles("../data/templates/card_template.html")
		card.ExecuteTemplate(&buf, "Card",nil)
	}

	vars["Content"] = template.HTML(buf.String())
}

func addHeader(vars map[string]interface{}) {
	navbar, _ := os.ReadFile("../data/templates/navbar.html")
	vars["Navbar"] = template.HTML(string(navbar))
}

func addFooter(vars map[string]interface{}) {
	footer, _ := os.ReadFile("../data/templates/footer.html")
	vars["Footer"] = template.HTML(string(footer))
}