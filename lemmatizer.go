// Package jargon offers tokenizers and lemmatizers, for use in text processing and NLP
package jargon

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/clipperhouse/jargon/stackexchange"
)

// Lemmatize transforms Tokens to their canonicalized ("lemmatized") terms.
//
// Tokens that are not canonicalized are returned as-is, e.g. for tokenized input:
//     "I", " ", "think", " ", "Ruby", " ", "on", " ", "Rails", " ", "is", " ", "great"
//
// lemmatized output:
//     "I", " ", "think", " ", "ruby-on-rails", " ", "is", " ", "great"
//
// Use token.IsLemma() to find out if a term was lemmatized
//
// Note that fewer tokens may be returned than were input. In this case, the five tokens
// representing Ruby<space>on<space>Rails are combined into a single token.
func (incoming *Tokens) Lemmatize(filters ...TokenFilter) *Tokens {
	if len(filters) == 0 {
		filters = append(filters, stackexchange.Tags)
	}

	var result = incoming

	// new lemmatizer for each filter, results as input into next lemmatizer
	for _, filter := range filters {
		lem := newLemmatizer(result, filter)
		result = &Tokens{
			Next: lem.next,
		}
	}

	return result
}

// LemmatizeString transforms words to their canonicalized ("lemmatized") terms
func LemmatizeString(s string, filters ...TokenFilter) string {
	r := strings.NewReader(s)
	tokens := Tokenize(r)
	lemmatized := tokens.Lemmatize(filters...)

	// We can elide the error because it's coming from a string, no real I/O
	result, _ := lemmatized.String()

	return result
}

func newLemmatizer(incoming *Tokens, filter TokenFilter) *lemmatizer {
	return &lemmatizer{
		incoming: incoming,
		filter:   filter,
		buffer:   &queue{},
		outgoing: &queue{},
	}
}

type lemmatizer struct {
	incoming *Tokens
	filter   TokenFilter
	buffer   *queue // for incoming tokens; no guarantee they will be emitted
	outgoing *queue
}

// next returns the next token; nil indicates end of data
func (lem *lemmatizer) next() (*Token, error) {
	for {
		if lem.outgoing.len() > 0 {
			return lem.outgoing.pop(), nil
		}

		err := lem.fill(1)
		if err != nil {
			if err == errInsufficient {
				// EOF, no problem
				return nil, nil
			}
			return nil, err
		}

		peek := lem.buffer.peek()
		if peek.IsPunct() || peek.IsSpace() {
			token := lem.buffer.pop()
			lem.outgoing.push(token)
			continue
		}

		// Else it's a word
		lem.ngrams()
	}
}

func (lem *lemmatizer) ngrams() error {
	// Try n-grams, longest to shortest (greedy)
	for desired := lem.filter.MaxGramLength(); desired > 0; desired-- {
		wordrun, err := lem.wordrun(desired)
		if err != nil {
			if err == errInsufficient {
				// No problem, try the next (smaller) n-gram
				continue
			}
			return err
		}

		canonical, found := lem.filter.Lookup(wordrun.words)

		if found {
			// if returned value is empty, interpret as "remove token", e.g. the stopwords filter
			if canonical == "" {
				lem.buffer.drop(wordrun.consumed)
				continue
			}

			// the canonical might have space or punct, so we want to re-tokenize

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
					//set it up to be emitted
					lem.outgoing.push(token)
				}
			} else {
				token := &Token{
					value: canonical,
					space: false,
					punct: false,
					lemma: true,
				}
				//set it up to be emitted
				lem.outgoing.push(token)
			}

			// discard the incoming tokens that comprised the lemma
			lem.buffer.drop(wordrun.consumed)
			return nil
		}

		if desired == 1 {
			// No n-grams, just emit the next token
			token := lem.buffer.pop()
			lem.outgoing.push(token)
			return nil
		}
	}
	return fmt.Errorf("did not find a token in ngrams. this should never happen")
}

// ensure that the buffer contains at least `desired` elements; returns false if channel is exhausted before achieving the count
func (lem *lemmatizer) fill(desired int) error {
	for lem.buffer.len() < desired {
		token, err := lem.incoming.Next()
		if err != nil {
			return err
		}
		if token == nil {
			// EOF
			return errInsufficient
		}
		lem.buffer.push(token)
	}
	return nil
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

		token := lem.buffer.tokens[consumed]
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
