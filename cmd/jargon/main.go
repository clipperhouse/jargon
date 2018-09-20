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
	"github.com/clipperhouse/jargon/numbers"
	"github.com/clipperhouse/jargon/stackexchange"
)

func main() {
	flag.Parse()

	lemmatizers = determineLemmatizers(tech, num)

	switch {
	case len(f) > 0:
		check(lemFile(f))
	case len(s) > 0:
		check(lemString(s))
	case len(u) > 0:
		check(lemURL(u))
	default:
		// No flags? Check to see if piped, otherwise print help.
		fi, err := os.Stdin.Stat()
		check(err)

		piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435/70613
		if piped {
			check(lemStdin())
		} else {
			flag.Usage()
		}
	}
}

var lemmatizers []*jargon.Lemmatizer

func determineLemmatizers(tech, num bool) []*jargon.Lemmatizer {
	// splitting this out into a func to allow testing

	var result []*jargon.Lemmatizer

	none := !tech && !num

	if tech || none {
		lem := jargon.NewLemmatizer(stackexchange.Dictionary, 3)
		result = append(result, lem)
	}

	if num {
		lem := jargon.NewLemmatizer(numbers.Dictionary, 3)
		result = append(result, lem)
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
var tech, num bool

func init() {
	flag.StringVar(&f, "f", "", "Input file path")
	flag.StringVar(&s, "s", "", "A (quoted) string to lemmatize")
	flag.StringVar(&u, "u", "", "A URL to fetch and lemmatize")
	flag.StringVar(&o, "o", "", "Output file path. If omitted, output goes to Stdout.")
	// flag.BoolVar(&tech, "tech", false, "Lemmatize technology terms using the StackExchange dictionary")
	// flag.BoolVar(&num, "num", false, `Lemmatize number phrases (e.g. "three hundred → "300")`)
	flag.Usage = func() {
		cmd := os.Args[0]
		out := flag.CommandLine.Output()

		usage := `
Usage:

%[1]s accepts piped UTF8 text from Stdin and pipes lemmatized text to Stdout
		
  Example: echo "I luv Rails" | %[1]s

Alternatively, use %[1]s 'standalone' by passing flags for inputs and outputs:

`
		fmt.Fprintf(out, usage, cmd)
		flag.PrintDefaults()
		fmt.Fprintf(out, "\n  Example: %s -f /path/to/original.txt -o /path/to/lemmatized.txt\n\n", cmd)
	}
}

func lemFile(filePath string) error {
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

	return lem(tokens)
}

func lemString(s string) error {
	r := strings.NewReader(s)
	tokens := jargon.Tokenize(r)
	return lem(tokens)
}

func lemURL(u string) error {
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

	return lem(tokens)
}

func lemStdin() error {
	tokens := jargon.Tokenize(os.Stdin)
	return lem(tokens)
}

func lem(tokens jargon.Tokens) error {

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

	tokens = lemAll(tokens, lemmatizers)

	for {
		t := tokens.Next()
		if t == nil {
			break
		}
		_, err := w.WriteString(t.String())
		check(err)
	}

	// Flush the buffer as a last step; return error if any
	return w.Flush()
}

func lemAll(tokens jargon.Tokens, lems []*jargon.Lemmatizer) jargon.Tokens {
	for _, lem := range lems {
		tokens = lem.Lemmatize(tokens)
	}
	return tokens
}
