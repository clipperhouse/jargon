package jargon

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/clipperhouse/jargon/is"
)

// Tokenize2 returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil.
//
// Tokenize2 returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func Tokenize2(r io.Reader) *Tokens {
	t := newTokenizer2(r)
	return &Tokens{
		Next: t.next,
	}
}

// TokenizeString2 returns an 'iterator' of Tokens. Call .Next() until it returns nil.
//
// It returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func TokenizeString2(s string) *Tokens {
	return Tokenize2(strings.NewReader(s))
}

type tokenizer2 struct {
	incoming *bufio.Reader
	buffer   bytes.Buffer
	outgoing *TokenQueue
}

func newTokenizer2(r io.Reader) *tokenizer2 {
	return &tokenizer2{
		incoming: bufio.NewReader(r),
		outgoing: &TokenQueue{},
	}
}

// next returns the next token. Call until it returns nil.
func (t *tokenizer2) next() (*Token, error) {
	if t.outgoing.Any() {
		// Punct or space accepted in previous call to readWord
		return t.outgoing.Pop(), nil
	}

	for {
		r, _, err := t.incoming.ReadRune()
		switch {
		case err != nil:
			if err == io.EOF {
				// No problem, we're done
				return nil, nil
			}
			return nil, err
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

func (t *tokenizer2) alphanumeric() (*Token, error) {
	for {
		r, _, err := t.incoming.ReadRune()
		switch {
		case err != nil:
			if err == io.EOF {
				// No problem
				return t.token(), nil
			}
			return nil, err
		case is.AHLetter(r):
			t.accept(r)
		case is.Numeric(r):
			t.accept(r)
		case is.MidLetter(r) || is.MidNumLetQ(r):
			// https://unicode.org/reports/tr29/#WB6 & WB7
			lookahead, eof, err := t.lookahead()
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
func (t *tokenizer2) numeric() (*Token, error) {
	for {
		r, _, err := t.incoming.ReadRune()
		switch {
		case err != nil:
			if err == io.EOF {
				// No problem
				return t.token(), nil
			}
			return nil, err
		case is.Numeric(r):
			t.accept(r)
		case is.MidNum(r) || is.MidNumLetQ(r):
			lookahead, eof, err := t.lookahead()
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
func (t *tokenizer2) katakana() (*Token, error) {
	for {
		r, _, err := t.incoming.ReadRune()
		switch {
		case err != nil:
			if err == io.EOF {
				// No problem
				return t.token(), nil
			}
			return nil, err
		case is.Katakana(r):
			t.accept(r)
		default:
			return t.token(), nil
		}
	}
}

func (t *tokenizer2) token() *Token {
	b := t.buffer.Bytes()

	// Got the bytes, can reset
	t.buffer.Reset()

	return NewToken(string(b), false)
}

func (t *tokenizer2) accept(r rune) {
	t.buffer.WriteRune(r)
}

// lookahead peeks the next rune
func (t *tokenizer2) lookahead() (r rune, eof bool, err error) {
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
