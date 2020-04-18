package jargon

import (
	"io"
	"strings"
)

// Filter processes a stream of tokens
type Filter func(*TokenStream) *TokenStream

// TokenStream represents an 'iterator' of Token, the result of a call to Tokenize or Filter. Call Next() until it returns nil.
type TokenStream struct {
	next func() (*Token, error)

	token *Token // stateful token when using Scan
	err   error  // stateful error when using Scan
}

// Next returns the next Token. If nil, the iterator is exhausted. Because it depends on I/O, callers should check errors.
func (stream *TokenStream) Next() (*Token, error) {
	return stream.next()
}

// NewTokenStream creates a new TokenStream
func NewTokenStream(next func() (*Token, error)) *TokenStream {
	stream := &TokenStream{}

	// Shim to ensure consistency with Scan
	wrapper := func() (*Token, error) {
		stream.token, stream.err = next()
		return stream.Token(), stream.Err()
	}

	stream.next = wrapper
	return stream
}

// Scan retrieves the next token and returns true if successful. The resulting token can be retrieved using
// the Token() method. Scan returns false at EOF or on error. Be sure to check the Err() method.
//	for stream.Scan() {
//		token := stream.Token()
//		// do stuff with token
//	}
//	if err := stream.Err(); err != nil {
//		// do something with err
//	}
func (stream *TokenStream) Scan() bool {
	stream.token, stream.err = stream.Next()
	return stream.token != nil && stream.err == nil
}

// Token returns the current Token in the stream, after calling Scan
func (stream *TokenStream) Token() *Token {
	return stream.token
}

// Err returns the current error in the stream, after calling Scan
func (stream *TokenStream) Err() error {
	return stream.err
}

type where struct {
	stream    *TokenStream
	predicate func(*Token) bool
}

// Where filters a stream of Tokens that match a predicate
func (stream *TokenStream) Where(predicate func(*Token) bool) *TokenStream {
	w := &where{
		stream:    stream,
		predicate: predicate,
	}
	return NewTokenStream(w.next)
}

func (w *where) next() (*Token, error) {
	for {
		token, err := w.stream.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			break
		}

		if w.predicate(token) {
			return token, nil
		}
	}

	return nil, nil
}

// ToSlice converts the Tokens iterator into a slice (array). Calling ToSlice will exhaust the iterator. For big files, putting everything into an array may cause memory pressure.
func (stream *TokenStream) ToSlice() ([]*Token, error) {
	var result []*Token

	for stream.Scan() {
		result = append(result, stream.Token())
	}

	if err := stream.Err(); err != nil {
		return nil, stream.Err()
	}

	return result, nil
}

// Filter applies one or more filters to a token stream
func (stream *TokenStream) Filter(filters ...Filter) *TokenStream {
	outgoing := stream
	for _, f := range filters {
		outgoing = f(outgoing)
	}
	return outgoing
}

func (stream *TokenStream) String() (string, error) {
	var b strings.Builder

	for stream.Scan() {
		token := stream.Token()
		b.WriteString(token.String())
	}

	if err := stream.Err(); err != nil {
		return "", stream.Err()
	}

	return b.String(), nil
}

// WriteTo writes all token string values to w
func (stream *TokenStream) WriteTo(w io.Writer) (int64, error) {
	var written int64
	for stream.Scan() {
		token := stream.Token()
		n, err := w.Write([]byte(token.String()))
		written += int64(n)

		if err != nil {
			return written, err
		}
	}

	if err := stream.Err(); err != nil {
		return written, err
	}

	return written, nil
}

// Words returns only all non-punctuation and non-space tokens
func (stream *TokenStream) Words() *TokenStream {
	isWord := func(t *Token) bool {
		return !t.IsPunct() && !t.IsSpace()
	}
	w := &where{
		stream:    stream,
		predicate: isWord,
	}
	return NewTokenStream(w.next)
}

// Lemmas returns only tokens which have been 'lemmatized', or in some way modified by a token filter
func (stream *TokenStream) Lemmas() *TokenStream {
	w := &where{
		stream:    stream,
		predicate: (*Token).IsLemma,
	}
	return NewTokenStream(w.next)
}

// Distinct return one token per occurence of a given value (string)
func (stream *TokenStream) Distinct() *TokenStream {
	seen := map[string]bool{}
	isDistinct := func(t *Token) bool {
		found := seen[t.String()]
		if !found {
			seen[t.String()] = true
		}
		return !found
	}

	w := &where{
		stream:    stream,
		predicate: isDistinct,
	}
	return NewTokenStream(w.next)
}

// Count counts all tokens. Note that it will consume all tokens, so you will not be able to iterate further after making this call.
func (stream *TokenStream) Count() (int, error) {
	var count int
	for stream.Scan() {
		count++
	}

	if err := stream.Err(); err != nil {
		return count, stream.Err()
	}

	return count, nil
}
