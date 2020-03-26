package jargon

import (
	"io"
	"strings"
)

// Tokens represents an 'iterator' of Token, the result of a call to Tokenize or Lemmatize. Call Next() until it returns nil.
type Tokens struct {
	// Next returns the next Token. If nil, the iterator is exhausted. Because it depends on I/O, callers should check errors.
	Next func() (*Token, error)
}

// ToSlice converts the Tokens iterator into a slice (array). Calling ToSlice will exhaust the iterator. For big files, putting everything into an array may cause memory pressure.
func (incoming *Tokens) ToSlice() ([]*Token, error) {
	var result []*Token

	for {
		t, err := incoming.Next()
		if err != nil {
			return result, err
		}
		if t == nil {
			break
		}
		result = append(result, t)
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

	for {
		t, err := incoming.Next()
		if err != nil {
			return b.String(), err
		}
		if t == nil {
			break
		}
		b.WriteString(t.String())
	}

	return b.String(), nil
}

// WriteTo writes all token string values to w
func (incoming *Tokens) WriteTo(w io.Writer) (int64, error) {
	var written int64
	for {
		t, err := incoming.Next()
		if err != nil {
			return written, err
		}
		if t == nil {
			break
		}

		n, err := w.Write([]byte(t.String()))
		written += int64(n)

		if err != nil {
			return written, err
		}
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
	for {
		t, err := incoming.Next()
		if err != nil {
			return 0, err
		}
		if t == nil {
			break
		}
		count++
	}
	return count, nil
}
