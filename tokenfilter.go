package jargon

// TokenFilter is a structure for processing a stream of tokens
type TokenFilter interface {
	Lookup(...string) (string, bool)
	MaxGramLength() int
}

// Filter is a structure for processing a stream of tokens
type Filter interface {
	Filter(*Tokens) *Tokens
}
