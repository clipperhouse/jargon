package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/ascii"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/stackoverflow"
	"github.com/clipperhouse/jargon/stemmer"
)

func main() {
	var config = &config{}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-html" {
			config.HTML = true
			continue
		}

		filter, found := filterMap[arg]
		if found {
			if arg == "-stem" {
				// Lookahead for language
				if len(args) > i+1 {
					stemlang := args[i+1]
					stemmer, ok := stemmerMap[stemlang]
					if ok {
						filter = stemmer
						i++
					}
				}
				// Otherwise defaults to English
			}

			config.Filters = append(config.Filters, filter)
		}

		switch {
		case i == 0 && !found:
			// Interpret first arg as input file path
			file, err := os.Open(arg)
			if err != nil {
				unknown := fmt.Errorf("%s is not a valid flag or file", arg)
				os.Stderr.WriteString(unknown.Error() + "\n")
			}
			check(err)
			defer file.Close()

			config.Filein = file
		case i == len(args)-1 && !found:
			// Interpret last arg as output file path
			file, err := os.Create(arg)
			if err != nil {
				unknown := fmt.Errorf("%s is not a valid flag or file", arg)
				os.Stderr.WriteString(unknown.Error() + "\n")
			}
			check(err)
			defer file.Close()

			config.Fileout = file
		case !found:
			err := fmt.Errorf("unknown flag %q", arg)
			check(err)
		}
	}

	// Piped?
	fi, err := os.Stdin.Stat()
	check(err)
	piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435/70613

	input := piped || config.Filein != nil
	if !input {
		if len(config.Filters) > 0 {
			os.Stderr.WriteString("indicate a file path as the first argument, or pipe in text\n")
		}
		os.Stderr.WriteString("Usage:\n")
		return
	}

	if piped && config.Filein != nil {
		err := fmt.Errorf("choose *either* an input file argument or piped input")
		check(err)
	}

	if piped {
		// We're limited by pipe buffer size, typically 64K with back pressure
		size := 64 * 1024
		if piped {
			config.Reader = bufio.NewReaderSize(os.Stdin, size)
		} else {
			config.Reader = bufio.NewReaderSize(config.Filein, size)
		}
	}

	if config.Reader == nil {
		fi, err := config.Filein.Stat()
		check(err)

		size := fi.Size()
		switch {
		case size <= 4*1024:
			// Minimum of 4K
			config.Reader = bufio.NewReaderSize(config.Filein, 4*1024)
		case size <= 1024*1024:
			// Aim for a right-sized buffer (single read, perhaps) up to 1MB
			config.Reader = bufio.NewReaderSize(config.Filein, int(size))
		default:
			// Otherwise, use 1MB buffer size, better perf over default, but not huge
			config.Reader = bufio.NewReaderSize(config.Filein, 1024*1024)
		}
	}

	if config.Fileout != nil {
		config.Writer = config.Fileout
	} else {
		config.Writer = os.Stdout
	}

	var tokens *jargon.Tokens
	if config.HTML {
		tokens = jargon.TokenizeHTML(config.Reader)
	} else {
		tokens = jargon.Tokenize(config.Reader)
	}
	for _, f := range config.Filters {
		tokens = tokens.Filter(f)
	}

	if _, err := tokens.WriteTo(config.Writer); err != nil {
		check(err)
	}
}

type config struct {
	Filein  *os.File
	Fileout *os.File
	Filters []jargon.Filter
	HTML    bool
	Reader  *bufio.Reader
	Writer  io.Writer
}

func check(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
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
