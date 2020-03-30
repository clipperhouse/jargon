package jargon

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/clipperhouse/jargon/is"
)

// Tokenize returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil.
//
// Its uses several specs from Unicode Text Segmentation https://unicode.org/reports/tr29/. It's not a full implementation, but a decent approximation for many mainstream cases.
//
// Tokenize returns all tokens (including white space), so text can be reconstructed with fidelity.
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
	outgoing *TokenQueue
}

func newTokenizer(r io.Reader) *tokenizer {
	return &tokenizer{
		incoming: bufio.NewReader(r),
		outgoing: &TokenQueue{},
	}
}

// next returns the next token. Call until it returns nil.
func (t *tokenizer) next() (*Token, error) {
	if t.outgoing.Any() {
		return t.outgoing.Pop(), nil
	}

	for {
		r, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return nil, nil
		case r == ' ', r == '\r', r == '\n', r == '\t':
			// An optimization to avoid hitting `is` methods
			token := NewToken(string(r), false)
			return token, nil
		case is.Leading(r):
			lookahead, eof, err := t.peekRune()
			if err != nil {
				return nil, err
			}
			if !eof && (is.ALetter(lookahead) || is.Numeric(lookahead)) {
				// It's leading
				t.accept(r)
				continue
			}
			// It's not leading
			token := NewToken(string(r), false)
			return token, nil
		case is.AHLetter(r):
			t.accept(r)
			return t.alphanumeric()
		case is.Numeric(r):
			t.accept(r)
			return t.numeric()
		case is.Katakana(r):
			t.accept(r)
			return t.katakana()
		default:
			// Everything else is its own token: punct, space, symbols, ideographs, controls, etc
			token := NewToken(string(r), false)
			return token, nil
		}
	}
}

func (t *tokenizer) alphanumeric() (*Token, error) {
	for {
		r, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return t.token(), nil
		case is.AHLetter(r):
			t.accept(r)
		case is.Numeric(r):
			t.accept(r)
		case is.MidLetter(r) || is.MidNumLetQ(r):
			// https://unicode.org/reports/tr29/#WB6 & WB7
			lookahead, eof, err := t.peekRune()
			if err != nil {
				return nil, err
			}

			if eof || !is.AHLetter(lookahead) {
				// r is trailing, not mid-word, and so a separate token; queue it for later
				trailing := NewToken(string(r), false)
				t.outgoing.Push(trailing)

				return t.token(), nil
			}

			// Otherwise, accept and continue
			t.accept(r)
		default:
			// Everything else is breaking

			// Current (breaking) rune is a token, queue it up for later
			breaking := NewToken(string(r), false)
			t.outgoing.Push(breaking)

			// Emit the buffered word
			return t.token(), nil
		}
	}
}

// https://unicode.org/reports/tr29/#WB11
func (t *tokenizer) numeric() (*Token, error) {
	for {
		r, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return t.token(), nil
		case is.Numeric(r):
			t.accept(r)
		case is.MidNum(r) || is.MidNumLetQ(r):
			lookahead, eof, err := t.peekRune()
			if err != nil {
				return nil, err
			}

			if eof || !is.Numeric(lookahead) {
				// r is trailing, not mid, and so a separate token; queue it for later
				trailing := NewToken(string(r), false)
				t.outgoing.Push(trailing)

				return t.token(), nil
			}

			// Otherwise, accept and continue
			t.accept(r)
		case is.AHLetter(r):
			t.accept(r)
			// Punt to general alpha
			return t.alphanumeric()
		default:
			return t.token(), nil
		}
	}
}

// https://unicode.org/reports/tr29/#WB13
func (t *tokenizer) katakana() (*Token, error) {
	for {
		r, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return t.token(), nil
		case is.Katakana(r):
			t.accept(r)
		default:
			return t.token(), nil
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

// readRune gets the next rune, advancing the reader
func (t *tokenizer) readRune() (r rune, eof bool, err error) {
	r, _, err = t.incoming.ReadRune()

	if err != nil {
		if err == io.EOF {
			return r, true, nil
		}
		return r, false, err
	}

	return r, false, nil
}

// peekRune peeks the next rune, without advancing the reader
func (t *tokenizer) peekRune() (r rune, eof bool, err error) {
	r, _, err = t.incoming.ReadRune()

	if err != nil {
		if err == io.EOF {
			return r, true, nil
		}
		return r, false, err
	}

	// Unread ASAP!
	err = t.incoming.UnreadRune()
	if err != nil {
		return r, false, err
	}

	return r, false, nil
}
