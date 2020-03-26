package jargon_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/stackoverflow"
)

func ExampleTokenize() {
	// Tokenize takes an io.Reader
	text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
	r := strings.NewReader(text)

	tokens := jargon.Tokenize(r)

	// Tokenize returns a Tokens iterator. Iterate by calling Next() until nil, which
	// indicates that the iterator is exhausted.
	for {
		token, err := tokens.Next()
		if err != nil {
			// Because the source is I/O, errors are possible
			log.Fatal(err)
		}
		if token == nil {
			break
		}

		// Do stuff with token
	}

	// Tokens is lazily evaluated; it does the tokenization work as you call Next.
	// This is done to ensure predictble memory usage and performance. It is
	// 'forward-only', which means that once you consume a token, you can't go back.

	// Usually, Tokenize serves as input to Lemmatize
}

func ExampleTokens_Filter() {
	// Lemmatize take tokens and attempts to find their canonical version

	// Lemmatize takes a Tokens iterator, and one or more token filters
	text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
	r := strings.NewReader(text)

	tokens := jargon.Tokenize(r)
	filtered := tokens.Filter(stackoverflow.Tags)

	// Lemmatize returns a Tokens iterator. Iterate by calling Next() until nil, which
	// indicates that the iterator is exhausted.
	for {
		token, err := filtered.Next()
		if err != nil {
			// Because the source is I/O, errors are possible
			log.Fatal(err)
		}
		if token == nil {
			break
		}

		// Do stuff with token
		if token.IsLemma() {
			fmt.Printf("found lemma: %s", token)
		}
	}
}
