package jargon

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/html"
)

// Tokenize returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil:
//	tokens := jargon.Tokenize(reader)
//	for {
//		token := tokens.Next()
//		if token == nil {
//			break
//		}
//
// 		// do stuff with token
//	}
//
// The tokenizer is targeted to English text that contains tech terms, so things like C++ and .Net are handled as single units, as are #hashtags and @handles.
//
// It generally relies on Unicode definitions of 'punctuation' and 'symbol', with some exceptions.
//
// Unlike some other common tokenizers, Tokenize returns all tokens (including white space), so text can be reconstructed with fidelity ("round tripped").
func Tokenize(r io.Reader) Tokens {
	t := newTokenizer(r)
	return Tokens{
		Next: t.next,
	}
}

type tokenizer struct {
	incoming *bufio.Reader
	outgoing bytes.Buffer
}

func newTokenizer(r io.Reader) *tokenizer {
	return &tokenizer{
		incoming: bufio.NewReaderSize(r, 4*4096),
	}
}

// next returns the next token. Call until it returns nil.
func (t *tokenizer) next() *Token {
	if t == nil {
		return nil
	}
	if t.outgoing.Len() > 0 {
		// Punct or space accepted in previous call to readWord
		return t.token()
	}
	for {
		switch r, _, err := t.incoming.ReadRune(); {
		case err != nil:
			if err == io.EOF {
				// No problem
				return t.token()
			}
			log.Print(err)
			return nil
		case unicode.IsSpace(r):
			t.accept(r)
			return t.token()
		case isPunct(r):
			t.accept(r)
			isLeadingPunct := leadingPunct.includes(r) && !t.peekTerminator()
			if isLeadingPunct {
				// Treat it as start of a word
				return t.readWord()
			}
			// Regular punct, emit it
			return t.token()
		default:
			// It's a letter
			t.accept(r)
			return t.readWord()
		}
	}
}

// Important that this function only gets entered from the Next() loop, which determines 'word start'
func (t *tokenizer) readWord() *Token {
	for {
		r, _, err := t.incoming.ReadRune()
		switch {
		case err == io.EOF:
			// Current word is terminated by EOF, send it back
			return t.token()
		case midPunct.includes(r):
			// Look ahead to see if it's followed by space or more punctuation
			followedByTerminator := t.peekTerminator()
			if followedByTerminator {
				// It's just regular punct, treat it as such

				// Get the current word token without the punct
				tok := t.token()

				// Accept the punct for later
				t.accept(r)

				// Emit the word token
				return tok
			}
			// Else, it's mid-word punct, treat it like a letter
			t.accept(r)
		case isPunct(r) || unicode.IsSpace(r):
			// Get the current word token without the punct
			tok := t.token()

			// Accept the punct for later
			t.accept(r)

			// Emit the word token
			return tok
		default:
			// Otherwise it's a letter, keep going
			t.accept(r)
		}
	}
}

func (t *tokenizer) token() *Token {
	b := t.outgoing.Bytes()
	if len(b) == 0 { // eof
		return nil
	}

	// Got the bytes, can reset
	t.outgoing.Reset()

	// Determine punct and/or space
	if utf8.RuneCount(b) == 1 {
		// Punct and space are always one rune in our usage
		r, _ := utf8.DecodeRune(b)

		known, ok := common[r]

		if ok {
			return known
		}

		return newTokenFromRune(r)
	}

	return &Token{
		value: string(b),
	}
}

func newTokenFromRune(r rune) *Token {
	return &Token{
		value: string(r),
		punct: isPunct(r),
		space: unicode.IsSpace(r),
	}
}

var common = make(map[rune]*Token)

func init() {
	runes := []rune{
		' ', '\r', '\n', '\t', '.', ',',
	}

	for _, r := range runes {
		common[r] = newTokenFromRune(r)
	}
}

func (t *tokenizer) accept(r rune) {
	t.outgoing.WriteRune(r)
}

// PeekTerminator looks to the next rune and determines if it breaks a word
func (t *tokenizer) peekTerminator() bool {
	r, _, err := t.incoming.ReadRune()

	if err != nil {
		if err == io.EOF {
			return true
		}
		panic(err)
	}

	// Unread ASAP!
	if uerr := t.incoming.UnreadRune(); uerr != nil {
		panic(uerr)
	}

	return isPunct(r) || unicode.IsSpace(r)
}

// TokenizeHTML tokenizes HTML. Text nodes are tokenized using jargon.Tokenize; everything else (tags, comments) are left verbatim.
// It returns a Tokens, intended to be iterated over by calling Next(), until nil
//	tokens := jargon.TokenizeHTML(reader)
//	for {
//		tok := tokens.Next()
//		if tok == nil {
//			break
//		}
//		// Do stuff with tok...
//	}
// It returns all tokens (including white space), so text can be reconstructed with fidelity. Ignoring (say) whitespace is a decision for the caller.
func TokenizeHTML(r io.Reader) Tokens {
	t := &htokenizer{
		html: html.NewTokenizer(r),
		text: dummy, // dummy to avoid nil
	}
	return Tokens{
		Next: t.Next,
	}
}

var dummy = Tokenize(strings.NewReader(""))

type htokenizer struct {
	html *html.Tokenizer
	text Tokens
}

// Next is the implementation of the Tokens interface. To iterate, call until it returns nil
func (t *htokenizer) Next() *Token {
	// Are we "inside" a text node?
	text := t.text.Next()
	if text != nil {
		return text
	}

	for {
		tt := t.html.Next()

		if tt == html.ErrorToken {
			// Presumably eof
			return nil
		}

		switch tok := t.html.Token(); {
		case tok.Type == html.TextToken:
			r := strings.NewReader(tok.Data)
			t.text = Tokenize(r)
			return t.text.Next()
		default:
			// Everything else is punct for our purposes
			new := &Token{
				value: tok.String(),
				punct: true,
				space: false,
			}
			return new
		}
	}
}
