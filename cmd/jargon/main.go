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
	"github.com/clipperhouse/jargon/stackexchange"
)

func main() {
	flag.Parse()

	dictionaries := determineDictionaries(tech, num, cont)

	switch {
	case len(f) > 0:
		check(lemFile(f, dictionaries))
	case len(s) > 0:
		check(lemString(s, dictionaries))
	case len(u) > 0:
		check(lemURL(u, dictionaries))
	default:
		// No flags? Check to see if piped, otherwise print help.
		fi, err := os.Stdin.Stat()
		check(err)

		piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435/70613
		if piped {
			check(lemStdin(dictionaries))
		} else {
			flag.Usage()
		}
	}
}

func determineDictionaries(tech, num, cont bool) []jargon.Dictionary {
	// splitting this out into a func to allow testing

	var result []jargon.Dictionary

	none := !tech && !num && !cont

	if tech || none {
		result = append(result, stackexchange.Dictionary)
	}

	if num {
		result = append(result, numbers.Dictionary)
	}

	if cont {
		result = append(result, contractions.Dictionary)
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
	flag.BoolVar(&tech, "tech", false, "Lemmatize technology terms using the StackExchange dictionary")
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
flags to choose other dictionaries.

  -tech
    	Lemmatize technology terms to Stack Overflow-style tags
    	(e.g. "Ruby on Rails" → "ruby-on-rails"). If no dictionary is
    	specified, this is the default.
  -cont
    	Expand contractions (e.g. "didn't" → "did not")
  -num
    	Lemmatize number phrases (e.g. "three hundred" → "300")

`
		fmt.Fprintf(out, usage, cmd)
	}
}

func lemFile(filePath string, dictionaries []jargon.Dictionary) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var tokens jargon.Tokens

	ext := filepath.Ext(filePath)
	switch ext {
	case ".html", ".htm":
		tokens = jargon.TokenizeHTML(file)
	default:
		tokens = jargon.Tokenize(file)
	}

	return lem(tokens, dictionaries)
}

func lemString(s string, dictionaries []jargon.Dictionary) error {
	r := strings.NewReader(s)
	tokens := jargon.Tokenize(r)
	return lem(tokens, dictionaries)
}

func lemURL(u string, dictionaries []jargon.Dictionary) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokens jargon.Tokens

	ct := resp.Header.Get("Content-Type")
	html := strings.HasPrefix(ct, "text/html")
	if html {
		tokens = jargon.TokenizeHTML(resp.Body)
	} else {
		tokens = jargon.Tokenize(resp.Body)
	}

	return lem(tokens, dictionaries)
}

func lemStdin(dictionaries []jargon.Dictionary) error {
	tokens := jargon.Tokenize(os.Stdin)
	return lem(tokens, dictionaries)
}

func lem(tokens jargon.Tokens, dictionaries []jargon.Dictionary) error {

	// switch tokens.(type) {
	// case *jargon.HTMLTokens:
	// 	fmt.Println("tokenized html")
	// case *jargon.TextTokens:
	// 	fmt.Println("tokenized plain text")
	// default:
	// 	panic("unknown text type")
	// }

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

	tokens = lemAll(tokens, dictionaries)

	_, err := tokens.WriteTo(w)
	check(err)

	// Flush the buffer as a last step; return error if any
	return w.Flush()
}

func lemAll(tokens jargon.Tokens, dictionaries []jargon.Dictionary) jargon.Tokens {
	for _, dictionary := range dictionaries {
		tokens = jargon.Lemmatize(tokens, dictionary)
	}
	return tokens
}
