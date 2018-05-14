// Derived in part from https://golang.org/src/text/template/parse/lex.go

// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tokenizers

import (
	"unicode"
	"unicode/utf8"
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
	token := NewToken(value, punct, space)
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
		case r == '.' && l.last(): // final dot
			l.emit(true, false)
		case r == '.': // leading dot is ok
			return lexWord
		case isPunct(r):
			l.emit(true, false)
		case unicode.IsSpace(r):
			l.emit(false, true)
		default:
			return lexWord
		}
	}

	return nil
}

func lexWord(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == '.':
			// Look ahead to see if it's a leading dot
			if l.last() || isTerminator(l.peek()) {
				// It's a legit terminator, not leading or mid-word (like ".net" or "Node.js")
				// Emit the word
				l.backup()
				l.emit(false, false)

				// Dot gets emitted in lexMain
				break Loop
			}
			// Otherwise continue, it's a leading dot
		case l.last():
			// Always emit
			l.emit(false, false)
			break Loop
		case isTerminator(r):
			// Always emit
			l.backup()
			l.emit(false, false)
			break Loop
		case r == eof:
			break Loop
		}

		// Otherwise absorb and continue
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
