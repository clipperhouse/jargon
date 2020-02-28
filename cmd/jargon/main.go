package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/numbers"
	"github.com/clipperhouse/jargon/stackoverflow"
)

func main() {
	flag.Parse()

	filters := determineFilters(tech, num, cont)

	switch {
	case len(f) > 0:
		check(lemFile(f, filters))
	case len(s) > 0:
		check(lemString(s, filters))
	case len(u) > 0:
		check(lemURL(u, filters))
	default:
		// No flags? Check to see if piped, otherwise print help.
		fi, err := os.Stdin.Stat()
		check(err)

		piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435/70613
		if piped {
			check(lemStdin(filters))
		} else {
			flag.Usage()
		}
	}
}

func determineFilters(tech, num, cont bool) []jargon.TokenFilter {
	// splitting this out into a func to allow testing

	var result []jargon.TokenFilter

	none := !tech && !num && !cont

	if tech || none {
		result = append(result, stackoverflow.Tags)
	}

	if num {
		result = append(result, numbers.Filter)
	}

	if cont {
		result = append(result, contractions.Expander)
	}

	return result
}

func check(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

var f, s, u, o string
var tech, num, cont bool

func init() {
	flag.StringVar(&f, "f", "", "Input file path")
	flag.StringVar(&s, "s", "", "A (quoted) string to lemmatize")
	flag.StringVar(&u, "u", "", "A URL to fetch and lemmatize")
	flag.StringVar(&o, "o", "", "Output file path. If omitted, output goes to Stdout.")
	flag.BoolVar(&tech, "tech", false, "Lemmatize technology terms using the Stack Overflow token filter")
	flag.BoolVar(&num, "num", false, `Lemmatize number phrases (e.g. "three hundred → "300")`)
	flag.BoolVar(&cont, "cont", false, `Expand contractions (e.g. "didn't → "did not")`)
	flag.Usage = func() {
		cmd := os.Args[0]
		out := flag.CommandLine.Output()

		usage := `
Usage:

%[1]s accepts piped UTF8 text from Stdin and pipes lemmatized text to Stdout
		
  Example: echo "I luv Rails" | %[1]s

Alternatively, use %[1]s 'standalone' by passing flags for inputs and outputs:

  -f string
    	Input file path
  -o string
    	Output file path. If omitted, output goes to Stdout.
  -s string
    	A (quoted) string to lemmatize
  -u string
    	A URL to fetch and lemmatize

  Example: %[1]s -f /path/to/original.txt -o /path/to/lemmatized.txt

By default, %[1]s uses a dictionary of technology terms. Pass the following 
flags to choose other filters.

  -tech
    	Lemmatize technology terms to Stack Overflow-style tags
    	(e.g. "Ruby on Rails" → "ruby-on-rails"). If no filter is
    	specified, this is the default.
  -cont
    	Expand contractions (e.g. "didn't" → "did not")
  -num
    	Lemmatize number phrases (e.g. "three hundred" → "300")

`
		fmt.Fprintf(out, usage, cmd)
	}
}

func lemFile(filePath string, filters []jargon.TokenFilter) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var tokens *jargon.Tokens

	ext := filepath.Ext(filePath)
	switch ext {
	case ".html", ".htm":
		tokens = jargon.TokenizeHTML(file)
	default:
		tokens = jargon.Tokenize(file)
	}

	return lem(tokens, filters)
}

func lemString(s string, filters []jargon.TokenFilter) error {
	r := strings.NewReader(s)
	tokens := jargon.Tokenize(r)
	return lem(tokens, filters)
}

func lemURL(u string, filters []jargon.TokenFilter) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokens *jargon.Tokens

	ct := resp.Header.Get("Content-Type")
	html := strings.HasPrefix(ct, "text/html")
	if html {
		tokens = jargon.TokenizeHTML(resp.Body)
	} else {
		tokens = jargon.Tokenize(resp.Body)
	}

	return lem(tokens, filters)
}

func lemStdin(filters []jargon.TokenFilter) error {
	tokens := jargon.Tokenize(os.Stdin)
	return lem(tokens, filters)
}

func lem(tokens *jargon.Tokens, filters []jargon.TokenFilter) error {
	var w *bufio.Writer

	if len(o) > 0 { // output file flag
		f, err := os.Create(o)
		check(err)
		defer f.Close()

		w = bufio.NewWriter(f)
	}

	if w == nil {
		w = bufio.NewWriter(os.Stdout)
	}

	tokens = lemAll(tokens, filters)

	_, err := tokens.WriteTo(w)
	check(err)

	// Flush the buffer as a last step; return error if any
	return w.Flush()
}

func lemAll(tokens *jargon.Tokens, filters []jargon.TokenFilter) *jargon.Tokens {
	for _, filter := range filters {
		tokens = tokens.Lemmatize(filter)
	}
	return tokens
}
