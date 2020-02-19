package jargon

import (
	"io"
	"strings"
)

// Tokens represents an 'iterator' of Token, the result of a call to Tokenize or Lemmatize. Call Next() until it returns nil.
type Tokens struct {
	// Next returns the next Token. If nil, the iterator is exhausted.
	Next func() (*Token, error)
}

// ToSlice converts the Tokens iterator into a slice (array). Calling ToSlice will exhaust the iterator. For big files, putting everything into an array may cause memory pressure.
func (tokens Tokens) ToSlice() ([]*Token, error) {
	var result []*Token

	for {
		t, err := tokens.Next()
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

func (tokens Tokens) String() (string, error) {
	var b strings.Builder

	for {
		t, err := tokens.Next()
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
func (tokens Tokens) WriteTo(w io.Writer) (int64, error) {
	var written int64
	for {
		t, err := tokens.Next()
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

// Token represents a piece of text with metadata.
type Token struct {
	value               string
	punct, space, lemma bool
}

// String is the string value of the token
func (t *Token) String() string {
	return t.value
}

// IsPunct indicates that the token should be considered 'breaking' of a run of words. Mostly uses
// Unicode's definition of punctuation, with some exceptions for our purposes.
func (t *Token) IsPunct() bool {
	return t.punct
}

// IsSpace indicates that the token consists entirely of white space, as defined by the unicode package.
//
//A token can be both IsPunct and IsSpace -- for example, line breaks and tabs are punctuation for our purposes.
func (t *Token) IsSpace() bool {
	return t.space
}

// IsLemma indicates that the token is a lemma, i.e., a canonical term that replaced original token(s).
func (t Token) IsLemma() bool {
	return t.lemma
}
