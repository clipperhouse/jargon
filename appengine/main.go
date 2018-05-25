package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/clipperhouse/jargon"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/jargon", mainHandler)
	appengine.Main()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/jargon" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")

	text := r.PostFormValue("text")
	original := jargon.TechProse.Tokenize(text)
	lemmatized := jargon.StackExchange.LemmatizeTokens(original)

	for _, t := range lemmatized {
		if t.IsLemma() {
			span.Execute(w, t)
		} else {
			w.Write([]byte(t.String()))
		}
	}
}

var span = template.Must(template.New("span").Parse(`<span class="lemma">{{ .String }}</span>`))

func isAjax(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
