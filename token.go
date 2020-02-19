package jargon

import "strings"

// Tokens represents an 'iterator' of Token. Call .Next() until it returns nil.
type Tokens struct {
	// Next returns the next Token. If nil, the iterator is exhausted.
	Next func() *Token
}

// ForEach iterates over all tokens and executes f. A convenience function, so you don't have to call Next and check nil. Call ForEach will exhaust the iterator.
func (tokens Tokens) ForEach(f func(t *Token)) {
	for {
		t := tokens.Next()
		if t == nil {
			break
		}
		f(t)
	}
}

// ToSlice converts the Tokens iterator into a slice (array). Calling ToSlice will exhaust the iterator. For big files, putting everything into an array may cause memory pressure.
func (tokens Tokens) ToSlice() []*Token {
	var result []*Token

	tokens.ForEach(func(t *Token) {
		result = append(result, t)
	})

	return result
}

func (tokens Tokens) String() string {
	var b strings.Builder

	tokens.ForEach(func(t *Token) {
		b.WriteString(t.value)
	})

	return b.String()
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
