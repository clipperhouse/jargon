package sigil

import (
	"fmt"

	"github.com/clipperhouse/jargon"
)

// NewFilter creates a new filter for leading characters. sigil is the leading character; legal defines legality for the following token.
func NewFilter(sigil string, legal func(s string) bool) jargon.Filter {
	f := &filter{
		sigil: sigil,
		legal: legal,
	}
	return f.filter
}

type filter struct {
	sigil string
	legal func(s string) bool
}

func (f *filter) filter(incoming *jargon.TokenStream) *jargon.TokenStream {
	t := &stream{
		filter:   f,
		incoming: incoming,
		outgoing: &jargon.TokenQueue{},
	}
	return jargon.NewTokenStream(t.next)
}

type stream struct {
	filter *filter

	incoming *jargon.TokenStream
	previous *jargon.Token
	outgoing *jargon.TokenQueue
}

func (s *stream) next() (*jargon.Token, error) {
	if s.outgoing.Any() {
		return s.outgoing.Pop(), nil
	}

	current, err := s.incoming.Next()
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil
	}

	// Previous token must not be a word
	boundaryOK := s.previous == nil || s.previous.IsSpace() || s.previous.IsPunct()
	if !boundaryOK {
		// Just send it back
		s.previous = current
		return current, nil
	}

	success, handle, err := s.try(s.filter.sigil, current, s.filter.legal)
	if err != nil {
		return nil, err
	}
	if success {
		// There should be nothing in outgoing
		if s.outgoing.Any() {
			return nil, fmt.Errorf("there should be nothing in outgoing, got %s", s.outgoing)
		}
		s.previous = handle
		return handle, nil
	}

	s.previous = current
	return current, nil
}

func (s *stream) try(sigil string, current *jargon.Token, legal func(string) bool) (bool, *jargon.Token, error) {
	if current.String() != sigil {
		return false, nil, nil
	}

	lookahead, err := s.incoming.Next()
	if err != nil {
		return false, nil, err
	}
	if lookahead == nil {
		// EOF
		return false, nil, nil
	}

	if legal(lookahead.String()) {
		// Drop current & lookahead, replace with new token
		s := sigil + lookahead.String()
		token := jargon.NewToken(s, true)
		return true, token, nil
	}

	// Queue the lookahead for later
	s.outgoing.Push(lookahead)

	return false, nil, nil
}
