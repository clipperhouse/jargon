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

// String is the string value of the token
func (t Token) String() string {
	return t.value
}

// IsPunct indicates that the token should be considered 'breaking' of a run of words; a delimiter.
func (t Token) IsPunct() bool {
	return t.punct
}

// IsSpace indicates that the token consists entirely of white space (as defined by the unicode package)
func (t Token) IsSpace() bool {
	return t.space
}

// Join reconstructs a slice of tokens into their original string (assuming the tokens preserved fidelity of the original, esp white space!)
func Join(tokens []Token) string {
	joined := make([]string, 0)
	for _, t := range tokens {
		joined = append(joined, t.String())
	}
	return strings.Join(joined, "")
}
