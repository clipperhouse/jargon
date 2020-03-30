package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/ascii"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/stackoverflow"
	"github.com/clipperhouse/jargon/stemmer"
)

func main() {
	tokenize := jargon.Tokenize
	var filters []jargon.Filter

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := strings.ToLower(args[i])

		t, ok := tokenizeMap[arg]
		if ok {
			tokenize = t
			continue
		}

		f, ok := filterMap[arg]

		if !ok {
			err := fmt.Errorf("unknown flag %q", arg)
			check(err)
		}

		if arg == "-stem" {
			// Lookahead for language
			if len(args) > i+1 {
				stemlang := args[i+1]
				stemmer, ok := stemmerMap[stemlang]
				if ok {
					f = stemmer
					i++
				}
			}
			// Otherwise defaults to English
		}

		filters = append(filters, f)
	}

	fi, err := os.Stdin.Stat()
	check(err)

	piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435/70613

	if piped {
		tokens := tokenize(os.Stdin)
		for _, f := range filters {
			tokens = tokens.Filter(f)
		}

		_, err := tokens.WriteTo(os.Stdout)
		check(err)
	}
}

func check(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}

var tokenizeMap = map[string](func(io.Reader) *jargon.Tokens){
	"-html": jargon.TokenizeHTML,
}

var filterMap = map[string]jargon.Filter{
	"-ascii":        ascii.Fold,
	"-contractions": contractions.Expander,
	"-stack":        stackoverflow.Tags,
	"-stem":         stemmer.English,
}

var stemmerMap = map[string]jargon.Filter{
	"english":   stemmer.English,
	"french":    stemmer.French,
	"norwegian": stemmer.Norwegian,
	"russian":   stemmer.Russian,
	"spanish":   stemmer.Spanish,
	"swedish":   stemmer.Swedish,
}
