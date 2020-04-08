package jargon_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/stackoverflow"
)

func ExampleTokenStream_Scan() {
	// TokensStream is an iterator resulting from a call to Tokenize or Filter

	text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
	stream := jargon.TokenizeString(text).Filter(stackoverflow.Tags)

	// Loop while Scan() returns true. Scan() will return false on error or end of tokens.
	for stream.Scan() {
		token := stream.Token()
		// Do stuff with token
		fmt.Print(token)
	}

	if err := stream.Err(); err != nil {
		// Because the source is I/O, errors are possible
		log.Fatal(err)
	}

	// As an iterator, TokenStream is 'forward-only', which means that
	// once you consume a token, you can't go back.

	// See also the convenience methods String, ToSlice, WriteTo
}

func ExampleTokenStream_Next() {
	// TokensStream is an iterator resulting from a call to Tokenize or Filter

	text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
	r := strings.NewReader(text)
	tokens := jargon.Tokenize(r)

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

	// As an iterator, TokenStream is 'forward-only', which means that
	// once you consume a token, you can't go back.

	// See also the convenience methods String, ToSlice, WriteTo
}
