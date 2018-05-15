package jargon

import (
	"strings"
)

// Token represents a piece of text with metadata.
type Token struct {
	value string
	punct bool
	space bool
}

func (t Token) Value() string {
	return t.value
}

func (t Token) String() string {
	return t.value
}

func (t Token) IsPunct() bool {
	return t.punct
}

func (t Token) IsSpace() bool {
	return t.space
}

func Join(tokens []Token, f func(Token) string) string {
	joined := make([]string, 0)
	for _, t := range tokens {
		joined = append(joined, t.Value())
	}
	return strings.Join(joined, "")
}
