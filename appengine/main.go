package main

import (
	"html/template"
	"net/http"
	"strings"

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

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 1 {
		http.NotFound(w, r)
		return
	}

	switch route := parts[1]; {
	case route == "text" || route == "jargon":
		textHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func textHandler(w http.ResponseWriter, r *http.Request) {
	text := r.PostFormValue("text")

	if len(strings.TrimSpace(text)) > 0 {
		tokens := jargon.TechProse.Tokenize(text)

		lemmatized := jargon.StackExchange.LemmatizeTokens(tokens)

		for t := range lemmatized {
			if t.IsLemma() {
				lemma.Execute(w, t)
			} else {
				w.Write([]byte(t.String()))
			}
		}
	}
}

var lemma = template.Must(template.New("span").Parse(`<span class="lemma">{{ .String }}</span>`))

func isAjax(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
