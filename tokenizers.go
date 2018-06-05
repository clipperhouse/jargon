package jargon

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Tokenize returns a channel of Tokens from a io.Reader, intended to be ranged over thus:
//	tokens := Tokenize(reader)
//	for t := range tokens {
// 		// do stuff
//	}
//
// The tokenizer is targeted to English text that contains tech terms, so things like C++ and .Net are handled as single units.
// It respects Unicode definitions of 'punctuation' and 'symbol', with a few hard-coded exceptions. Symbols are treated as word characters (like alphanumerics), allowing things like email addresses, hashtags and @-handles to be understood as a single token.
// Dots are allowed to lead words, and to appear mid-word, allowing things like .Net and Node.js. Trailing dots are considered end-of-sentence.
// [@, #, -, *, %, /, \] are hard-coded as symbols, even though Unicode specifies them as punctuation. See http://www.unicode.org/faq/punctuation_symbols.html
// All other punctuation terminates words, as does white space.
//
// It returns all tokens (including white space), so text can be reconstructed with fidelity. Ignoring (say) whitespace is a decision for the caller.
func Tokenize(r io.Reader) chan Token {
	b := newReader(r)
	return b.tokens
}

// TokenizeHTML tokenizes HTML. Text nodes are tokenized using jargon.Tokenize; everything else (tags, comments) are left verbatim.
// It returns a channel of Tokens, intended to be ranged over thus:
//	tokens := TokenizeHTML(string)
//	for t := range tokens {
// 		// do stuff
//	}
//
// It returns all tokens (including white space), so text can be reconstructed with fidelity. Ignoring (say) whitespace is a decision for the caller.
func TokenizeHTML(r io.Reader) chan Token {
	result := make(chan Token, 0)
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
				r := strings.NewReader(tok.Data)
				words := Tokenize(r)
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
