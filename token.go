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

func (t Token) String() string {
	return t.value
}

// IsPunct indicates that the token should
func (t Token) IsPunct() bool {
	return t.punct
}

func (t Token) IsSpace() bool {
	return t.space
}

func Join(tokens []Token) string {
	joined := make([]string, 0)
	for _, t := range tokens {
		joined = append(joined, t.String())
	}
	return strings.Join(joined, "")
}
