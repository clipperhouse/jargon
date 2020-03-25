package jargon

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/blevesearch/segment"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Tokenize returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil:
//
// The tokenizer is targeted to English text that contains tech terms, so things like C++ and .Net are handled as single units, as are #hashtags and @handles.
//
// It generally relies on Unicode definitions of 'punctuation' and 'symbol', with some exceptions.
//
// Tokenize returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func Tokenize(r io.Reader) *Tokens {
	t := newTokenizer(r)
	return &Tokens{
		Next: t.next,
	}
}

// TokenizeString returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil:
//
// The tokenizer is targeted to English text that contains tech terms, so things like C++ and .Net are handled as single units, as are #hashtags and @handles.
//
// It generally relies on Unicode definitions of 'punctuation' and 'symbol', with some exceptions.
//
// Tokenize returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func TokenizeString(s string) *Tokens {
	return Tokenize(strings.NewReader(s))
}

type tokenizer struct {
	segmenter *segment.Segmenter
	buffer    []seg
	outgoing  *TokenQueue
}

func newTokenizer(r io.Reader) *tokenizer {
	return &tokenizer{
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

func (t *tokenizer) segment() seg {
	return seg{
		Bytes: t.segmenter.Bytes(),
		Type:  t.segmenter.Type(),
		Err:   t.segmenter.Err(),
	}
}
func (t *tokenizer) next2() (*Token, error) {
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

// next returns the next token. Call until it returns nil.
func (t *tokenizer) next() (*Token, error) {
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

func (t *tokenizer) accept(s seg) {
	t.buffer = append(t.buffer, s)
}

func (t *tokenizer) emit() {
	if len(t.buffer) > 0 {
		t.outgoing.Push(t.token())
	}
}

func (t *tokenizer) token() *Token {
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

// TokenizeHTML tokenizes HTML. Text nodes are tokenized using jargon.Tokenize; everything else (tags, comments) are left verbatim.
// It returns a Tokens, intended to be iterated over by calling Next(), until nil.
// It returns all tokens (including white space), so text can be reconstructed with fidelity. Ignoring (say) whitespace is a decision for the caller.
func TokenizeHTML(r io.Reader) *Tokens {
	t := &htokenizer{
		html: html.NewTokenizer(r),
		text: dummy, // dummy to avoid nil
	}
	return &Tokens{
		Next: t.next,
	}
}

var dummy = &Tokens{Next: func() (*Token, error) { return nil, nil }}

type htokenizer struct {
	html                    *html.Tokenizer
	text                    *Tokens
	suspendTextTokenization bool
	parent                  atom.Atom
}

// next is the implementation of the Tokens interface. To iterate, call until it returns nil
func (t *htokenizer) next() (*Token, error) {
	// Are we "inside" a text node?
	text, err := t.text.Next()
	if err != nil {
		return nil, err
	}
	if text != nil {
		return text, nil
	}

	for {
		tt := t.html.Next()

		if tt == html.ErrorToken {
			err := t.html.Err()
			if err == io.EOF {
				// No problem
				return nil, nil
			}
			return nil, err
		}

		htoken := t.html.Token()

		switch htoken.Type {
		case html.StartTagToken:
			// Record that we are entering script or style blocks; don't tokenize text
			if htoken.DataAtom == atom.Script || htoken.DataAtom == atom.Style {
				t.parent = htoken.DataAtom
			}
		case html.TextToken:
			switch t.parent {
			case atom.Script, atom.Style:
				// Don't tokenize script and style blocks, just return as one big string
				token := &Token{
					value: htoken.String(),
					punct: false,
					space: false,
				}
				return token, nil
			default:
				t.text = TokenizeString(htoken.String())
				return t.text.Next()
			}
		case html.EndTagToken:
			if htoken.DataAtom == t.parent {
				htoken.DataAtom = 0
			}
		}

		// Everything else is punct for our purposes
		token := &Token{
			value: htoken.String(),
			punct: true,
			space: false,
		}
		return token, nil
	}
}
