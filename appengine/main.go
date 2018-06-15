// A demo of jargon for use on Google App Engine
package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/clipperhouse/jargon/stackexchange"

	"github.com/clipperhouse/jargon"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", mainHandler)
	appengine.Main()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	jargonHandler(w, r)
}

var lemmatizer = jargon.NewLemmatizer(stackexchange.Dictionary, 3)

func jargonHandler(w http.ResponseWriter, r *http.Request) {
	text := r.PostFormValue("text")
	html := r.PostFormValue("html")

	var tokens chan jargon.Token
	if len(text) > 0 {
		r := strings.NewReader(text)
		tokens = jargon.Tokenize(r)
	} else {
		r := strings.NewReader(html)
		tokens = jargon.TokenizeHTML(r)
	}

	lemmatized := lemmatizer.Lemmatize(tokens)

	for t := range lemmatized {
		if t.IsLemma() {
			lemma.Execute(w, t)
		} else {
			plain.Execute(w, t)
		}
	}
}

var lemma = template.Must(template.New("lemma").Parse(`<span class="lemma">{{ . }}</span>`))
var plain = template.Must(template.New("plain").Parse(`{{ . }}`))

func isAjax(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
