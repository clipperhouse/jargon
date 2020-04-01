package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/ascii"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/stackoverflow"
	"github.com/clipperhouse/jargon/stemmer"
)

var filein = flag.String("file", "", "input file path (if none, stdin is used as input)")
var fileout = flag.String("out", "", "output file path (if none, stdout is used as input)")
var html = flag.Bool("html", false, "parse input as html (keep tags whole")
var lang = flag.String("lang", "english", "language of input, relevant when used with -stem. options:\n"+strings.Join(langs, ", "))

// These values aren't actually used, see filters loop below
var xascii = flag.Bool("ascii", false, "a filter to replace diacritics with ascii equivalents, e.g. café → cafe")
var xcontractions = flag.Bool("contractions", false, "a filter to expand contractions, e.g. Would've → Would have")
var xstack = flag.Bool("stack", false, "a filter to recognize tech terms as Stack Overflow tags, e.g. Ruby on Rails → ruby-on-rails")
var xstem = flag.Bool("stem", false, "a filter to stem words using snowball stemmer, e.g. management|manager → manag")

var filterMap = map[string]jargon.Filter{
	"-ascii":        ascii.Fold,
	"-contractions": contractions.Expander,
	"-stack":        stackoverflow.Tags,
	"-stem":         stemmer.English,
}

var langs = []string{"english", "french", "norwegian", "russian", "spanish", "swedish"}
var stemmerMap = map[string]jargon.Filter{
	"english":   stemmer.English,
	"french":    stemmer.French,
	"norwegian": stemmer.Norwegian,
	"russian":   stemmer.Russian,
	"spanish":   stemmer.Spanish,
	"swedish":   stemmer.Swedish,
}

func main() {
	flag.Parse()

	var config = &config{
		HTML: *html,
	}

	//
	// Input
	//

	if *filein != "" {
		// Try to open it
		file, err := os.Open(*filein)
		check(err)
		defer file.Close()

		config.Filein = file
	}

	in, err := os.Stdin.Stat()
	check(err)
	pipedin := (in.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435/70613

	// If no input, display usage
	input := pipedin || config.Filein != nil
	if !input {
		os.Stderr.WriteString(flag.CommandLine.Name() + " takes text from std input and processes it with one or more filters\n\n")
		os.Stderr.WriteString("Flags:\n")
		flag.PrintDefaults()
		return
	}

	// Choose one input *or* the other
	if pipedin && config.Filein != nil {
		err := fmt.Errorf("choose *either* an input -file argument *or* piped input")
		check(err)
	}

	//
	// Filters
	//

	// Loop through filters; order matters, so can't use flag package
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		filter, found := filterMap[arg]
		if found {
			if filter == stemmer.English {
				// Look for a language specification
				if *lang != "" {
					stem, found := stemmerMap[*lang]
					if found {
						filter = stem
					} else {
						var langs []string
						for k := range stemmerMap {
							langs = append(langs, k)
						}
						err := fmt.Errorf("lang %q is not known by %s; options are %s; leave it unspecified to default to english", *lang, flag.CommandLine.Name(), strings.Join(langs, ", "))
						check(err)
					}
				}
			}

			config.Filters = append(config.Filters, filter)
		}
	}

	//
	// Output
	//

	if *fileout != "" {
		// Interpret last arg as output file path
		file, err := os.Create(*fileout)
		check(err)
		defer file.Close()

		config.Fileout = file
	}
	pipedout := config.Fileout == nil

	//
	// Reader
	//

	if pipedin || pipedout {
		// We're limited by the OS pipe buffer, typically 64K with back pressure
		// Using anything larger doesn't buy us anything
		size := 64 * 1024
		if pipedin {
			config.Reader = bufio.NewReaderSize(os.Stdin, size)
		} else {
			config.Reader = bufio.NewReaderSize(config.Filein, size)
		}
	}

	if config.Filein != nil {
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

	if config.Writer == nil {
		// Match the input buffer size; mismatch doesn't buy us anything
		size := config.Reader.Size()
		if pipedout {
			config.Writer = bufio.NewWriterSize(os.Stdout, size)
		} else {
			config.Writer = bufio.NewWriterSize(config.Fileout, size)
		}
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
	if err := config.Writer.Flush(); err != nil {
		check(err)
	}
}

type config struct {
	Filein  *os.File
	Fileout *os.File
	Filters []jargon.Filter
	HTML    bool
	Reader  *bufio.Reader
	Writer  *bufio.Writer
}

func check(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}
