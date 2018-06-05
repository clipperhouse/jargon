// Package jargon offers tokenizers and lemmatizers, for use in text processing and NLP
package jargon

import (
	"strings"
)

// Lemmatizer is the main structure for looking up canonical tags
type Lemmatizer struct {
	values        map[string]string
	maxGramLength int
	normalize     func(string) string
}

// NewLemmatizer creates and populates a new Lemmatizer for the purpose of looking up canonical tags.
// Data and rules mostly live in the Dictionary interface, which is usually imported.
func NewLemmatizer(d Dictionary) *Lemmatizer {
	lem := &Lemmatizer{
		values:        make(map[string]string),
		maxGramLength: d.MaxGramLength(),
		normalize:     d.Normalize,
	}
	tags := d.Lemmas()
	for _, tag := range tags {
		key := lem.normalize(tag)
		lem.values[key] = tag
	}
	synonyms := d.Synonyms()
	for synonym, canonical := range synonyms {
		key := lem.normalize(synonym)
		lem.values[key] = canonical
	}
	return lem
}

// Lemmatize transforms a stream of tokens to their canonicalized terms.
// Tokens that are not canonicalized are returned as-is, e.g.
//     "I", " ", "think", " ", "Ruby", " ", "on", " ", "Rails", " ", "is", " ", "great"
// becomes
//     "I", " ", "think", " ", "ruby-on-rails", " ", "is", " ", "great"
// Note that fewer tokens may be returned than were input, and that correct lemmatization depends on correct tokenization!
func (lem *Lemmatizer) Lemmatize(tokens chan Token) chan Token {
	sc := newScanner(tokens)
	go lem.run(sc)
	return sc.outgoing
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
	close(sc.outgoing)
}

func (lem *Lemmatizer) ngrams(sc *scanner) {
	// Try n-grams, longest to shortest (greedy)
	for take := lem.maxGramLength; take > 0; take-- {
		run, consumed, ok := sc.wordrun(take)

		if !ok {
			continue // on to the next n-gram
		}

		gram := join(run)
		key := lem.normalize(gram)
		canonical, found := lem.values[key]

		if found {
			// Emit new token, replacing consumed tokens
			lemma := Token{
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

func join(tokens []Token) string {
	joined := make([]string, 0)
	for _, t := range tokens {
		joined = append(joined, t.String())
	}
	return strings.Join(joined, "")
}

type scanner struct {
	incoming, outgoing chan Token
	buffer             []Token
}

func newScanner(tokens chan Token) *scanner {
	return &scanner{
		incoming: tokens,
		outgoing: make(chan Token, 0),
		buffer:   make([]Token, 0),
	}
}
func (sc *scanner) emit(t Token) {
	sc.outgoing <- t
}

// drop (truncate) the first `n` elements of the buffer
// remember, a token being in the buffer does not imply that we will emit it
func (sc *scanner) drop(n int) {
	sc.buffer = sc.buffer[n:]
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
func (sc *scanner) wordrun(take int) ([]Token, int, bool) {
	taken := make([]Token, 0)
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
