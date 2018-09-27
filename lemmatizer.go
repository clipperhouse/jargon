// Package jargon offers tokenizers and lemmatizers, for use in text processing and NLP
package jargon

import (
	"fmt"
	"strings"
	"unicode"
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

// Lemmatize transforms a tokens to their canonicalized terms.
// It returns an 'iterator' of Tokens, given input Tokens. Call .Next() until it returns nil:
//	tokens := Tokenize(reader)
//	lem := NewLemmatizer(stackexchange.Dictionary, 3)
//	lemmas := lem.Lemmatize(tokens)
//	for {
//		lemma := lemmas.Next()
//		if lemma == nil {
//			break
//		}
//
// 		// do stuff with lemma
//	}
// Tokens that are not canonicalized are returned as-is, e.g. for input:
//     "I", " ", "think", " ", "Ruby", " ", "on", " ", "Rails", " ", "is", " ", "great"
// lemmatized output:
//     "I", " ", "think", " ", "ruby-on-rails", " ", "is", " ", "great"
// Note that fewer tokens may be returned than were input, and that correct lemmatization depends on correct tokenization!
func (lem *Lemmatizer) Lemmatize(tokens Tokens) *LemmaTokens {
	return newLemmaTokens(lem, tokens)
}

// LemmaTokens is an "iterator" for the results of lemmatization; keep calling .Next() until it returns nil, indicating the end
// LemmaTokens implements the Tokens interface
type LemmaTokens struct {
	lem      *Lemmatizer
	incoming Tokens
	buffer   []*Token // for incoming tokens; no guarantee they will be emitted
	outgoing []*Token
}

func newLemmaTokens(lem *Lemmatizer, incoming Tokens) *LemmaTokens {
	return &LemmaTokens{
		lem:      lem,
		incoming: incoming,
	}
}

// Next returns the next token; nil indicates end of data
func (t *LemmaTokens) Next() *Token {
	if t == nil {
		return nil
	}
	for {

		if len(t.outgoing) > 0 {
			return t.emit()
		}

		t.fill(1) // ok to ignore this error

		if len(t.buffer) == 0 {
			return nil
		}

		switch tok := t.buffer[0]; {
		case tok.IsPunct() || tok.IsSpace():
			// Emit it straight from the incoming buffer
			t.drop(1)
			return tok
		default:
			// Else it's a word
			t.ngrams()
		}
	}
}

func (t *LemmaTokens) ngrams() {
	// Try n-grams, longest to shortest (greedy)
	for take := t.lem.maxGramLength; take > 0; take-- {
		run, consumed, ok := t.wordrun(take)

		if !ok {
			continue // on to the next n-gram
		}

		canonical, found := t.lem.Lookup(run)

		if found {
			// the canonical can have space or punct, so we want to return separate tokens

			// optimization: check if tokenization is needed, avoid expense if not
			var tokenize bool
			for _, r := range canonical {
				if unicode.IsSpace(r) || isPunct(r) {
					tokenize = true
					break
				}
			}

			if tokenize {
				r := strings.NewReader(canonical)
				tokens := Tokenize(r)
				for {
					tok := tokens.Next()
					if tok == nil {
						break
					}
					tok.lemma = true
					t.stage(tok) //set it up to be emitted
				}
			} else {
				tok := &Token{
					value: canonical,
					space: false,
					punct: false,
					lemma: true,
				}
				t.stage(tok) //set it up to be emitted
			}

			t.drop(consumed) // discard the incoming tokens that comprised the lemma
			return
		}

		if take == 1 {
			// No n-grams, just emit
			original := t.buffer[0]
			t.stage(original) // set it up to be emitted
			t.drop(1)         // take it out of the buffer
			return
		}
	}
	err := fmt.Errorf("did not find a token. this should never happen")
	panic(err)
}

func join(tokens []*Token) string {
	joined := make([]string, 0)
	for _, t := range tokens {
		joined = append(joined, t.String())
	}
	return strings.Join(joined, "")
}

// drop (truncate) the first `n` elements of the buffer
// remember, a token being in the buffer does not imply that we will emit it
func (t *LemmaTokens) drop(n int) {
	copy(t.buffer, t.buffer[n:])
	t.buffer = t.buffer[:len(t.buffer)-n]
}

// ensure that the buffer contains at least `count` elements; returns false if channel is exhausted before achieving the count
func (t *LemmaTokens) fill(count int) bool {
	for count >= len(t.buffer) {
		token := t.incoming.Next()
		if token == nil {
			// EOF
			return false
		}
		t.buffer = append(t.buffer, token)
	}
	return true
}

func (t *LemmaTokens) stage(tok *Token) {
	t.outgoing = append(t.outgoing, tok)
}

func (t *LemmaTokens) emit() *Token {
	n := 1
	tok := t.outgoing[0]
	copy(t.outgoing, t.outgoing[n:])
	t.outgoing = t.outgoing[:len(t.outgoing)-n]
	return tok
}

// Analogous to tokens.Take(take) in Linq
func (t *LemmaTokens) wordrun(take int) ([]string, int, bool) {
	var (
		taken []string // the words
		count int      // tokens consumed, not necessarily equal to take
	)

	for len(taken) < take {
		ok := t.fill(count)
		if !ok {
			// Not enough (buffered) tokens to continue
			// So, a word run of length `take` is impossible
			return nil, 0, false
		}

		token := t.buffer[count]
		switch {
		case token.IsPunct():
			// Note: test for punct before space; newlines and tabs can be
			// considered both punct and space (depending on the tokenizer!)
			// and we want to treat them as breaking word runs.
			return nil, 0, false
		case token.IsSpace():
			// Ignore and continue
			count++
		default:
			// Found a word
			taken = append(taken, token.String())
			count++
		}
	}

	return taken, count, true
}
