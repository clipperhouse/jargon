package jargon

import (
	"bufio"
	"io"
	"strings"

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
	err      error
	buffer   []rune
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
			return t.token(), nil
		case t.wb3(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb3ab(r):
			// true indicates break
			token := t.token()
			token2 := NewToken(string(r), false)

			if token != nil {
				t.outgoing.Push(token2)
				return token, nil
			}
			return token2, nil
		case t.wb5(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb6(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb7(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb7a(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb7b(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb7c(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb8(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb9(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb10(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb11(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb12(r):
			// true indicates continue
			t.accept(r)
			continue
		case t.wb13(r):
			// true indicates continue
			t.accept(r)
			continue
		}

		if t.err != nil {
			err := t.err
			t.err = nil
			return nil, err
		}

		// https://unicode.org/reports/tr29/#WB999
		// Everything else is its own token: punct, space, symbols, ideographs, controls, etc

		token := t.token()
		token2 := NewToken(string(r), false)

		if token != nil {
			t.outgoing.Push(token2)
			return token, nil
		}
		return token2, nil
	}
}

// https://unicode.org/reports/tr29/#WB3
func (t *tokenizer) wb3(r rune) (continues bool) {
	// If it's a new token and CR
	if len(t.buffer) == 0 {
		return is.Cr(r)
	}

	// If it's LF and previous was CR
	if is.Lf(r) {
		previous := t.buffer[len(t.buffer)-1]
		return is.Cr(previous)
	}

	return false
}

// https://unicode.org/reports/tr29/#WB3a
func (t *tokenizer) wb3ab(r rune) (breaks bool) {
	return is.Cr(r) || is.Lf(r) || is.Newline(r)
}

// https://unicode.org/reports/tr29/#WB5
func (t *tokenizer) wb5(r rune) (continues bool) {
	// If it's a new token and AHLetter
	if len(t.buffer) == 0 {
		return is.AHLetter(r)
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.AHLetter(previous) && is.AHLetter(r)
}

// https://unicode.org/reports/tr29/#WB6
func (t *tokenizer) wb6(r rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	lookahead, eof, err := t.peekRune()
	if err != nil {
		t.err = err
		return false
	}
	if eof {
		return false
	}

	return is.AHLetter(previous) && (is.MidLetter(r) || is.MidNumLetQ(r)) && is.AHLetter(lookahead)
}

// https://unicode.org/reports/tr29/#WB7
func (t *tokenizer) wb7(r rune) (continues bool) {
	if len(t.buffer) < 2 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	preprevious := t.buffer[len(t.buffer)-2]

	return is.AHLetter(preprevious) && (is.MidLetter(previous) || is.MidNumLetQ(previous)) && is.AHLetter(r)
}

// https://unicode.org/reports/tr29/#WB7a
func (t *tokenizer) wb7a(r rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.HebrewLetter(previous) && is.SingleQuote(r)
}

// https://unicode.org/reports/tr29/#WB7b
func (t *tokenizer) wb7b(r rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	lookahead, eof, err := t.peekRune()
	if err != nil {
		t.err = err
		return false
	}
	if eof {
		return false
	}

	return is.AHLetter(previous) && is.DoubleQuote(r) && is.HebrewLetter(lookahead)
}

// https://unicode.org/reports/tr29/#WB7c
func (t *tokenizer) wb7c(r rune) (continues bool) {
	if len(t.buffer) < 2 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	preprevious := t.buffer[len(t.buffer)-2]

	return is.HebrewLetter(preprevious) && is.DoubleQuote(previous) && is.HebrewLetter(r)
}

// https://unicode.org/reports/tr29/#WB8
func (t *tokenizer) wb8(r rune) (continues bool) {
	// If it's a new token and Numeric
	if len(t.buffer) == 0 {
		return is.Numeric(r)
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.Numeric(previous) && is.Numeric(r)
}

// https://unicode.org/reports/tr29/#WB9
func (t *tokenizer) wb9(r rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.AHLetter(previous) && is.Numeric(r)
}

// https://unicode.org/reports/tr29/#WB9
func (t *tokenizer) wb10(r rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.Numeric(previous) && is.AHLetter(r)
}

// https://unicode.org/reports/tr29/#WB11
func (t *tokenizer) wb11(r rune) (continues bool) {
	if len(t.buffer) < 2 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	preprevious := t.buffer[len(t.buffer)-2]

	return is.Numeric(preprevious) && (is.MidNum(previous) || is.MidNumLetQ(previous)) && is.Numeric(r)
}

// https://unicode.org/reports/tr29/#WB12
func (t *tokenizer) wb12(r rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	lookahead, eof, err := t.peekRune()
	if err != nil {
		t.err = err
		return false
	}
	if eof {
		return false
	}

	return is.Numeric(previous) && (is.MidNum(r) || is.MidNumLetQ(r)) && is.Numeric(lookahead)
}

// https://unicode.org/reports/tr29/#WB13
func (t *tokenizer) wb13(r rune) (continues bool) {
	// If it's a new token and Katakana
	if len(t.buffer) == 0 {
		return is.Katakana(r)
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.Katakana(previous) && is.Katakana(r)
}

func (t *tokenizer) token() *Token {
	if len(t.buffer) == 0 {
		return nil
	}

	s := string(t.buffer)
	t.buffer = t.buffer[:0]

	token := NewToken(s, false)
	return token
}

func (t *tokenizer) accept(r rune) {
	t.buffer = append(t.buffer, r)
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

func (t *tokenizer) unreadRune() error {
	err := t.incoming.UnreadRune()

	if err != nil {
		return err
	}

	return nil
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
