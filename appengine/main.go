package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/clipperhouse/jargon"

	"google.golang.org/appengine"
)

func main() {
	loadTemplates()
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
		// http.NotFound(w, r)
		// return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	model := textModel{}

	if r.Method == "POST" {
		text := r.PostFormValue("text")
		model.Original = text
		model.Result = jargon.StackExchange.Lemmatize(text)
	} else {
		model.Original = demo
	}

	var tmpl *template.Template

	if isAjax(r) {
		tmpl = _result
	} else {
		tmpl = layout
	}

	tmpl.Execute(w, model)
}

type textModel struct {
	Path     string
	Original string
	Result   string
}

var layout, _result *template.Template

func loadTemplates() {
	t, err := template.ParseFiles("layout.html", "_result.html")
	if err != nil {
		panic(err)
	}
	layout = t.Lookup("layout")
	_result = t.Lookup("_result")
}

func isAjax(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
