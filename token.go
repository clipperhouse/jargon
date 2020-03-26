package jargon

import (
	"unicode"
)

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
func (t *Token) IsLemma() bool {
	return t.lemma
}

func NewToken(s string, isLemma bool) *Token {
	token, found := common[s]

	if found {
		return token
	}

	r, ok := tryRuneInString(s)

	return &Token{
		value: s,
		punct: ok && isPunct(r),
		space: ok && unicode.IsSpace(r),
		lemma: isLemma,
	}
}

var common = make(map[string]*Token)

func init() {
	runes := []string{
		" ", "\r", "\n", "\t", ".", ",",
	}

	for _, r := range runes {
		common[r] = NewToken(r, false)
	}
}
