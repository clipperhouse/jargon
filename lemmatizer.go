// Package jargon offers tokenizers and lemmatizers, for use in text processing and NLP
package jargon

import (
	"github.com/clipperhouse/jargon/stackexchange"
)

// Lemmatizer is the main structure for looking up canonical tags
type Lemmatizer struct {
	values        map[string]string
	maxGramLength int
	normalize     func(string) string
	buffer        []Token
	tokens        chan Token
}

// StackExchange is a built-in *Lemmatizer, using tag and synonym data from the following Stack Exchange sites: Stack Overflow,
// Server Fault, Game Dev and Data Science. It's indended to identify canonical tags (technologies),
// e.g. Ruby on Rails (3 words) will be replaced with ruby-on-rails (1 word).
// It looks for word runs (n-grams) up to length 3, ignoring spaces.
var StackExchange = NewLemmatizer(stackexchange.Dictionary)

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

// LemmatizeTokens takes a slice of tokens and returns tokens with canonicalized terms.
// Terms (tokens) that are not canonicalized are returned as-is, e.g.
//     ["I", " ", "think", " ", "Ruby", " ", "on", " ", "Rails", " ", "is", " ", "great"]
// becomes
//     ["I", " ", "think", " ", "ruby-on-rails", " ", "is", " ", "great"]
// Note that fewer tokens may be returned than were input.
// A lot depends on the original tokenization, so make sure that it's right!
func (lem *Lemmatizer) LemmatizeTokens(tokens chan Token) chan Token {
	lem.tokens = make(chan Token, 0)
	go lem.scan(tokens)
	return lem.tokens
}

func (lem *Lemmatizer) scan(tokens chan Token) {
	lem.buffer = make([]Token, 0)
	for {
		lem.fill(tokens, 1) // ok to ignore this error

		if len(lem.buffer) == 0 {
			break
		}

		switch t := lem.buffer[0]; {
		case t.IsPunct() || t.IsSpace():
			// Emit it
			lem.emit(t)
			lem.drop(1)
		default:
			// Else it's a word, try n-grams, longest to shortest (greedy)
			for take := lem.maxGramLength; take > 0; take-- {
				run, consumed, ok := lem.wordrun(tokens, take)
				if ok {
					gram := Join(run)
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
						lem.emit(lemma)
						lem.drop(consumed) // discard the incoming tokens that comprised the lemma
						break
					}

					if take == 1 {
						// No n-grams, just emit
						lem.emit(t)
						lem.drop(1)
					}
				}
			}
		}
	}
	lem.buffer = nil
	close(lem.tokens)
}

func (lem *Lemmatizer) emit(t Token) {
	lem.tokens <- t
}

// drop (truncate) the first `n` elements of the buffer
// remember, a token being in the buffer does not imply that we will emit it
func (lem *Lemmatizer) drop(n int) {
	lem.buffer = lem.buffer[n:]
}

// ensure that the buffer contains at least `count` elements; returns false if channel is exhausted before achieving the count
func (lem *Lemmatizer) fill(tokens chan Token, count int) bool {
	for count >= len(lem.buffer) {
		token, ok := <-tokens
		if !ok {
			// Incoming token channel is closed
			return false
		}
		lem.buffer = append(lem.buffer, token)
	}
	return true
}

// Analogous to tokens.Take(take) in Linq
func (lem *Lemmatizer) wordrun(tokens chan Token, take int) ([]Token, int, bool) {
	taken := make([]Token, 0)
	count := 0 // tokens consumed, not necessarily equal to take

	for len(taken) < take {
		ok := lem.fill(tokens, count)
		if !ok {
			// Not enough (buffered) tokens to continue
			// So, a word run of length `take` is impossible
			return nil, 0, false
		}

		token := lem.buffer[count]
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
