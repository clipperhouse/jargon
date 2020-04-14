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
		current, eof, err := t.readRune()
		switch {
		case err != nil:
			return nil, err
		case eof:
			return t.token(), nil
		}

		// Some funcs below require lookahead; better to do I/O here than there
		// (we don't care about eof for lookahead, irrelevant)
		lookahead, _, err := t.peekRune()
		if err != nil {
			return nil, err
		}

		switch {
		case t.wb3(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb3ab(current):
			// true indicates break
			token := t.token()
			token2 := NewToken(string(current), false)

			if token != nil {
				t.outgoing.Push(token2)
				return token, nil
			}
			return token2, nil
		case t.wb5(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb6(current, lookahead):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb7(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb7a(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb7b(current, lookahead):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb7c(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb8(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb9(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb10(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb11(current):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb12(current, lookahead):
			// true indicates continue
			t.accept(current)
			continue
		case t.wb13(current):
			// true indicates continue
			t.accept(current)
			continue
		}

		// https://unicode.org/reports/tr29/#WB999
		// Everything else is its own token: punct, space, symbols, ideographs, controls, etc

		token := t.token()
		token2 := NewToken(string(current), false)

		if token != nil {
			t.outgoing.Push(token2)
			return token, nil
		}
		return token2, nil
	}
}

// https://unicode.org/reports/tr29/#WB3
func (t *tokenizer) wb3(current rune) (continues bool) {
	// If it's a new token and CR
	if len(t.buffer) == 0 {
		return is.Cr(current)
	}

	// If it's LF and previous was CR
	if is.Lf(current) {
		previous := t.buffer[len(t.buffer)-1]
		return is.Cr(previous)
	}

	return false
}

// https://unicode.org/reports/tr29/#WB3a
func (t *tokenizer) wb3ab(current rune) (breaks bool) {
	return is.Cr(current) || is.Lf(current) || is.Newline(current)
}

// https://unicode.org/reports/tr29/#WB5
func (t *tokenizer) wb5(current rune) (continues bool) {
	// If it's a new token and AHLetter
	if len(t.buffer) == 0 {
		return is.AHLetter(current)
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.AHLetter(previous) && is.AHLetter(current)
}

// https://unicode.org/reports/tr29/#WB6
func (t *tokenizer) wb6(current, lookahead rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.AHLetter(previous) && (is.MidLetter(current) || is.MidNumLetQ(current)) && is.AHLetter(lookahead)
}

// https://unicode.org/reports/tr29/#WB7
func (t *tokenizer) wb7(current rune) (continues bool) {
	if len(t.buffer) < 2 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	preprevious := t.buffer[len(t.buffer)-2]

	return is.AHLetter(preprevious) && (is.MidLetter(previous) || is.MidNumLetQ(previous)) && is.AHLetter(current)
}

// https://unicode.org/reports/tr29/#WB7a
func (t *tokenizer) wb7a(current rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.HebrewLetter(previous) && is.SingleQuote(current)
}

// https://unicode.org/reports/tr29/#WB7b
func (t *tokenizer) wb7b(current, lookahead rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.AHLetter(previous) && is.DoubleQuote(current) && is.HebrewLetter(lookahead)
}

// https://unicode.org/reports/tr29/#WB7c
func (t *tokenizer) wb7c(current rune) (continues bool) {
	if len(t.buffer) < 2 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	preprevious := t.buffer[len(t.buffer)-2]

	return is.HebrewLetter(preprevious) && is.DoubleQuote(previous) && is.HebrewLetter(current)
}

// https://unicode.org/reports/tr29/#WB8
func (t *tokenizer) wb8(current rune) (continues bool) {
	// If it's a new token and Numeric
	if len(t.buffer) == 0 {
		return is.Numeric(current)
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.Numeric(previous) && is.Numeric(current)
}

// https://unicode.org/reports/tr29/#WB9
func (t *tokenizer) wb9(current rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.AHLetter(previous) && is.Numeric(current)
}

// https://unicode.org/reports/tr29/#WB9
func (t *tokenizer) wb10(current rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.Numeric(previous) && is.AHLetter(current)
}

// https://unicode.org/reports/tr29/#WB11
func (t *tokenizer) wb11(current rune) (continues bool) {
	if len(t.buffer) < 2 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	preprevious := t.buffer[len(t.buffer)-2]

	return is.Numeric(preprevious) && (is.MidNum(previous) || is.MidNumLetQ(previous)) && is.Numeric(current)
}

// https://unicode.org/reports/tr29/#WB12
func (t *tokenizer) wb12(current, lookahead rune) (continues bool) {
	if len(t.buffer) == 0 {
		return false
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.Numeric(previous) && (is.MidNum(current) || is.MidNumLetQ(current)) && is.Numeric(lookahead)
}

// https://unicode.org/reports/tr29/#WB13
func (t *tokenizer) wb13(current rune) (continues bool) {
	// If it's a new token and Katakana
	if len(t.buffer) == 0 {
		return is.Katakana(current)
	}

	previous := t.buffer[len(t.buffer)-1]
	return is.Katakana(previous) && is.Katakana(current)
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
