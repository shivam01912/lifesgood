package blog

import (
	"github.com/gomarkdown/markdown"
	"html/template"
	"io"
	"lifesgood/app/util"
	"log"
	"net/http"
	"strings"
	"time"
)

func PreviewBlog(w http.ResponseWriter, r *http.Request) {
	if !util.ValidateCookie(w, r) {
		return
	}

	file, _, err := r.FormFile("content")
	if err != nil {
		log.Println("Error Retrieving the File", err)
		return
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	html := markdown.ToHTML(fileBytes, nil, nil)

	tags := strings.Split(r.PostForm.Get("tags"), ",")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}

	previewVars := map[string]interface{}{
		"Link":          "/blog/likes?id=dummy",
		"Title":         r.PostForm.Get("title"),
		"Content":       template.HTML(html),
		"Tags":          tags,
		"Date":          time.Unix(time.Now().Unix(), 0).Format("2 Jan, 2006"),
		"Likes":         0,
		"IsPreviewFlow": true,
	}

	util.PopulateBasePageVars(previewVars)

	t, _ := template.ParseFiles("./data/templates/blog_template.gohtml")

	err = t.ExecuteTemplate(w, "Blog", previewVars)
	if err != nil {
		log.Println("Error in executing blog template : ", err)
		return
	}

}
