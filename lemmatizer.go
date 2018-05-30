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
		buffer:        make([]Token, 0),
		tokens:        make(chan Token, 0),
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
	go func() {
		for {
			current, ok := <-tokens
			if ok {
				lem.buffer = append(lem.buffer, current)
			}
			if len(lem.buffer) == 0 {
				break
			}

			switch t := lem.buffer[0]; {
			case t.IsPunct() || t.IsSpace():
				// Emit it
				lem.tokens <- t
				lem.buffer = lem.buffer[1:]
			default:
				// Else it's a word, try n-grams, longest to shortest (greedy)
				for take := lem.maxGramLength; take > 0; take-- {
					run, consumed, ok := lem.wordrun(tokens, take)
					if ok {
						gram := Join(run)
						key := lem.normalize(gram)
						canonical, found := lem.values[key]

						if found {
							// Emit token, replacing consumed tokens
							lemma := Token{
								value: canonical,
								space: false,
								punct: false,
								lemma: true,
							}
							lem.tokens <- lemma
							// Discard the incoming tokens that comprise the lemma
							lem.buffer = lem.buffer[consumed:]
							break
						}

						if take == 1 {
							// No n-grams, just emit
							lem.tokens <- lem.buffer[0]
							lem.buffer = lem.buffer[1:]
						}
					}
				}
			}
		}
		lem.buffer = nil
		close(lem.tokens)
	}()

	return lem.tokens
}

// Analogous to tokens.Take(take) in Linq
func (lem *Lemmatizer) wordrun(tokens chan Token, take int) ([]Token, int, bool) {
	taken := make([]Token, 0)
	count := 0 // tokens consumed, not necessarily equal to take

	for len(taken) < take {
		for count >= len(lem.buffer) {
			token, ok := <-tokens
			if !ok {
				// Incoming token channel is closed
				// Hard stop
				return nil, 0, false
			}
			lem.buffer = append(lem.buffer, token)
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
