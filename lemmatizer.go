package jargon

import (
	"strings"
)

// Lemmatizer is the main structure for looking up canonical tags
type Lemmatizer struct {
	values map[string]string
}

// NewLemmatizer creates and populates a new Lemmatizer for the purpose of looking up canonical tags
func NewLemmatizer(d Dictionary) *Lemmatizer {
	result := &Lemmatizer{
		values: make(map[string]string),
	}
	tags := d.GetTags()
	for _, tag := range tags {
		key := normalize(tag)
		result.values[key] = tag
	}
	synonyms := d.GetSynonyms()
	for synonym, canonical := range synonyms {
		key := normalize(synonym)
		result.values[key] = canonical
	}
	return result
}

var gramLengths = []int{3, 2, 1}

// LemmatizeTokens takes a slice of tokens and returns tokens with canonicalized terms.
// Terms (tokens) that are not canonicalized are returned as-is, e.g.
// ["I", " ", "think", " ", "Ruby", " ", "on", " ", "Rails", " ", "is", " ", "great"] →
//    ["I", " ", "think", " ", "ruby-on-rails", " ", "is", " ", "great"]
// Note that fewer tokens may be returned than were input.
func (lem *Lemmatizer) LemmatizeTokens(tokens []Token) []Token {
	lemmatized := make([]Token, 0)
	pos := 0

	for pos < len(tokens) {
		switch current := tokens[pos]; {
		case current.Punct() || current.Space():
			// Emit it
			lemmatized = append(lemmatized, current)
			pos++
		default:
		Grams:
			// Else it's a word, try n-grams
			for _, take := range gramLengths {
				run, consumed, ok := wordrun(tokens, pos, take)
				if ok {
					gram := Join(run, Token.Value)
					key := normalize(gram)
					canonical, found := lem.values[key]

					if found {
						// Emit token, replacing consumed tokens
						token := NewToken(canonical, false, false)
						lemmatized = append(lemmatized, token)
						pos += consumed
						break Grams
					}

					if take == 1 {
						// No n-grams, just emit
						token := tokens[pos]
						lemmatized = append(lemmatized, token)
						pos++
					}
				}
			}
		}
	}

	return lemmatized
}

// normalize returns a string suitable as a key for tag lookup, removing dots, dashes and forward slashes, and converting to lowercase
func normalize(s string) string {
	result := make([]rune, 0)

	for index, value := range s {
		if index == 0 {
			// Leading dots are meaningful and should not be removed, for example ".net"
			result = append(result, value)
			continue
		}
		if value == '.' || value == '-' || value == '/' {
			continue
		}
		result = append(result, value)
	}
	return strings.ToLower(string(result))
}

// Analogous to tokens.Skip(skip).Take(take) in Linq
func wordrun(tokens []Token, skip, take int) ([]Token, int, bool) {
	taken := make([]Token, 0)
	consumed := 0 // tokens consumed, not necessarily equal to take

	for len(taken) < take {
		end := skip + consumed
		eof := end >= len(tokens)

		if eof {
			// Hard stop
			return nil, 0, false
		}

		candidate := tokens[end]
		switch {
		// Note: test for punct before space; newlines and tabs can be
		// considered both punct and space (depending on the tokenizer!)
		// and we want to treat them as breaking word runs.
		case candidate.Punct():
			// Hard stop
			return nil, 0, false
		case candidate.Space():
			// Ignore and continue
			consumed++
		default:
			// Found a word
			taken = append(taken, candidate)
			consumed++
		}
	}

	return taken, consumed, true
}
