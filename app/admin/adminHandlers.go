package admin

import (
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"lifesgood/app/util"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := util.ConstructBasePageVars()

	t, _ := template.ParseFiles("./data/templates/admin_login_template.gohtml")

	err := t.ExecuteTemplate(w, "Login", vars)
	if err != nil {
		log.Println("Error in executing Login template : ", err)
		return
	}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ProcessLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	credentials := util.FetchAdminCredentials(r.PostForm.Get("username"))

	userPass := r.PostForm.Get("password")

	match := CheckPasswordHash(userPass, credentials.Password)

	if !match {
		w.Write([]byte("You are not an Admin"))
		return
	}

	util.SetCookieHandler(w, r, credentials.Username, credentials.Password)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

	http.Redirect(w, r, "/admin/home", http.StatusSeeOther)
}

func AdminHome(w http.ResponseWriter, r *http.Request) {

	if !util.ValidateCookie(r, "session") {
		w.Write([]byte("You are not an admin"))
		return
	}

	vars := util.ConstructBasePageVars()

	t, _ := template.ParseFiles("./data/templates/admin_page_template.gohtml")

	err := t.ExecuteTemplate(w, "Admin", vars)
	if err != nil {
		log.Println("Error in executing admin template : ", err)
		return
	}
}
