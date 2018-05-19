package main

import (
	"html/template"
	"net/http"

	"github.com/clipperhouse/jargon"

	"google.golang.org/appengine"
)

func main() {
	getTemplate()
	http.HandleFunc("/", mainHandler)
	appengine.Main()
}

var demo = `We might have some prose regarding Ruby on Rails or NODEJS, like a job description.

Or some structured data:
{
    "language": "ObjC"
}`

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	model := textModel{
		Path: r.URL.Path,
	}

	if r.Method == "POST" {
		text := r.PostFormValue("text")
		model.Original = text
		model.Result = jargon.Lemmatize(text)
	} else {
		model.Original = demo
	}

	tmpl := getTemplate()
	tmpl.Execute(w, model)
}

type textModel struct {
	Path     string
	Original string
	Result   string
}

var mainTmpl *template.Template

func getTemplate() *template.Template {
	if mainTmpl == nil {
		t, err := template.ParseFiles("layout.html", "_result.html")
		if err != nil {
			panic(err)
		}
		mainTmpl = t
	}

	return mainTmpl
}
