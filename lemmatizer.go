// Package jargon offers tokenizers and lemmatizers, for use in text processing and NLP
package jargon

import (
	"io"
	"strings"
)

// Lemmatizer is the main structure for looking up canonical tags
type Lemmatizer struct {
	Dictionary
	maxGramLength int
}

// NewLemmatizer creates and populates a new Lemmatizer for the purpose of looking up and replacing canonical tags.
func NewLemmatizer(d Dictionary, maxGramLength int) *Lemmatizer {
	lem := &Lemmatizer{
		Dictionary:    d,
		maxGramLength: maxGramLength,
	}
	return lem
}

// Lemmatize transforms a stream of tokens to their canonicalized terms.
// Tokens that are not canonicalized are returned as-is, e.g.
//     "I", " ", "think", " ", "Ruby", " ", "on", " ", "Rails", " ", "is", " ", "great"
// becomes
//     "I", " ", "think", " ", "ruby-on-rails", " ", "is", " ", "great"
// Note that fewer tokens may be returned than were input, and that correct lemmatization depends on correct tokenization!
func (lem *Lemmatizer) Lemmatize(tokens <-chan *Token) <-chan *Token {
	outgoing := make(chan *Token, 0)
	emit := func(t *Token) {
		outgoing <- t
	}

	sc := newScanner(tokens, emit)
	go func() {
		// Need closure to close outgoing channel on completion of run()
		lem.run(sc)
		close(outgoing)
	}()

	return outgoing
}

// LemmatizeAndWrite transforms a stream of tokens to their canonicalized terms, and writes them to w.
// An error is returned if the writer returns an error
// Tokens that are not canonicalized are returned as-is, e.g.
//     "I", " ", "think", " ", "Ruby", " ", "on", " ", "Rails", " ", "is", " ", "great"
// becomes
//     "I", " ", "think", " ", "ruby-on-rails", " ", "is", " ", "great"
// Note that fewer tokens may be returned than were input, and that correct lemmatization depends on correct tokenization!
func (lem *Lemmatizer) LemmatizeAndWrite(tokens <-chan *Token, w io.Writer) error {
	errchan := make(chan error, 0)
	emit := func(t *Token) {
		b := []byte(t.String())
		_, err := w.Write(b)
		if err != nil {
			errchan <- err
		}
	}

	sc := newScanner(tokens, emit)
	go func() {
		// Need closure to close error channel on completion of run()
		lem.run(sc)
		close(errchan)
	}()

	for err := range errchan {
		return err
	}

	return nil
}

func (lem *Lemmatizer) run(sc *scanner) {
	for {
		sc.fill(1) // ok to ignore this error

		if len(sc.buffer) == 0 {
			break
		}

		switch t := sc.buffer[0]; {
		case t.IsPunct() || t.IsSpace():
			// Emit it
			sc.emit(t)
			sc.drop(1)
		default:
			// Else it's a word
			lem.ngrams(sc)
		}
	}
	sc.buffer = nil
}

func (lem *Lemmatizer) ngrams(sc *scanner) {
	// Try n-grams, longest to shortest (greedy)
	for take := lem.maxGramLength; take > 0; take-- {
		run, consumed, ok := sc.wordrun(take)

		if !ok {
			continue // on to the next n-gram
		}

		gram := join(run)
		canonical, found := lem.Lookup(gram)

		if found {
			// Emit new token, replacing consumed tokens
			lemma := &Token{
				value: canonical,
				space: false,
				punct: false,
				lemma: true,
			}
			sc.emit(lemma)
			sc.drop(consumed) // discard the incoming tokens that comprised the lemma
			break
		}

		if take == 1 {
			// No n-grams, just emit
			sc.emit(run[0])
			sc.drop(1)
		}
	}
}

func join(tokens []*Token) string {
	joined := make([]string, 0)
	for _, t := range tokens {
		joined = append(joined, t.String())
	}
	return strings.Join(joined, "")
}

type scanner struct {
	incoming <-chan *Token
	buffer   []*Token
	emit     func(*Token)
}

func newScanner(incoming <-chan *Token, emit func(*Token)) *scanner {
	return &scanner{
		incoming: incoming,
		emit:     emit,
	}
}

// drop (truncate) the first `n` elements of the buffer
// remember, a token being in the buffer does not imply that we will emit it
func (sc *scanner) drop(n int) {
	l := len(sc.buffer[n:])
	b := make([]*Token, l, 2*l)
	copy(b, sc.buffer[n:])
	sc.buffer = b
}

// ensure that the buffer contains at least `count` elements; returns false if channel is exhausted before achieving the count
func (sc *scanner) fill(count int) bool {
	for count >= len(sc.buffer) {
		token, ok := <-sc.incoming
		if !ok {
			// Incoming token channel is closed
			return false
		}
		sc.buffer = append(sc.buffer, token)
	}
	return true
}

// Analogous to tokens.Take(take) in Linq
func (sc *scanner) wordrun(take int) ([]*Token, int, bool) {
	taken := make([]*Token, 0)
	count := 0 // tokens consumed, not necessarily equal to take

	for len(taken) < take {
		ok := sc.fill(count)
		if !ok {
			// Not enough (buffered) tokens to continue
			// So, a word run of length `take` is impossible
			return nil, 0, false
		}

		token := sc.buffer[count]
		switch {
		// Note: test for punct before space; newlines and tabs can be
		// considered both punct and space (depending on the tokenizer!)
		// and we want to treat them as breaking word runs.
		case token.IsPunct():
			// Hard stop
			return nil, 0, false
		case token.IsSpace():
			// Ignore and continue
			count++
		default:
			// Found a word
			taken = append(taken, token)
			count++
		}
	}

	return taken, count, true
}
