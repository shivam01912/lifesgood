package util

import (
	"html/template"
	"os"
)

func ConstructBasePageVars() map[string]interface{} {
	basePageVars := map[string]interface{}{}
	addHeader(basePageVars)
	addFooter(basePageVars)
	return basePageVars
}

func PopulateBasePageVars(vars map[string]interface{}) {
	addHeader(vars)
	addFooter(vars)
}

func addHeader(vars map[string]interface{}) {
	navbar, _ := os.ReadFile("./data/templates/navbar.gohtml")
	vars["Navbar"] = template.HTML(string(navbar))
}

func addFooter(vars map[string]interface{}) {
	footer, _ := os.ReadFile("./data/templates/footer.gohtml")
	vars["Footer"] = template.HTML(string(footer))
}
