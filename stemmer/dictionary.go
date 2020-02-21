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

type dictionary struct {
	stemmer func(string, bool) string
}

// English is a Snowball stemmer for English, implemented as a jargon.Dictionary
var English = &dictionary{
	stemmer: english.Stem,
}

// French is a Snowball stemmer for French, implemented as a jargon.Dictionary
var French = &dictionary{
	stemmer: french.Stem,
}

// Norwegian is a Snowball stemmer for Norwegian, implemented as a jargon.Dictionary
var Norwegian = &dictionary{
	stemmer: norwegian.Stem,
}

// Russian is a Snowball stemmer for Russian, implemented as a jargon.Dictionary
var Russian = &dictionary{
	stemmer: russian.Stem,
}

// Spanish is a Snowball stemmer for Spanish, implemented as a jargon.Dictionary
var Spanish = &dictionary{
	stemmer: spanish.Stem,
}

// Swedish is a Snowball stemmer for Swedish, implemented as a jargon.Dictionary
var Swedish = &dictionary{
	stemmer: swedish.Stem,
}

func (d *dictionary) Lookup(s []string) (string, bool) {
	word := s[0]
	stem := d.stemmer(word, true)

	if stem != word {
		return stem, true
	}

	return word, false
}

func (d *dictionary) MaxGramLength() int {
	return 1
}
