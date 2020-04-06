package twitter

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/clipperhouse/jargon"
)

// Filter looks for Twitter-style @handles and #hashtags
func Filter(incoming *jargon.TokenStream) *jargon.TokenStream {
	t := &tokens{
		incoming: incoming,
		outgoing: &jargon.TokenQueue{},
	}
	return jargon.NewTokenStream(t.next)
}

type tokens struct {
	incoming *jargon.TokenStream
	previous *jargon.Token
	outgoing *jargon.TokenQueue
}

func (t *tokens) next() (*jargon.Token, error) {
	if t.outgoing.Any() {
		return t.outgoing.Pop(), nil
	}

	current, err := t.incoming.Next()
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil
	}

	// Previous token must not be a word
	boundaryOK := t.previous == nil || t.previous.IsSpace() || t.previous.IsPunct()
	if !boundaryOK {
		// Just send it back
		t.previous = current
		return current, nil
	}

	success, handle, err := t.try("@", current, legalHandle)
	if err != nil {
		return nil, err
	}
	if success {
		// There should be nothing in outgoing
		if t.outgoing.Any() {
			return nil, fmt.Errorf("there should be nothing in outgoing, got %s", t.outgoing)
		}
		t.previous = handle
		return handle, nil
	}

	success, hashtag, err := t.try("#", current, legalHashtag)
	if err != nil {
		return nil, err
	}
	if success {
		// There should be nothing in outgoing, verify
		if t.outgoing.Any() {
			return nil, fmt.Errorf("there should be nothing in outgoing, got %s", t.outgoing)
		}
		t.previous = hashtag
		return hashtag, nil
	}

	t.previous = current
	return current, nil
}

func (t *tokens) try(sigil string, current *jargon.Token, legal func(string) bool) (bool, *jargon.Token, error) {
	if current.String() != sigil {
		return false, nil, nil
	}

	lookahead, err := t.incoming.Next()
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
	t.outgoing.Push(lookahead)

	return false, nil, nil
}

// https://help.twitter.com/en/managing-your-account/twitter-username-rules
func legalHandle(s string) bool {
	length := 0
	for _, r := range s {
		length++

		switch {
		case length > 15:
			return false
		case
			'A' <= r && r <= 'Z',
			'a' <= r && r <= 'r',
			'0' <= r && r <= '9',
			r == '_':
			continue
		default:
			return false
		}
	}

	return true
}

// Determined by playing with Twitter's web UI to see what got highlighted ;)
func legalHashtag(s string) bool {
	// One-character hashtags need to be a letter
	length := utf8.RuneCountInString(s)
	if length == 1 {
		r, _ := utf8.DecodeRuneInString(s)
		return unicode.IsLetter(r)
	}

	for _, r := range s {
		switch {
		case
			unicode.IsLetter(r),
			unicode.IsNumber(r),
			r == '_':
			continue
		default:
			return false
		}
	}

	return true
}
