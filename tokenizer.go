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
	buffer   bytes.Buffer
}

func newTokenizer(r io.Reader) *tokenizer {
	return &tokenizer{
		incoming: bufio.NewReader(r),
	}
}

func (t *tokenizer) mightBeLeading(r rune) bool {
	switch r {
	case
		'.',
		'_',
		'#',
		'@':
		return true
	}
	return false
}

func (t *tokenizer) isLeadingPunct(r rune) (bool, error) {
	if t.mightBeLeading(r) {
		followedByTerminator, err := t.lookaheadIsTerminator()
		if err != nil {
			return false, err
		}
		return !followedByTerminator, nil
	}
	return false, nil
}

func (t *tokenizer) mightBeMidPunct(r rune) bool {
	switch r {
	case
		'.',
		'_',
		'\'',
		'’',
		'/':
		return true
	}
	return false
}

func (t *tokenizer) isMidPunct(r rune) (bool, error) {
	if t.mightBeMidPunct(r) {
		terminated, err := t.lookaheadIsTerminator()
		if err != nil {
			return false, err
		}
		return !terminated, nil
	}
	return false, nil
}

func (t *tokenizer) mightBeTrailingPunct(r rune) bool {
	switch r {
	case
		'+',
		'#',
		'_':
		return true
	}
	return false
}

func (t *tokenizer) isTrailingPunct(r rune) (bool, error) {
	if t.mightBeTrailingPunct(r) {
		terminated, err := t.lookaheadIsTerminator()
		if err != nil {
			return false, err
		}
		return terminated, nil
	}
	return false, nil
}

// next returns the next token. Call until it returns nil.
func (t *tokenizer) next() (*Token, error) {
	if t.buffer.Len() > 0 {
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
			// no need to buffer it
			token := NewToken(string(r), false)
			return token, nil
		case isPunct(r):
			leading, err := t.isLeadingPunct(r)
			if err != nil {
				return nil, err
			}
			if leading {
				// Treat it as start of a word
				t.accept(r)
				return t.word()
			}

			// Regular punct, emit it (no need to buffer)
			token := NewToken(string(r), false)
			return token, nil
		default:
			// It's alphanumeric
			t.accept(r)
			return t.word()
		}
	}
}

// Important that this function only gets entered from the Next() loop, which determines 'word start'
func (t *tokenizer) word() (*Token, error) {
	for {
		r, _, err := t.incoming.ReadRune()
		switch {
		case err != nil:
			if err == io.EOF {
				// No problem
				return t.token(), nil
			}
			return nil, err
		case isPunct(r):
			mid, err := t.isMidPunct(r)
			if err != nil {
				return nil, err
			}
			if mid {
				// It's mid-word punct, treat it like a letter
				t.accept(r)
				continue
			}

			trailing, err := t.isTrailingPunct(r)
			if err != nil {
				return nil, err
			}
			if trailing {
				// It's trailing punct, treat it like a letter
				t.accept(r)
				continue
			}

			// It's just regular punct

			// Get the current word token without the punct
			token := t.token()

			// Accept the punct for later
			t.accept(r)

			// Emit the word token
			return token, nil
		case unicode.IsSpace(r):
			// Get the current word token without the punct
			token := t.token()

			// Accept the punct for later
			t.accept(r)

			// Emit the word token
			return token, nil
		default:
			// Otherwise it's alphanumeric, keep going
			t.accept(r)
		}
	}
}

func (t *tokenizer) token() *Token {
	b := t.buffer.Bytes()

	// Got the bytes, can reset
	t.buffer.Reset()

	return NewToken(string(b), false)
}

func (t *tokenizer) accept(r rune) {
	t.buffer.WriteRune(r)
}

// PeekTerminator looks to the next rune and determines if it breaks a word
func (t *tokenizer) lookaheadIsTerminator() (bool, error) {
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
