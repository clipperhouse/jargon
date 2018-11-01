// A demo of jargon for use on Google App Engine
package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/stackexchange"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/_ah/health", healthCheckHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
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

	var b bytes.Buffer

	for {
		t := lemmatized.Next()
		if t == nil {
			break
		}

		// we buffer (instead of writing directly to Response) because the Body
		// will be closed if we read and write concurrently:
		// https://github.com/golang/go/issues/15527
		if t.IsLemma() {
			lemma.Execute(&b, t)
		} else {
			plain.Execute(&b, t)
		}
	}

	_, err := b.WriteTo(w)
	if err != nil {
		log.Print(err)
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
