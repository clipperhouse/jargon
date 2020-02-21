// Package stemmer offers the Snowball stemmer in several languages.
package stemmer

import (
	"github.com/kljensen/snowball/english"
	"github.com/kljensen/snowball/french"
	"github.com/kljensen/snowball/norwegian"
	"github.com/kljensen/snowball/russian"
	"github.com/kljensen/snowball/spanish"
	"github.com/kljensen/snowball/swedish"
)

type filter struct {
	stemmer func(string, bool) string
}

// English is a Snowball stemmer for English, implemented as a jargon.TokenFilter
var English = &filter{
	stemmer: english.Stem,
}

// French is a Snowball stemmer for French, implemented as a jargon.TokenFilter
var French = &filter{
	stemmer: french.Stem,
}

// Norwegian is a Snowball stemmer for Norwegian, implemented as a jargon.TokenFilter
var Norwegian = &filter{
	stemmer: norwegian.Stem,
}

// Russian is a Snowball stemmer for Russian, implemented as a jargon.TokenFilter
var Russian = &filter{
	stemmer: russian.Stem,
}

// Spanish is a Snowball stemmer for Spanish, implemented as a jargon.TokenFilter
var Spanish = &filter{
	stemmer: spanish.Stem,
}

// Swedish is a Snowball stemmer for Swedish, implemented as a jargon.TokenFilter
var Swedish = &filter{
	stemmer: swedish.Stem,
}

func (f *filter) Lookup(s []string) (string, bool) {
	word := s[0]
	stem := f.stemmer(word, true)

	if stem != word {
		return stem, true
	}

	return word, false
}

func (f *filter) MaxGramLength() int {
	return 1
}
