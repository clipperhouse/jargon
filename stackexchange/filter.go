package stackexchange

import (
	"fmt"

	"github.com/clipperhouse/jargon/synonyms"
)

//go:generate go run generate/main.go

// Tags is the main exported Tags of Stack Exchange tags and synonyms, from the following Stack Exchange sites: Stack Overflow,
// Server Fault, Game Dev and Data Science. It's indended to identify canonical tags (technologies),
// e.g. Ruby on Rails (3 words) will be replaced with ruby-on-rails (1 word).
var Tags *synonyms.Filter

func init() {
	ignore := []rune{' ', '-', '.', '/'}
	filter, err := synonyms.NewFilter(mappings, synonyms.IgnoreCase, synonyms.Ignore(ignore...))
	fmt.Println(filter.MaxGramLength())
	if err != nil {
		panic(err)
	}
	Tags = filter
}
