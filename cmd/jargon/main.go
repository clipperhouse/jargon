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

func check(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

var f, s, u, o string

func init() {
	flag.StringVar(&f, "f", "", "Input file path")
	flag.StringVar(&s, "s", "", "A (quoted) string to lemmatize")
	flag.StringVar(&u, "u", "", "A URL to fetch and lemmatize")
	flag.StringVar(&o, "o", "", "Output file path")
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

	return lem(file)
}

func lemString(s string) error {
	r := strings.NewReader(s)
	return lem(r)
}

func lemURL(u string) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return lem(resp.Body)
}

func lemStdin() error {
	return lem(os.Stdin)
}

var lemmatizer = jargon.NewLemmatizer(stackexchange.Dictionary, 3)

func lem(r io.Reader) error {
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

	br := bufio.NewReader(r)

	tokens := jargon.Tokenize(br)
	lemmas := lemmatizer.Lemmatize(tokens)

	for {
		t := lemmas.Next()
		if t == nil {
			break
		}
		_, err := w.WriteString(t.String())
		check(err)
	}

	// Flush the buffer as a last step; return error if any
	return w.Flush()
}
