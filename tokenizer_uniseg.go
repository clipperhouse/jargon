package jargon

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/blevesearch/segment"
)

// TokenizeUniseg tokenizes according to Unicode Text Segmentation https://unicode.org/reports/tr29/
//
// It returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil.
//
// Experimental; segmenter dependency has issues
func TokenizeUniseg(r io.Reader) *Tokens {
	t := newUnisegTokenizer(r)
	return &Tokens{
		Next: t.next,
	}
}

// TokenizeUnisegString tokenizes according to Unicode Text Segmentation https://unicode.org/reports/tr29/
//
// It returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil.
func TokenizeUnisegString(s string) *Tokens {
	return TokenizeUniseg(strings.NewReader(s))
}

type unisegTokenizer struct {
	segmenter *segment.Segmenter
	buffer    []seg
	outgoing  *TokenQueue
}

func newUnisegTokenizer(r io.Reader) *unisegTokenizer {
	return &unisegTokenizer{
		segmenter: segment.NewSegmenter(r),
		outgoing:  &TokenQueue{},
	}
}

type seg struct {
	Bytes []byte
	Type  int
	Err   error
}

func (seg seg) Is(typ int) bool {
	return seg.Type == typ
}

func (t *unisegTokenizer) segment() seg {
	return seg{
		Bytes: t.segmenter.Bytes(),
		Type:  t.segmenter.Type(),
		Err:   t.segmenter.Err(),
	}
}

// next returns the next token. Call until it returns nil.
func (t *unisegTokenizer) next() (*Token, error) {
	// First, look for something to send back
	if t.outgoing.Len() > 0 {
		return t.outgoing.Pop(), nil
	}

	// Else, pull new segment(s)
	for t.segmenter.Segment() {
		current := t.segment()

		if err := current.Err; err != nil {
			return nil, err
		}

	handle_current:

		// Something like a word or a grapheme?
		isWord := !current.Is(segment.None)
		if isWord {
			t.accept(current)

			// We continue to look for allowed middle and trailing chars, such as wishy-washy or C++
			continue
		}

		// At this point, it must be a rune (right?)
		// Guard statement, in case that's wrong
		r, ok := tryRuneInBytes(current.Bytes)
		if !ok {
			return nil, fmt.Errorf("should be a rune, but it's %q, this is likely a bug in the tokenizer", current)
		}

		if unicode.IsSpace(r) {
			// Space is always terminating

			// Anything in the buffer can go out
			t.emit()

			// Accept the space & emit it
			t.accept(current)
			t.emit()

			// We know everything in outgoing is a complete token
			return t.outgoing.Pop(), nil
		}

		// At this point, it's punct

		// Expressions like .Net, #hashtags and @handles
		// Must be one of our leading chars, and must be start of a new token
		mightBeLeading := len(t.buffer) == 0 && leadings[r]

		if mightBeLeading {
			// Look ahead
			if t.segmenter.Segment() {
				lookahead := t.segment()

				// Is it a word/grapheme?
				isLeading := !lookahead.Is(segment.None)
				if isLeading {
					// If so, we can concatenate in the buffer
					t.accept(current)
					t.accept(lookahead)

					// But we don't know if it's a complete token
					// Might be something like .net-core
					continue
				}

				// Else, consider it regular punctuation
				// Gotta handle the lookahead, we've consumed it

				// Current rune is not leading, therefore it is a token
				t.accept(current)
				t.emit()

				// Lookahead is space or punct, let the main loop handle it
				current = lookahead
				goto handle_current
			}
		}

		// Expressions like F# and C++
		// Must be one of our trailing chars, and must not be start of a new token
		mightBeTrailing := len(t.buffer) > 0 && trailings[r]

		if mightBeTrailing {
			// Look ahead
			if t.segmenter.Segment() {
				lookahead := t.segment()

				// May precede another (identical) trailing, like C++
				lr, ok := tryRuneInBytes(lookahead.Bytes)
				if ok && r == lr {
					// Append them both & emit
					t.accept(current)
					t.accept(lookahead)
					t.emit()
					continue
				}

				// Complete the token & queue it
				t.accept(current)
				t.emit()

				// Lookahead can be anything, let the main loop handle it
				current = lookahead
				goto handle_current
			}
		}

		// Truly terminating punct at this point

		// Queue the existing buffer
		t.emit()

		t.accept(current)
		t.emit()

		return t.outgoing.Pop(), nil
	}

	if err := t.segmenter.Err(); err != nil {
		return nil, err
	}

	// Anything left
	if len(t.buffer) > 0 {
		t.emit()
	}

	if t.outgoing.Len() > 0 {
		return t.outgoing.Pop(), nil
	}

	return nil, nil
}

func (t *unisegTokenizer) accept(s seg) {
	t.buffer = append(t.buffer, s)
}

func (t *unisegTokenizer) emit() {
	if len(t.buffer) > 0 {
		t.outgoing.Push(t.token())
	}
}

func (t *unisegTokenizer) token() *Token {
	// Avoid an allocation if possible
	if len(t.buffer) == 1 {
		b := t.buffer[0].Bytes
		// Got the bytes, can reset
		t.buffer = t.buffer[:0]
		return NewToken(string(b), false)
	}

	b := []byte{}

	for _, seg := range t.buffer {
		b = append(b, seg.Bytes...)
	}

	// Got the bytes, can reset
	t.buffer = t.buffer[:0]

	return NewToken(string(b), false)
}

func tryRuneInBytes(b []byte) (rune, bool) {
	ok := utf8.RuneCount(b) == 1

	if ok {
		r, _ := utf8.DecodeRune(b)
		return r, true
	}

	return utf8.RuneError, false
}

func tryRuneInString(s string) (rune, bool) {
	ok := utf8.RuneCountInString(s) == 1

	if ok {
		r, _ := utf8.DecodeRuneInString(s)
		return r, true
	}

	return utf8.RuneError, false
}

var leadingString = map[string]bool{
	".": true,
	"#": true,
	"@": true,
}

var leadings = runeSet{
	'.': true,
	'#': true,
	'@': true,
}

var trailings = runeSet{
	'+': true,
	'#': true,
}

// sketch of something simpler?
func (t *unisegTokenizer) next2() (*Token, error) {
	// First, look for something to send back
	if t.outgoing.Len() > 0 {
		return t.outgoing.Pop(), nil
	}

	seg := t.segmenter.Segment()
	if err := t.segmenter.Err(); err != nil {
		return nil, err
	}
	if !seg {
		return nil, nil
	}

	current := t.segmenter.Bytes()

	mightBeLeading := leadingString[string(current)]
	if mightBeLeading {
		seg := t.segmenter.Segment()
		if err := t.segmenter.Err(); err != nil {
			return nil, err
		}
		if !seg {
			// EOF
			goto emit
		}

		lookahead := t.segmenter.Bytes()
		typ := t.segmenter.Type()
		isLeading := typ == segment.Letter || typ == segment.Number

		if isLeading {
			// We have one token, concatenate current + lookahead
			current = append(current, lookahead...)
			goto emit
		}

		// Else, we have two tokens (current, lookahead), queue the lookahead
		ltoken := NewToken(string(lookahead), false)
		t.outgoing.Push(ltoken)
		goto emit
	}

emit:
	token := NewToken(string(current), false)
	return token, nil
}
