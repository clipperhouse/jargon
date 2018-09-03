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
		w.WriteByte('\n')
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

var f, s, u string

func init() {
	flag.StringVar(&f, "f", "", "A file path to lemmatize")
	flag.StringVar(&s, "s", "", "A (quoted) string to lemmatize")
	flag.StringVar(&u, "u", "", "A URL to fetch and lemmatize")
	flag.Usage = func() {
		cmd := os.Args[0]
		out := flag.CommandLine.Output()

		usage := `
Usage: %[1]s accepts piped text from tools such as cat, curl or echo, via Stdin
		
  Example: echo "I luv Rails" | %[1]s

Alternatively, use %[1]s 'standalone' by passing flags for text sources:

`
		fmt.Fprintf(out, usage, cmd)
		flag.PrintDefaults()
		fmt.Fprintf(out, "\n  Example: jargon -f /path/to/file.txt\n\n")
		fmt.Fprintf(out, "Results are piped to Stdout (regardless of input)\n\n")
	}
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

func lemStdin() error {
	return lem(os.Stdin, w)
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
