package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/stackexchange"
)

func main() {
	flag.Parse()

	var err error

	switch {
	case len(f) > 0:
		err = lemFile(f)
	case len(s) > 0:
		lemString(s)
		fmt.Print("\n")
	case len(u) > 0:
		err = lemURL(u)
	default:
		flag.PrintDefaults()
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

var f, s, u string

func init() {
	flag.StringVar(&f, "f", "", "A file path to lemmatize")
	flag.StringVar(&s, "s", "", "A (quoted) string to lemmatize")
	flag.StringVar(&u, "u", "", "A URL to fetch and lemmatize")
}

// turns out that buffering on the way out performs ~40% better, at least on my machine
var w = bufio.NewWriter(os.Stdout)

func lemFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	lem(file)
	return nil
}

func lemString(s string) {
	r := strings.NewReader(s)
	lem(r)
}

func lemURL(u string) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	lem(resp.Body)
	return nil
}

var lemmatizer = jargon.NewLemmatizer(stackexchange.Dictionary, 3)

func lem(r io.Reader) {
	br := bufio.NewReader(r)
	b, _ := br.Peek(512) // ignore the error here, it usually means we can't get 512 bytes, but returns what is gotten anyway
	c := http.DetectContentType(b)

	var tokens <-chan jargon.Token

	if strings.HasPrefix(c, "text/html") {
		tokens = jargon.TokenizeHTML(br)
	} else {
		tokens = jargon.Tokenize(br)
	}

	lemmas := lemmatizer.Lemmatize(tokens)

	for l := range lemmas {
		w.WriteString(l.String())
	}
}
