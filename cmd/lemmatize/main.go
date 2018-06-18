package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/stackexchange"
)

func main() {
	flag.Parse()

	switch {
	case len(f) > 0:
		err := lemFile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case len(s) > 0:
		lemString(s)
		fmt.Print("\n")
	default:
		flag.PrintDefaults()
	}
}

var f, s string

func init() {
	flag.StringVar(&f, "f", "", "A file path to lemmatize")
	flag.StringVar(&s, "s", "", "A (quoted) string to lemmatize")
}

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

var lemmatizer = jargon.NewLemmatizer(stackexchange.Dictionary, 3)

func lem(r io.Reader) {
	tokens := jargon.Tokenize(r)
	lemmas := lemmatizer.Lemmatize(tokens)

	for l := range lemmas {
		fmt.Print(l.String())
	}
}
