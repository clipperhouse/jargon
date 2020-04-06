package jargon

// Filter processes a stream of tokens
type Filter func(*TokenStream) *TokenStream
