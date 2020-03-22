package jargon

import (
	"log"
	"strings"
)

func ExampleTokens() {
	// Tokens is an iterator resulting from a call to Tokenize or Filter

	text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
	r := strings.NewReader(text)
	tokens := Tokenize(r)

	// Iterate by calling Next() until nil, which indicates that the iterator is exhausted.
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

	// As an iterator, Tokens is 'forward-only', which means that
	// once you consume a token, you can't go back.

	// See also the convenience methods String, ToSlice, WriteTo
}
