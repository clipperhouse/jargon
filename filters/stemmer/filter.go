// Package stemmer offers the Snowball stemmer in several languages
package stemmer

import (
	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/mapper"
	"github.com/kljensen/snowball/english"
	"github.com/kljensen/snowball/french"
	"github.com/kljensen/snowball/norwegian"
	"github.com/kljensen/snowball/russian"
	"github.com/kljensen/snowball/spanish"
	"github.com/kljensen/snowball/swedish"
)

// English is a Snowball stemmer for English, implemented as a jargon.Filter
var English jargon.Filter = newStemmer(english.Stem)

// French is a Snowball stemmer for French, implemented as a jargon.Filter
var French = newStemmer(french.Stem)

// Norwegian is a Snowball stemmer for Norwegian, implemented as a jargon.Filter
var Norwegian = newStemmer(norwegian.Stem)

// Russian is a Snowball stemmer for Russian, implemented as a jargon.Filter
var Russian = newStemmer(russian.Stem)

// Spanish is a Snowball stemmer for Spanish, implemented as a jargon.Filter
var Spanish = newStemmer(spanish.Stem)

// Swedish is a Snowball stemmer for Swedish, implemented as a jargon.Filter
var Swedish = newStemmer(swedish.Stem)

// newStemmer creates a new stemmer
func newStemmer(stem func(string, bool) string) jargon.Filter {
	f := func(token *jargon.Token) *jargon.Token {
		// Only interested in stemming words
		if token.IsPunct() || token.IsSpace() {
			return token
		}

		stemmed := stem(token.String(), true)

		if stemmed == token.String() {
			// Had no effect, send back the original
			return token
		}

		return jargon.NewToken(stemmed, true)
	}

	return mapper.NewFilter(f)
}
