package jargon

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

// Tokenize returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil.
//
// Tokenize returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func Tokenize(r io.Reader) *Tokens {
	t := newTokenizer(r)
	return &Tokens{
		Next: t.next,
	}
}

// TokenizeString returns an 'iterator' of Tokens. Call .Next() until it returns nil.
//
// It returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func TokenizeString(s string) *Tokens {
	return Tokenize(strings.NewReader(s))
}

type tokenizer struct {
	incoming *bufio.Reader
	outgoing bytes.Buffer
}

func newTokenizer(r io.Reader) *tokenizer {
	return &tokenizer{
		incoming: bufio.NewReaderSize(r, 64*1024),
	}
}

// next returns the next token. Call until it returns nil.
func (t *tokenizer) next() (*Token, error) {
	if t.outgoing.Len() > 0 {
		// Punct or space accepted in previous call to readWord
		return t.token(), nil
	}
	for {
		switch r, _, err := t.incoming.ReadRune(); {
		case err != nil:
			if err == io.EOF {
				// No problem, we're done
				return nil, nil
			}
			return nil, err
		case unicode.IsSpace(r):
			t.accept(r)
			return t.token(), nil
		case isPunct(r):
			t.accept(r)

			followedByTerminator, err := t.peekTerminator()
			if err != nil {
				return nil, err
			}

			isLeadingPunct := isLeadingPunct(r) && !followedByTerminator
			if isLeadingPunct {
				// Treat it as start of a word
				return t.readWord()
			}
			// Regular punct, emit it
			return t.token(), nil
		default:
			// It's a letter
			t.accept(r)
			return t.readWord()
		}
	}
}

// Important that this function only gets entered from the Next() loop, which determines 'word start'
func (t *tokenizer) readWord() (*Token, error) {
	for {
		r, _, err := t.incoming.ReadRune()
		switch {
		case err != nil:
			if err == io.EOF {
				// No problem
				return t.token(), nil
			}
			return nil, err
		case isMidPunct(r):
			// Look ahead to see if it's followed by space or more punctuation
			followedByTerminator, err := t.peekTerminator()
			if err != nil {
				return nil, err
			}

			if followedByTerminator {
				// It's just regular punct, treat it as such

				// Get the current word token without the punct
				token := t.token()

				// Accept the punct for later
				t.accept(r)

				// Emit the word token
				return token, nil
			}

			// Else, it's mid-word punct, treat it like a letter
			t.accept(r)
		case isPunct(r) || unicode.IsSpace(r):
			// Get the current word token without the punct
			token := t.token()

			// Accept the punct for later
			t.accept(r)

			// Emit the word token
			return token, nil
		default:
			// Otherwise it's a letter, keep going
			t.accept(r)
		}
	}
}

func (t *tokenizer) token() *Token {
	b := t.outgoing.Bytes()

	// Got the bytes, can reset
	t.outgoing.Reset()

	return NewToken(string(b), false)
}

func (t *tokenizer) accept(r rune) {
	t.outgoing.WriteRune(r)
}

// PeekTerminator looks to the next rune and determines if it breaks a word
func (t *tokenizer) peekTerminator() (bool, error) {
	r, _, err := t.incoming.ReadRune()

	if err != nil {
		if err == io.EOF {
			return true, nil
		}
		return false, err
	}

	// Unread ASAP!
	if err := t.incoming.UnreadRune(); err != nil {
		return false, err
	}

	return isPunct(r) || unicode.IsSpace(r), nil
}
