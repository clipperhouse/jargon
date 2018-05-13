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
// [@, #,-, *, %, /, \] hard-coded as symbols, even though Unicode specifies them as punctuation. See http://www.unicode.org/faq/punctuation_symbols.html
// All other punctuation terminates words, as does white space.
// TODO: URLs
// TODO: mid-word apostrophes?
var TechProse = &techProse{}

func (t *techProse) Tokenize(text string) []string {
	lex := lex(text)
	return lex.items
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	input string   // the string being scanned
	state stateFn  // the next lexing function to enter
	pos   int      // current position in the input
	start int      // start position of this item
	width int      // width of last rune read from input
	items []string // channel of scanned items
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
func (l *lexer) emit() {
	l.items = append(l.items, l.input[l.start:l.pos])
	l.start = l.pos
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make([]string, 0),
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
		case r == '.': // leading dot is ok
			return lexWord
		case isPunct(r):
			l.emit()
		case unicode.IsSpace(r):
			l.emit()
		default:
			return lexWord
		}
	}

	return nil
}

func lexWord(l *lexer) stateFn {
	for {
		r := l.next()

		if r == '.' {
			// Look ahead
			if l.last() || isTerminator(l.peek()) {
				// It's a legit terminator, not mid-word
				l.backup()
				l.emit()
				return lexMain
			}
			continue // skip following checks
		}

		if isTerminator(r) {
			// Always emit
			l.backup()
			l.emit()
			return lexMain
		}

		if r == eof {
			return lexMain
		}

		// Otherwise absorb and continue
	}
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
	'/':  exists,
	'\\': exists,
}

func isPunctException(r rune) bool {
	_, ok := punctExceptions[r]
	return ok
}
