// Derived in part from https://golang.org/src/text/template/parse/lex.go
// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jargon

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/html"
)

type techProse struct{}

// TechProse tokenizer attempts to tokenize English text that contains tech terms.
// It respects Unicode definitions of 'punctuation' and 'symbol', with a few hard-coded exceptions (below).
// Symbols are treated as word characters (like alphanumerics), allowing things like email addresses, hashtags and @-handles to be understood as a single token.
// Dots are allowed to lead words, and to appear mid-word, allowing things like .Net and Node.js. Trailing dots are considered end-of-sentence.
// [@, #, -, *, %, /, \] are hard-coded as symbols, even though Unicode specifies them as punctuation. See http://www.unicode.org/faq/punctuation_symbols.html
// All other punctuation terminates words, as does white space.
// Like the other tokenizers in this package, it returns all tokens (including white space), so text can be reconstructed with fidelity. If callers don't want white space, they'll need to filter.
// TODO: URLs
// TODO: mid-word apostrophes?
var TechProse = &techProse{}

func (t *techProse) Tokenize(text string) []Token {
	lex := lex(text)
	return lex.tokens
}

type techHTML struct{}

// TechHTML is a tokenizer for HTML text. Text nodes are tokenized using TechProse; tags and comments left verbatim.
var TechHTML = &techHTML{}

func (t *techHTML) Tokenize(text string) []Token {
	result := make([]Token, 0)
	r := strings.NewReader(text)
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			// Presumably eof
			break
		}

		switch tok := z.Token(); {
		case tok.Type == html.TextToken:
			words := TechProse.Tokenize(tok.Data)
			result = append(result, words...)
		default:
			// Everything else is punct for our purposes
			new := Token{
				value: tok.String(),
				punct: true,
				space: false,
			}
			result = append(result, new)
		}
	}

	return result
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	input  string  // the string being scanned
	state  stateFn // the next lexing function to enter
	pos    int     // current position in the input
	start  int     // start position of this item
	width  int     // width of last rune read from input
	tokens []Token // channel of scanned items
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) last() bool {
	return l.pos == len(l.input)
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(punct, space bool) {
	value := l.input[l.start:l.pos]
	token := Token{
		value: value,
		punct: punct,
		space: space,
	}
	l.tokens = append(l.tokens, token)
	l.start = l.pos
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make([]Token, 0),
	}
	l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexMain; l.state != nil; {
		l.state = l.state(l)
	}
}

func lexMain(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == eof:
			break Loop
		case isLeadingPunct(r) && l.last():
			// Final dot
			l.emit(true, false)
		case isLeadingPunct(r):
			// Leading dot might be ok
			return lexWord
		case isPunct(r):
			l.emit(true, false)
		case unicode.IsSpace(r):
			// For our purposes, newlines and tabs should be considered punctuation, i.e.,
			// they break a word run. Lemmatizers should test for punct before testing for space.
			punct := r == '\r' || r == '\n' || r == '\t'
			l.emit(punct, true)
		default:
			return lexWord
		}
	}

	return nil
}

// Important that this function only gets entered from the lexMain loop; lexMain determines 'word start'
func lexWord(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isLeadingPunct(r) || isMidPunct(r):
			// Could be a leading or mid-word dot,
			// or a mid-word apostrophe

			// Look ahead to see if dot or apostrophe is acting as punctuation,
			// by being the last char, or being followed by space or more punctuation.
			// (Test last before testing peek, peek will throw if eof)
			if l.last() || isTerminator(l.peek()) {
				// Emit the word
				l.backup()
				l.emit(false, false)

				// Dot gets emitted in lexMain
				break Loop
			}
			// Otherwise continue, it's a leading or mid-word dot, or mid-word apostrophe
		case l.last():
			// Always emit
			l.emit(false, false)
			break Loop
		case isTerminator(r):
			// Always emit
			// Terminator will be handled by lexMain
			l.backup()
			l.emit(false, false)
			break Loop
		default:
			// Otherwise absorb and continue
		}
	}
	return lexMain
}

func isTerminator(r rune) bool {
	return isPunct(r) || unicode.IsSpace(r)
}

func isPunct(r rune) bool {
	return unicode.IsPunct(r) && !isPunctException(r)
}

var exists = struct{}{}
var punctExceptions = map[rune]struct{}{
	// In some cases, we want to consider it a symbol, even though Unicode defines it as punctuation
	// See See http://www.unicode.org/faq/punctuation_symbols.html
	'-':  exists,
	'#':  exists,
	'@':  exists,
	'*':  exists,
	'%':  exists,
	'_':  exists,
	'/':  exists,
	'\\': exists,
}

func isPunctException(r rune) bool {
	_, ok := punctExceptions[r]
	return ok
}

var leadingPunct = map[rune]struct{}{
	// Punctuation that can lead a word, like .Net
	'.': exists,
}

func isLeadingPunct(r rune) bool {
	_, ok := leadingPunct[r]
	return ok
}

var midPunct = map[rune]struct{}{
	// Punctuation that can appear mid-word
	'.':  exists,
	'\'': exists,
	'’':  exists,
}

func isMidPunct(r rune) bool {
	_, ok := midPunct[r]
	return ok
}
