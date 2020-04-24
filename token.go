package jargon

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

// Token represents a piece of text with metadata.
type Token struct {
	value               []byte
	punct, space, lemma bool
}

// Bytes returned the bytes of the token
func (t *Token) Bytes() []byte {
	return t.value
}

// String is the string value of the token
func (t *Token) String() string {
	return string(t.value)
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
func NewToken(s []byte, isLemma bool) *Token {
	token, found := common[string(s)][isLemma]

	if found {
		return token
	}

	if len(s) == 0 {
		return nil
	}

	var punct, space bool

	switch {
	case reflect.DeepEqual(s, []byte("\r\n")):
		punct = true
		space = true
	default:
		punct = true
		i := 0
		for i < len(s) {
			r, w := utf8.DecodeRune(s[i:])
			i += w
			if !isPunct(r) {
				punct = false
				break
			}
		}

		space = true
		i = 0
		for i < len(s) {
			r, w := utf8.DecodeRune(s[i:])
			i += w
			if !unicode.IsSpace(r) {
				space = false
				break
			}
		}
	}

	return &Token{
		value: s,
		punct: punct,
		space: space,
		lemma: isLemma,
	}
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
			true:  NewToken([]byte(s), true),
			false: NewToken([]byte(s), false),
		}
	}
}
