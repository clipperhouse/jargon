// Package jargon offers tokenizers and lemmatizers, for use in text processing and NLP
package jargon

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/clipperhouse/jargon/stackexchange"
)

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
func Lemmatize(incoming Tokens, dictionaries ...Dictionary) Tokens {
	if len(dictionaries) == 0 {
		dictionaries = append(dictionaries, stackexchange.Dictionary)
	}

	var result = incoming

	// new lemmatizer for each dictionary, results as input into next lemmatizer
	for _, dictionary := range dictionaries {
		lem := &lemmatizer{
			incoming:   result,
			dictionary: dictionary,
		}
		result = Tokens{
			Next: lem.next,
		}
	}

	return result
}

type lemmatizer struct {
	dictionary Dictionary
	incoming   Tokens
	buffer     []*Token // for incoming tokens; no guarantee they will be emitted
	outgoing   []*Token
}

// next returns the next token; nil indicates end of data
func (t *lemmatizer) next() (*Token, error) {
	if t == nil {
		return nil, nil
	}
	for {

		if len(t.outgoing) > 0 {
			return t.emit(), nil
		}

		_, err := t.fill(1)
		if err != nil {
			return nil, err
		}

		if len(t.buffer) == 0 {
			return nil, nil
		}

		switch token := t.buffer[0]; {
		case token.IsPunct() || token.IsSpace():
			// Emit it straight from the incoming buffer
			t.drop(1)
			return token, nil
		default:
			// Else it's a word
			t.ngrams()
		}
	}
}

func (t *lemmatizer) ngrams() error {
	// Try n-grams, longest to shortest (greedy)
	for take := t.dictionary.MaxGramLength(); take > 0; take-- {
		wordrun, err := t.wordrun(take)
		if err != nil {
			if err == errInsufficient {
				// No problem, try the next n-gram
				continue
			}
			return err
		}

		canonical, found := t.dictionary.Lookup(wordrun.words)

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
					tok, err := tokens.Next()
					if err != nil {
						return err
					}
					if tok == nil {
						break
					}
					tok.lemma = true
					t.stage(tok) //set it up to be emitted
				}
			} else {
				token := &Token{
					value: canonical,
					space: false,
					punct: false,
					lemma: true,
				}
				t.stage(token) //set it up to be emitted
			}

			t.drop(wordrun.consumed) // discard the incoming tokens that comprised the lemma
			return nil
		}

		if take == 1 {
			// No n-grams, just emit
			original := t.buffer[0]
			t.stage(original) // set it up to be emitted
			t.drop(1)         // take it out of the buffer
			return nil
		}
	}
	return fmt.Errorf("did not find a token in ngrams. this should never happen")
}

// drop (truncate) the first `n` elements of the buffer
// remember, a token being in the buffer does not imply that we will emit it
func (t *lemmatizer) drop(n int) {
	copy(t.buffer, t.buffer[n:])
	t.buffer = t.buffer[:len(t.buffer)-n]
}

// ensure that the buffer contains at least `count` elements; returns false if channel is exhausted before achieving the count
func (t *lemmatizer) fill(count int) (bool, error) {
	for count >= len(t.buffer) {
		token, err := t.incoming.Next()
		if err != nil {
			return false, err
		}
		if token == nil {
			// EOF
			return false, nil
		}
		t.buffer = append(t.buffer, token)
	}
	return true, nil
}

func (t *lemmatizer) stage(tok *Token) {
	t.outgoing = append(t.outgoing, tok)
}

func (t *lemmatizer) emit() *Token {
	n := 1
	tok := t.outgoing[0]
	copy(t.outgoing, t.outgoing[n:])
	t.outgoing = t.outgoing[:len(t.outgoing)-n]
	return tok
}

type wordrun struct {
	words    []string
	consumed int
}

var empty = wordrun{}
var errInsufficient = errors.New("could not find word run of desired length")

func (t *lemmatizer) wordrun(take int) (wordrun, error) {
	var (
		taken []string // the words
		count int      // tokens consumed, not necessarily equal to take
	)

	for len(taken) < take {
		ok, err := t.fill(count)
		if err != nil {
			return empty, err
		}
		if !ok {
			// Not enough (buffered) tokens to continue
			// So, a word run of length `take` is impossible
			return empty, errInsufficient
		}

		token := t.buffer[count]
		switch {
		case token.IsPunct():
			// Note: test for punct before space; newlines and tabs can be
			// considered both punct and space (depending on the tokenizer!)
			// and we want to treat them as breaking word runs.
			return empty, errInsufficient
		case token.IsSpace():
			// Ignore and continue
			count++
		default:
			// Found a word
			taken = append(taken, token.String())
			count++
		}
	}

	result := wordrun{
		words:    taken,
		consumed: count,
	}

	return result, nil
}
