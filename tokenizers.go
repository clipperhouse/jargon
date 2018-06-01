package jargon

import (
	"strings"

	"golang.org/x/net/html"
)

type techProse struct{}

// TechProse tokenizer attempts to tokenize English text that contains tech terms.
// It respects Unicode definitions of 'punctuation' and 'symbol', with a few hard-coded exceptions (below).
// Symbols are treated as word characters (like alphanumerics), allowing things like email addresses, hashtags and @-handles to be understood as a single token.
// Dots are allowed to lead words, and to appear mid-word, allowing things like .Net and Node.js. Trailing dots are considered end-of-sentence.
// [@, #, -, *, %, /, \] are hard-coded as symbols, even though Unicode specifies them as punctuation. See http://www.unicode.org/faq/punctuation_symbols.html
// All other punctuation terminates words, as does white space.
// Like the other tokenizers in this package, it returns all tokens (including white space), so text can be reconstructed with fidelity. If callers don't want white space, they'll need to filter.
// TODO: URLs
var TechProse = &techProse{}

// Tokenize returns a channel of Tokens, intended to be ranged over thus:
//	tokens := TechProse.Tokenize(string)
//	for t := range tokens {
// 		// do stuff
//	}
func (t *techProse) Tokenize(text string) chan Token {
	s := strings.NewReader(text)
	b := newReader(s)
	return b.tokens
}

type techHTML struct{}

// TechHTML is a tokenizer for HTML text. Text nodes are tokenized using TechProse; tags and comments left verbatim.
var TechHTML = &techHTML{}

// Tokenize returns a channel of Tokens, intended to be ranged over thus:
//	tokens := TechHTML.Tokenize(string)
//	for t := range tokens {
// 		// do stuff
//	}
func (t *techHTML) Tokenize(text string) chan Token {
	result := make(chan Token, 0)
	r := strings.NewReader(text)
	z := html.NewTokenizer(r)

	go func() {
		for {
			tt := z.Next()

			if tt == html.ErrorToken {
				// Presumably eof
				close(result)
				break
			}

			switch tok := z.Token(); {
			case tok.Type == html.TextToken:
				words := TechProse.Tokenize(tok.Data)
				for w := range words {
					result <- w
				}
			default:
				// Everything else is punct for our purposes
				new := Token{
					value: tok.String(),
					punct: true,
					space: false,
				}
				result <- new
			}
		}
	}()

	return result
}
