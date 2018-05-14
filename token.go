package jargon

import (
	"strings"
)

type Token struct {
	value string
	punct bool
	space bool
}

func NewToken(value string, punct, space bool) Token {
	return Token{value, punct, space}
}

func (t Token) Value() string {
	return t.value
}

func (t Token) String() string {
	return t.value
}

func (t Token) Punct() bool {
	return t.punct
}

func (t Token) Space() bool {
	return t.space
}

func Join(tokens []Token, f func(Token) string) string {
	joined := make([]string, 0)
	for _, t := range tokens {
		joined = append(joined, t.Value())
	}
	return strings.Join(joined, "")
}
