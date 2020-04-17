// Package mapper provides a convenience builder for filters that map inputs to outputs, one-to-one
package mapper

import (
	"github.com/clipperhouse/jargon"
)

type filter struct {
	funcs []func(*jargon.Token) *jargon.Token
}

// NewFilter creates a filter which applies one or more funcs to each token
func NewFilter(funcs ...func(*jargon.Token) *jargon.Token) jargon.Filter {
	// Save the parameters for lazy loading (below)
	f := &filter{
		funcs: funcs,
	}
	return f.Filter
}

// Filter applies mapping func(s) to each incoming token
func (f *filter) Filter(incoming *jargon.TokenStream) *jargon.TokenStream {
	t := &tokens{
		incoming: incoming,
		filter:   f,
	}

	return jargon.NewTokenStream(t.next)
}

type tokens struct {
	// incoming stream of tokens from another source, such as a tokenizer
	incoming *jargon.TokenStream
	filter   *filter
}

func (t *tokens) next() (*jargon.Token, error) {
	token, err := t.incoming.Next()
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, nil
	}

	for _, f := range t.filter.funcs {
		token = f(token)
	}

	return token, nil
}
