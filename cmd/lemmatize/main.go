package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/stackexchange"
)

var lemmatizer = jargon.NewLemmatizer(stackexchange.Dictionary, 3)

func main() {
	flag.Parse()

	if len(filePath) > 0 {
		err := lemFile(filePath)
		if err != nil {
			panic(err)
		}
	}
}

func lemFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	tokens := jargon.Tokenize(file)
	lemmas := lemmatizer.Lemmatize(tokens)

	for l := range lemmas {
		fmt.Print(l.String())
	}

	return nil
}

var filePath string

func init() {
	flag.StringVar(&filePath, "f", "", "")
}
