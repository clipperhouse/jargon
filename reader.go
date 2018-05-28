package jargon

import (
	"bufio"
	"io"
	"unicode"
)

type reader struct {
	*bufio.Reader
	buffer []rune
	tokens chan Token
	state  state
}

// PeekTerminator looks to the next rune and determines if it breaks a word
func (b *reader) PeekTerminator() bool {
	r, _, err := b.ReadRune()

	if err != nil {
		if err == io.EOF {
			return true
		}
		panic(err)
	}

	// Unread ASAP!
	if uerr := b.UnreadRune(); uerr != nil {
		panic(uerr)
	}

	return isPunct(r) || unicode.IsSpace(r)
}

type state func(*reader) state

func newReader(r io.Reader) *reader {
	b := &reader{
		Reader: bufio.NewReader(r),
		tokens: make(chan Token, 0),
	}
	go b.run()
	return b
}

func (b *reader) run() {
	for b.state = readMain; b.state != nil; {
		b.state = b.state(b)
	}
	close(b.tokens)
}

func (b *reader) accept(r rune) {
	b.buffer = append(b.buffer, r)
}

func (b *reader) emit() {
	value := string(b.buffer)
	token := Token{
		value: value,
	}

	// Determine punct and/or space
	if len(b.buffer) == 1 {
		// Punct and space are always one rune in our usage
		r := b.buffer[0]

		// For our purposes, newlines and tabs should be considered punctuation, i.e.,
		// they break a word run. Lemmatizers should test for punct *before* testing for space.
		token.punct = isPunct(r) || r == '\r' || r == '\n' || r == '\t'
		token.space = unicode.IsSpace(r)
	}

	b.tokens <- token
	b.buffer = nil
}

func readMain(b *reader) state {
	for {
		switch r, _, err := b.ReadRune(); {
		case err != nil:
			if err == io.EOF {
				// No problem
				return nil
			}
			// Problem
			panic(err)
		case isPunct(r) || unicode.IsSpace(r):
			if mightBeLeadingPunct(r) {
				// Look to the next character
				if !b.PeekTerminator() {
					// Treat it as a word
					b.accept(r)
					return readWord
				}
			}
			b.accept(r)
			b.emit()
		default:
			b.accept(r)
			return readWord
		}
	}
}

// Important that this function only gets entered from the lexMain loop; lexMain determines 'word start'
func readWord(b *reader) state {
	for {
		r, _, err := b.ReadRune()
		switch {
		case err == io.EOF:
			b.emit()
			return readMain
		case mightBeMidPunct(r):
			// Look ahead to see if dot or apostrophe is acting as punctuation,
			// by being the last char, or being followed by space or more punctuation.
			if b.PeekTerminator() {
				// Emit the word
				b.emit()

				// Accept and emit the punct
				// (We'd rather UnreadRune here but PeekTerminator invalidates that.)
				b.accept(r)
				b.emit()

				return readMain
			}
			// Otherwise accept, it's a mid-word dot or apostrophe
		case isPunct(r) || unicode.IsSpace(r):
			// Emit the word
			b.emit()

			// Terminator will be handled by readMain
			b.UnreadRune()
			return readMain
		}

		// Otherwise accept and continue
		b.accept(r)
	}
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

func mightBeLeadingPunct(r rune) bool {
	_, ok := leadingPunct[r]
	return ok
}

var midPunct = map[rune]struct{}{
	// Punctuation that can appear mid-word
	'.':  exists,
	'\'': exists,
	'’':  exists,
}

func mightBeMidPunct(r rune) bool {
	_, ok := midPunct[r]
	return ok
}
