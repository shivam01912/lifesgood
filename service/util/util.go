package util

import (
	"html/template"
	"os"
)

func AddHeader(vars map[string]interface{}) {
	navbar, _ := os.ReadFile("../data/templates/navbar.html")
	vars["Navbar"] = template.HTML(string(navbar))
}

func AddFooter(vars map[string]interface{}) {
	footer, _ := os.ReadFile("../data/templates/footer.html")
	vars["Footer"] = template.HTML(string(footer))
}