// Package stemmer offers the Snowball stemmer in several languages.
package stemmer

import (
	"github.com/clipperhouse/jargon"
	"github.com/kljensen/snowball/english"
	"github.com/kljensen/snowball/french"
	"github.com/kljensen/snowball/norwegian"
	"github.com/kljensen/snowball/russian"
	"github.com/kljensen/snowball/spanish"
	"github.com/kljensen/snowball/swedish"
)

type stemmer struct {
	stem func(string, bool) string
}

// English is a Snowball stemmer for English, implemented as a jargon.Filter
var English = newStemmer(english.Stem).Stem

// French is a Snowball stemmer for French, implemented as a jargon.Filter
var French = newStemmer(french.Stem).Stem

// Norwegian is a Snowball stemmer for Norwegian, implemented as a jargon.Filter
var Norwegian = newStemmer(norwegian.Stem).Stem

// Russian is a Snowball stemmer for Russian, implemented as a jargon.Filter
var Russian = newStemmer(russian.Stem).Stem

// Spanish is a Snowball stemmer for Spanish, implemented as a jargon.Filter
var Spanish = newStemmer(spanish.Stem).Stem

// Swedish is a Snowball stemmer for Swedish, implemented as a jargon.Filter
var Swedish = newStemmer(swedish.Stem).Stem

// newStemmer creates a new stemmer
func newStemmer(stem func(string, bool) string) *stemmer {
	return &stemmer{
		stem: stem,
	}
}

func (st *stemmer) Stem(incoming *jargon.TokenStream) *jargon.TokenStream {
	t := &tokens{
		stemmer:  st,
		incoming: incoming,
	}
	return jargon.NewTokenStream(t.next)
}

type tokens struct {
	stemmer  *stemmer
	incoming *jargon.TokenStream
}

func (t *tokens) next() (*jargon.Token, error) {
	for {
		token, err := t.incoming.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			return nil, nil
		}

		// Only interested in stemming words
		if token.IsPunct() || token.IsSpace() {
			return token, nil
		}

		stemmed := t.stemmer.stem(token.String(), true)

		if stemmed == token.String() {
			// Had no effect, send back the original
			return token, nil
		}

		return jargon.NewToken(stemmed, true), nil
	}
}
