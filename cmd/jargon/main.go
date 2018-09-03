package main

import (
	"bufio"
	"flag"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/stackexchange"
)

func main() {
	flag.Parse()

	switch {
	case len(f) > 0:
		check(lemFile(f))
	case len(s) > 0:
		check(lemString(s))
		w.WriteByte('\n')
	case len(u) > 0:
		check(lemURL(u))
	default:
		flag.PrintDefaults()
	}
}

func check(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
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

	return lem(file, w)
}

func lemString(s string) error {
	r := strings.NewReader(s)
	return lem(r, w)
}

func lemURL(u string) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return lem(resp.Body, w)
}

var lemmatizer = jargon.NewLemmatizer(stackexchange.Dictionary, 3)

func lem(r io.Reader, w *bufio.Writer) error {
	br := bufio.NewReader(r)

	tokens := jargon.Tokenize(br)
	lemmas := lemmatizer.Lemmatize(tokens)

	for {
		t := lemmas.Next()
		if t == nil {
			break
		}
		w.WriteString(t.String())
	}

	// Flush the buffer as a last step; return error if any
	return w.Flush()
}
