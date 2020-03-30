package jargon

import (
	"unicode"
	"unicode/utf8"
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

// NewToken creates a new token, and calculates whether the token is space or punct.
func NewToken(s string, isLemma bool) *Token {
	token, found := common[s][isLemma]

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

func tryRuneInString(s string) (rune, bool) {
	ok := utf8.RuneCountInString(s) == 1

	if ok {
		r, _ := utf8.DecodeRuneInString(s)
		return r, true
	}

	return utf8.RuneError, false
}

var common = make(map[string]map[bool]*Token)

func init() {
	ss := []string{
		" ", "\r", "\n", "\t", ".", ",",
		"A", "a",
		"An", "an",
		"The", "the",
		"And", "and",
		"Or", "or",
		"Not", "not",
		"Of", "of",
		"In", "in",
		"On", "on",
		"To", "to",
		"Be", "be",
		"Is", "is",
		"Are", "are",
		"Has", "has",
		"Have", "have",
		"It", "it",
		"Do", "do",
	}

	for _, s := range ss {
		common[s] = map[bool]*Token{
			true:  NewToken(s, true),
			false: NewToken(s, false),
		}
	}
}
