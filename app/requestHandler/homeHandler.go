package requestHandler

import (
	"html/template"
	"lifesgood/app/config"
	"lifesgood/app/provider"
	"lifesgood/app/util"
	"log"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	homeVars := map[string]interface{}{}

	util.PopulateBasePageVars(homeVars)
	provider.AddCards(homeVars, config.HOME) // TODO: Fix repeated calls to this function every time blog refreshes

	t, err := template.ParseFiles("./data/templates/home_template.gohtml")

	if err != nil {
		log.Println("Error parsing template : ", err)
	}

	err = t.ExecuteTemplate(w, "Home", homeVars)
	if err != nil {
		log.Println("Error in executing home template : ", err)
		return
	}
}
