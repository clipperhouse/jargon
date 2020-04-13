package jargon

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/clipperhouse/jargon/is"
)

// Tokenize tokenizes a reader into a stream of tokens. Iterate through the stream by calling Scan() or Next().
//
// Its uses several specs from Unicode Text Segmentation https://unicode.org/reports/tr29/. It's not a full implementation, but a decent approximation for many mainstream cases.
//
// Tokenize returns all tokens (including white space), so text can be reconstructed with fidelity.
func Tokenize(r io.Reader) *TokenStream {
	t := newTokenizer(r, false)
	return NewTokenStream(t.next)
}

// TokenizeString tokenizes a string into a stream of tokens. Iterate through the stream by calling Scan() or Next().
//
// It returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func TokenizeString(s string) *TokenStream {
	return Tokenize(strings.NewReader(s))
}

type tokenizer struct {
	incoming *bufio.Reader
	buffer   bytes.Buffer
	outgoing *TokenQueue

	// guard is a debugging flag to verify assumptions (aka guard statements)
	guard bool
}

func newTokenizer(r io.Reader, guard bool) *tokenizer {
	return &tokenizer{
		incoming: bufio.NewReaderSize(r, 64*1024),
		outgoing: &TokenQueue{},
		guard:    guard,
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
		case r == ' ', r == '\t':
			// An optimization to avoid hitting `is` methods
			token := NewToken(string(r), false)
			return token, nil
		case is.Cr(r):
			// https://unicode.org/reports/tr29/#WB3
			t.accept(r)
			return t.cr()
		case r == '\r' || r == '\n' || is.Newline(r):
			// https://unicode.org/reports/tr29/#WB3a
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
		case is.HebrewLetter(r):
			t.accept(r)
			return t.hebrew()
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
	// Assumes an AHLetter is already in the buffer
	if t.guard {
		b := t.buffer.Bytes()
		if len(b) == 0 {
			return nil, fmt.Errorf(`buffer should be have one or more runes; %q; this is likely a bug in the tokenizer`, string(b))
		}
	}

	for {
		r, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return t.token()
		case is.HebrewLetter(r):
			t.accept(r)
			return t.hebrew()
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

				return t.token()
			}

			// Otherwise, accept and continue
			t.accept(r)
		default:
			// Everything else is breaking

			// Current (breaking) rune is a token, queue it up for later
			breaking := NewToken(string(r), false)
			t.outgoing.Push(breaking)

			// Emit the buffered word
			return t.token()
		}
	}
}

func (t *tokenizer) cr() (*Token, error) {
	// Assumes \r was previously accepted in main loop
	if t.guard {
		b := t.buffer.Bytes()
		if len(b) != 1 || b[0] != '\r' {
			return nil, fmt.Errorf(`buffer should be '\r'; this is likely a bug in the tokenizer`)
		}
	}

	lookahead, eof, err := t.peekRune()
	switch {
	case err != nil:
		return nil, err
	case eof:
		return t.token()
	case lookahead == '\n':
		// It's CRLF, which we want to be a single token
		t.accept(lookahead)

		// Act as if we read instead of peeked
		if _, err := t.incoming.Discard(utf8.RuneLen(lookahead)); err != nil {
			// Based on successful peek above, error should be impossible?
			return nil, err
		}

		return t.token()
	default:
		// CR is it's own token, then
		return t.token()
	}
}

// https://unicode.org/reports/tr29/#WB11
func (t *tokenizer) numeric() (*Token, error) {
	// Assumes an Number is already in the buffer
	if t.guard {
		b := t.buffer.Bytes()
		if len(b) == 0 {
			return nil, fmt.Errorf(`buffer should be have one or more runes; this is likely a bug in the tokenizer`)
		}
		last, _ := utf8.DecodeLastRune(b)
		if !is.Numeric(last) {
			return nil, fmt.Errorf(`last rune should be numeric; this is likely a bug in the tokenizer`)
		}
	}

	for {
		r, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return t.token()
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

				return t.token()
			}

			// Otherwise, accept and continue
			t.accept(r)
		case is.AHLetter(r):
			t.accept(r)
			// Punt to general alpha
			return t.alphanumeric()
		default:
			return t.token()
		}
	}
}

// https://unicode.org/reports/tr29/#WB13
func (t *tokenizer) katakana() (*Token, error) {
	// Assumes a Katakana character is already in the buffer
	if t.guard {
		b := t.buffer.Bytes()
		if len(b) == 0 {
			return nil, fmt.Errorf(`katakana: buffer should be have one or more runes; this is likely a bug in the tokenizer`)
		}
		last, _ := utf8.DecodeLastRune(b)
		if !is.Katakana(last) {
			return nil, fmt.Errorf(`last rune should be katakana; this is likely a bug in the tokenizer`)
		}
	}

	for {
		r, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return t.token()
		case is.Katakana(r):
			t.accept(r)
		default:
			return t.token()
		}
	}
}

// https://unicode.org/reports/tr29/#WB7a
func (t *tokenizer) hebrew() (*Token, error) {
	// Assumes a Hebrew character is already in the buffer
	if t.guard {
		b := t.buffer.Bytes()
		if len(b) == 0 {
			return nil, fmt.Errorf(`hebrew: buffer should be have one or more runes; this is likely a bug in the tokenizer`)
		}
		last, _ := utf8.DecodeLastRune(b)
		if !is.HebrewLetter(last) {
			return nil, fmt.Errorf(`last rune should be hebrew; this is likely a bug in the tokenizer`)
		}
	}

	for {
		r, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return t.token()
		case r == '\'':
			t.accept(r)
		case r == '"':
			lookahead, eof, err := t.peekRune()
			if err != nil {
				return nil, err
			}

			if eof || !is.HebrewLetter(lookahead) {
				// Terminate the token
				return t.token()
			}

			t.accept(r)
		case is.HebrewLetter(r):
			t.accept(r)
		case is.AHLetter(r):
			t.accept(r)
			return t.alphanumeric()
		default:
			return t.token()
		}
	}
}

func (t *tokenizer) token() (*Token, error) {
	b := t.buffer.Bytes()

	if len(b) == 0 {
		return nil, fmt.Errorf(`token: buffer should be have one of more runes; this is likely a bug in the tokenizer`)
	}

	// Got the bytes, can reset
	t.buffer.Reset()

	token := NewToken(string(b), false)
	return token, nil
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
