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

	if r.Method == "OPTIONS" {
		return
	}

	jargonHandler(w, r)
}

var lemmatizer = jargon.NewLemmatizer(stackexchange.Dictionary, 3)

func jargonHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	route := parts[1]

	var tokens chan jargon.Token

	switch route {
	case "text":
		tokens = jargon.TokenizeHTML(r.Body)
		//		tokens = jargon.Tokenize(r.Body)
	case "html":
		tokens = jargon.TokenizeHTML(r.Body)
	default:
		http.NotFound(w, r)
		return
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
