package jargon

// TokenFilter is a structure for processing a stream of tokens
type TokenFilter interface {
	Lookup([]string) (string, bool)
	MaxGramLength() int
}
