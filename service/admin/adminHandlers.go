package admin

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"lifesgood/db/mongo"
	"lifesgood/model"
	"lifesgood/service/util"
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

	credentials := fetchAdminCredentials(r.PostForm.Get("username"))

	userPass := r.PostForm.Get("password")

	match := CheckPasswordHash(userPass, credentials.Password)

	if !match {
		w.Write([]byte("You are not an Admin"))
		return
	}

	vars := util.ConstructBasePageVars()

	t, _ := template.ParseFiles("./data/templates/admin_page_template.gohtml")

	err = t.ExecuteTemplate(w, "Admin", vars)
	if err != nil {
		log.Println("Error in executing admin template : ", err)
		return
	}
}

func fetchAdminCredentials(username string) model.Credentials {
	var filter, option interface{}
	filter = bson.D{{"username", username}}
	option = bson.D{{"_id", 0}}

	client, ctx, cancel := mongo.Connect()
	defer mongo.Close(client, ctx, cancel)

	result := mongo.FindOne(client, ctx, "lifesgood", "credentials", filter, option)

	data, err := bson.Marshal(result)
	if err != nil {
		log.Println("Unable to marshal Credentials")
	}

	var credentials model.Credentials
	err = bson.Unmarshal(data, &credentials)
	if err != nil {
		log.Println("Unable to Unmarshal Credentials")
		log.Println(err)
	}

	return credentials
}
