// A demo of jargon for use on Google App Engine
package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/clipperhouse/jargon/stackexchange"

	"github.com/clipperhouse/jargon"
)

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/_ah/health", healthCheckHandler)

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	cors(w)

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

	var tokens jargon.Tokens

	switch route {
	case "text":
		tokens = jargon.Tokenize(r.Body)
	case "html":
		tokens = jargon.TokenizeHTML(r.Body)
	default:
		http.NotFound(w, r)
		return
	}

	lemmatized := lemmatizer.Lemmatize(tokens)

	for {
		t := lemmatized.Next()
		if t == nil {
			break
		}
		if t.IsLemma() {
			lemma.Execute(w, t)
		} else {
			plain.Execute(w, t)
		}
	}
}

var lemma = template.Must(template.New("lemma").Parse(`<span class="lemma">{{ . }}</span>`))
var plain = template.Must(template.New("plain").Parse(`{{ . }}`))

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	cors(w)
	w.Write([]byte("ok"))
}

func cors(w http.ResponseWriter) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
}
