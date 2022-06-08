package jargon

import (
	"io"
	"strings"

	"github.com/clipperhouse/uax29/iterators"
	"github.com/clipperhouse/uax29/words"
)

// Tokenize tokenizes a reader into a stream of tokens. Iterate through the stream by calling Scan() or Next().
//
// Its uses several specs from Unicode Text Segmentation https://unicode.org/reports/tr29/. It's not a full implementation, but a decent approximation for many mainstream cases.
//
// Tokenize returns all tokens (including white space), so text can be reconstructed with fidelity.
func Tokenize(r io.Reader) *TokenStream {
	t := newTokenizer(r)
	return NewTokenStream(t.next)
}

// TokenizeString tokenizes a string into a stream of tokens. Iterate through the stream by calling Scan() or Next().
//
// It returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func TokenizeString(s string) *TokenStream {
	return Tokenize(strings.NewReader(s))
}

type tokenizer struct {
	sc *iterators.Scanner
}

func newTokenizer(r io.Reader) *tokenizer {
	return &tokenizer{
		sc: words.NewScanner(r),
	}
}

// next returns the next token. Call until it returns nil.
func (t *tokenizer) next() (*Token, error) {
	if t.sc.Scan() {
		token := NewToken(t.sc.Text(), false)
		return token, nil
	}
	if err := t.sc.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}
