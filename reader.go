package jargon

import (
	"bufio"
	"io"
	"unicode"
)

type reader struct {
	*bufio.Reader
	buffer []rune
	tokens []Token
	state  state
}

// PeekTerminator looks to the next rune and determines if it breaks a word
func (b *reader) PeekTerminator() bool {
	r, _, err := b.ReadRune()

	// Unread immediately!
	if uerr := b.UnreadRune(); uerr != nil {
		panic(uerr)
	}

	return err != nil || isTerminator(r)
}

type state func(*reader) state

func newReader(r io.Reader) *reader {
	b := &reader{
		Reader: bufio.NewReader(r),
		buffer: make([]rune, 0),
		tokens: make([]Token, 0),
	}
	b.run()
	return b
}

func (b *reader) run() {
	for b.state = readMain; b.state != nil; {
		b.state = b.state(b)
	}
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
		// they break a word run. Lemmatizers should test for punct before testing for space.
		token.punct = isPunct(r) || r == '\r' || r == '\n' || r == '\t'
		token.space = unicode.IsSpace(r)
	}
	b.tokens = append(b.tokens, token)
	b.buffer = make([]rune, 0)
}

func readMain(b *reader) state {
	for {
		r, _, err := b.ReadRune()
		switch {
		case err == io.EOF:
			return nil
		case mightBeLeadingPunct(r):
			// Look to the next character
			if !b.PeekTerminator() {
				// Treat it as a word
				b.accept(r)
				return readWord
			}
			// Not leading punct, just regular punct
			b.accept(r)
			b.emit()
		case isPunct(r):
			b.accept(r)
			b.emit()
		case unicode.IsSpace(r):
			// For our purposes, newlines and tabs should be considered punctuation, i.e.,
			// they break a word run. Lemmatizers should test for punct before testing for space.
			b.accept(r)
			b.emit()
		default:
			// Punt to readWord
			b.UnreadRune()
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
		case isTerminator(r):
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
