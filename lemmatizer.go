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
func (lem *lemmatizer) next() (*Token, error) {
	for {
		if len(lem.outgoing) > 0 {
			return lem.emit(), nil
		}

		err := lem.fill(1)
		if err != nil {
			if err == errInsufficient {
				// EOF, no problem
				return nil, nil
			}
			return nil, err
		}

		switch token := lem.buffer[0]; {
		case token.IsPunct() || token.IsSpace():
			// Emit it straight from the incoming buffer
			lem.drop(1)
			return token, nil
		default:
			// Else it's a word
			lem.ngrams()
		}
	}
}

func (lem *lemmatizer) ngrams() error {
	// Try n-grams, longest to shortest (greedy)
	for desired := lem.dictionary.MaxGramLength(); desired > 0; desired-- {
		wordrun, err := lem.wordrun(desired)
		if err != nil {
			if err == errInsufficient {
				// No problem, try the next (smaller) n-gram
				continue
			}
			return err
		}

		canonical, found := lem.dictionary.Lookup(wordrun.words)

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
					token, err := tokens.Next()
					if err != nil {
						return err
					}
					if token == nil {
						break
					}
					token.lemma = true
					lem.stage(token) //set it up to be emitted
				}
			} else {
				token := &Token{
					value: canonical,
					space: false,
					punct: false,
					lemma: true,
				}
				lem.stage(token) //set it up to be emitted
			}

			lem.drop(wordrun.consumed) // discard the incoming tokens that comprised the lemma
			return nil
		}

		if desired == 1 {
			// No n-grams, just emit
			original := lem.buffer[0]
			lem.stage(original) // set it up to be emitted
			lem.drop(1)         // take it out of the buffer
			return nil
		}
	}
	return fmt.Errorf("did not find a token in ngrams. this should never happen")
}

// drop (truncate) the first `n` elements of the buffer
// remember, a token being in the buffer does not imply that we will emit it
func (lem *lemmatizer) drop(n int) {
	copy(lem.buffer, lem.buffer[n:])
	lem.buffer = lem.buffer[:len(lem.buffer)-n]
}

// ensure that the buffer contains at least `desired` elements; returns false if channel is exhausted before achieving the count
func (lem *lemmatizer) fill(desired int) error {
	for len(lem.buffer) < desired {
		token, err := lem.incoming.Next()
		if err != nil {
			return err
		}
		if token == nil {
			// EOF
			return errInsufficient
		}
		lem.buffer = append(lem.buffer, token)
	}
	return nil
}

func (lem *lemmatizer) stage(token *Token) {
	lem.outgoing = append(lem.outgoing, token)
}

func (lem *lemmatizer) emit() *Token {
	n := 1
	token := lem.outgoing[0]
	copy(lem.outgoing, lem.outgoing[n:])
	lem.outgoing = lem.outgoing[:len(lem.outgoing)-n]
	return token
}

type wordrun struct {
	words    []string
	consumed int
}

var (
	empty           = wordrun{}
	errInsufficient = errors.New("could not find word run of desired length")
)

func (lem *lemmatizer) wordrun(desired int) (wordrun, error) {
	var (
		words    []string
		consumed int // tokens consumed or 'seen', not necessarily equal to desired
	)

	for len(words) < desired {
		err := lem.fill(consumed + 1)
		if err != nil {
			// If errInsufficient, not enough (buffered) tokens to continue,
			// so a word run of desired length is impossible; be handled by ngrams().
			// Other errors are just errors; pass 'em back.
			return empty, err
		}

		token := lem.buffer[consumed]
		switch {
		case token.IsPunct():
			// Note: test for punct before space; newlines and tabs can be
			// considered both punct and space (depending on the tokenizer!)
			// and we want to treat them as breaking word runs.
			return empty, errInsufficient
		case token.IsSpace():
			// Ignore and continue
			consumed++
		default:
			// Found a word
			words = append(words, token.String())
			consumed++
		}
	}

	result := wordrun{
		words:    words,
		consumed: consumed,
	}

	return result, nil
}
