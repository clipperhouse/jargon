package jargon

// Filter is a structure for processing a stream of tokens
type Filter interface {
	Filter(*TokenStream) *TokenStream
}
