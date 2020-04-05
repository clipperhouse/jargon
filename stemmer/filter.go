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

type filter struct {
	stem func(string, bool) string
}

// English is a Snowball stemmer for English, implemented as a jargon.Filter
var English = &filter{
	stem: english.Stem,
}

// French is a Snowball stemmer for French, implemented as a jargon.Filter
var French = &filter{
	stem: french.Stem,
}

// Norwegian is a Snowball stemmer for Norwegian, implemented as a jargon.Filter
var Norwegian = &filter{
	stem: norwegian.Stem,
}

// Russian is a Snowball stemmer for Russian, implemented as a jargon.Filter
var Russian = &filter{
	stem: russian.Stem,
}

// Spanish is a Snowball stemmer for Spanish, implemented as a jargon.Filter
var Spanish = &filter{
	stem: spanish.Stem,
}

// Swedish is a Snowball stemmer for Swedish, implemented as a jargon.Filter
var Swedish = &filter{
	stem: swedish.Stem,
}

func (f *filter) Filter(incoming *jargon.TokenStream) *jargon.TokenStream {
	t := &tokens{
		filter:   f,
		incoming: incoming,
	}
	return jargon.NewTokenStream(t.next)
}

type tokens struct {
	filter   *filter
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

		stemmed := t.filter.stem(token.String(), true)

		if stemmed == token.String() {
			// Had no effect, send back the original
			return token, nil
		}

		return jargon.NewToken(stemmed, true), nil
	}
}
