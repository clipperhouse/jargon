package jargon

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/html"
)

// Tokenize returns an 'iterator' of Tokens from a io.Reader. Call .Next() until it returns nil:
//	tokens := Tokenize(reader)
//	for {
//		token := tokens.Next()
//		if token == nil {
//			break
//		}
//
// 		// do stuff with token
//	}
//
// The tokenizer is targeted to English text that contains tech terms, so things like C++ and .Net are handled as single units.
// It respects Unicode definitions of 'punctuation' and 'symbol', with a few hard-coded exceptions. Symbols are treated as word characters (like alphanumerics), allowing things like email addresses, hashtags and @-handles to be understood as a single token.
// Dots are allowed to lead words, and to appear mid-word, allowing things like .Net and Node.js. Trailing dots are considered end-of-sentence.
// [@, #, -, *, %, /, \] are hard-coded as symbols, even though Unicode specifies them as punctuation. See http://www.unicode.org/faq/punctuation_symbols.html
// All other punctuation terminates words, as does white space.
//
// It returns all tokens (including white space), so text can be reconstructed with fidelity. Ignoring (say) whitespace is a decision for the caller.
func Tokenize(r io.Reader) *TextTokens {
	return newTextTokens(r)
}

// TextTokens is an "iterator" for the results of lemmatization; keep calling .Next() until it returns nil, indicating the end
// TextTokens implements the Tokens interface
type TextTokens struct {
	incoming *bufio.Reader
	buffer   bytes.Buffer
}

func newTextTokens(r io.Reader) *TextTokens {
	return &TextTokens{
		incoming: bufio.NewReaderSize(r, 4*4096),
	}
}

// Next returns the next token; nil indicates end of data
func (t *TextTokens) Next() *Token {
	if t == nil {
		return nil
	}
	if t.buffer.Len() > 0 {
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
			// Problem
			panic(err)
		case unicode.IsSpace(r):
			t.accept(r)
			return t.token()
		case isPunct(r):
			t.accept(r)
			isLeadingPunct := mightBeLeadingPunct(r) && !t.peekTerminator()
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
func (t *TextTokens) readWord() *Token {
	for {
		r, _, err := t.incoming.ReadRune()
		switch {
		case err == io.EOF:
			// Current word is terminated by EOF, send it back
			return t.token()
		case mightBeMidPunct(r):
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

func (t *TextTokens) token() *Token {
	b := t.buffer.Bytes()
	if len(b) == 0 { // eof
		return nil
	}

	// Got the bytes, can reset
	t.buffer.Reset()

	// Determine punct and/or space
	if utf8.RuneCount(b) == 1 {
		// Punct and space are always one rune in our usage
		r, _ := utf8.DecodeRune(b)

		known, ok := knownTokens[r]

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
		punct: isPunct(r) || r == '\r' || r == '\n' || r == '\t',
		space: unicode.IsSpace(r),
	}
}

var knownTokens = make(map[rune]*Token)

func init() {
	runes := []rune{
		' ', '\r', '\n', '\t', '.', ',',
	}

	for _, r := range runes {
		knownTokens[r] = newTokenFromRune(r)
	}
}

func (t *TextTokens) accept(r rune) {
	t.buffer.WriteRune(r)
}

// PeekTerminator looks to the next rune and determines if it breaks a word
func (t *TextTokens) peekTerminator() bool {
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
// It returns a channel of Tokens, intended to be ranged over thus:
//	tokens := TokenizeHTML(string)
//	for t := range tokens {
// 		// do stuff
//	}
//
// It returns all tokens (including white space), so text can be reconstructed with fidelity. Ignoring (say) whitespace is a decision for the caller.
func TokenizeHTML(r io.Reader) *HTMLTokens {
	h := html.NewTokenizer(r)
	t := &HTMLTokens{h, nil}
	return t
}

type HTMLTokens struct {
	html *html.Tokenizer
	text *TextTokens
}

func (t *HTMLTokens) Next() *Token {
	// Are we "inside" a text node after previous call?
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
