package jargon

import (
	"io"
	"strings"
)

// Tokens represents an 'iterator' of Token, the result of a call to Tokenize or Filter. Call Next() until it returns nil.
type Tokens struct {
	// Next returns the next Token. If nil, the iterator is exhausted. Because it depends on I/O, callers should check errors.
	Next func() (*Token, error)

	token *Token // stateful token when using Scan
	err   error  // stateful error when using Scan
}

func newTokens(next func() (*Token, error)) *Tokens {
	return &Tokens{
		Next: next,
	}
}

func (t *Tokens) Scan() bool {
	t.token, t.err = t.Next()
	return t.token != nil && t.err == nil
}

func (t *Tokens) Token() *Token {
	return t.token
}

func (t *Tokens) Err() error {
	return t.err
}

// ToSlice converts the Tokens iterator into a slice (array). Calling ToSlice will exhaust the iterator. For big files, putting everything into an array may cause memory pressure.
func (incoming *Tokens) ToSlice() ([]*Token, error) {
	var result []*Token

	for incoming.Scan() {
		result = append(result, incoming.Token())
	}

	if err := incoming.Err(); err != nil {
		return nil, incoming.Err()
	}

	return result, nil
}

// Filter applies one or more filters to a token stream
func (incoming *Tokens) Filter(filters ...Filter) *Tokens {
	outgoing := incoming
	for _, f := range filters {
		outgoing = f.Filter(outgoing)
	}
	return outgoing
}

func (incoming *Tokens) String() (string, error) {
	var b strings.Builder

	for incoming.Scan() {
		token := incoming.Token()
		b.WriteString(token.String())
	}

	if err := incoming.Err(); err != nil {
		return "", incoming.Err()
	}

	return b.String(), nil
}

// WriteTo writes all token string values to w
func (incoming *Tokens) WriteTo(w io.Writer) (int64, error) {
	var written int64
	for incoming.Scan() {
		token := incoming.Token()
		n, err := w.Write([]byte(token.String()))
		written += int64(n)

		if err != nil {
			return written, err
		}
	}

	if err := incoming.Err(); err != nil {
		return written, incoming.Err()
	}

	return written, nil
}

// Words returns only all non-punctuation and non-space tokens
func (incoming *Tokens) Words() *Tokens {
	isWord := func(t *Token) bool {
		return !t.IsPunct() && !t.IsSpace()
	}
	w := &where{
		incoming:  incoming,
		predicate: isWord,
	}
	return &Tokens{
		Next: w.next,
	}
}

// Lemmas returns only tokens which have been 'lemmatized', or in some way modified by a token filter
func (incoming *Tokens) Lemmas() *Tokens {
	w := &where{
		incoming:  incoming,
		predicate: (*Token).IsLemma,
	}
	return &Tokens{
		Next: w.next,
	}
}

// Count counts all tokens. Note that it will consume all tokens, so you will not be able to iterate further after making this call.
func (incoming *Tokens) Count() (int, error) {
	var count int
	for incoming.Scan() {
		count++
	}

	if err := incoming.Err(); err != nil {
		return count, incoming.Err()
	}

	return count, nil
}
