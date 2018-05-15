package jargon

import (
	"strings"

	"golang.org/x/net/html"
)

type techHTML struct{}

// TechHTML is a tokenizer for HTML text. Text nodes are tokenized using TechProse; tags and comments left verbatim.
var TechHTML = &techHTML{}

func (t *techHTML) Tokenize(text string) []Token {
	result := make([]Token, 0)
	r := strings.NewReader(text)
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			// Presumably eof
			break
		}

		switch tok := z.Token(); {
		case tok.Type == html.TextToken:
			words := TechProse.Tokenize(tok.Data)
			result = append(result, words...)
		default:
			// Everything else is punct for our purposes
			new := Token{
				value: tok.String(),
				punct: true,
				space: false,
			}
			result = append(result, new)
		}
	}

	return result
}
