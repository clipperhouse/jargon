package stack

import (
	"strings"

	"github.com/clipperhouse/jargon"
)

type Synonyms struct {
	trie     *TokenTrie
	maxWords int
}

func NewSynonyms(mappings map[string]string, ignoreCase bool, ignoreRunes []rune) (*Synonyms, error) {
	trie := newTokenTrie(ignoreCase, ignoreRunes)
	maxWords := 1
	for k, v := range mappings {
		synonyms := strings.Split(k, ",")
		for _, synonym := range synonyms {

			synonym = strings.TrimSpace(synonym)
			r := strings.NewReader(synonym)
			tokens, err := jargon.Tokenize(r).ToSlice()
			if err != nil {
				return nil, err
			}

			words := 0
			for _, token := range tokens {
				if !token.IsSpace() && !token.IsPunct() {
					words++
				}
			}
			if words > maxWords {
				maxWords = words
			}

			trie.Add(tokens, v)
		}
	}

	return &Synonyms{
		trie:     trie,
		maxWords: maxWords,
	}, nil
}

func (syns *Synonyms) Filter(incoming *jargon.Tokens) *jargon.Tokens {
	f := newFilter(incoming, syns.trie, syns.maxWords)
	return &jargon.Tokens{
		Next: f.next,
	}
}

func newFilter(incoming *jargon.Tokens, trie *TokenTrie, maxWords int) *filter {
	return &filter{
		incoming: incoming,
		buffer:   &jargon.TokenQueue{},
		outgoing: &jargon.TokenQueue{},
		trie:     trie,
		maxWords: maxWords,
	}
}

type filter struct {
	// incoming stream of tokens from another source, such as a tokenizer
	incoming *jargon.Tokens
	// a 'lookahead' buffer for incoming tokens
	buffer *jargon.TokenQueue
	// outgoing queue of filtered tokens
	outgoing *jargon.TokenQueue
	trie     *TokenTrie
	maxWords int
}

// next returns the next token; nil indicates end of data
func (f *filter) next() (*jargon.Token, error) {
	// Clear out any outgoing
	if f.outgoing.Len() > 0 {
		return f.outgoing.Pop(), nil
	}

	// Consume all the words
	for {
		err := f.fill()
		if err != nil {
			return nil, err
		}

		run := f.wordrun()
		if len(run) == 0 {
			// No more words
			break
		}

		// Try to lemmatize
		found, canonical, consumed := f.trie.SearchCanonical(run...)
		if found {
			if canonical != "" {
				token := jargon.NewToken(canonical, true)
				f.outgoing.Push(token)
			}
			f.buffer.Drop(consumed)
			continue
		}

		f.buffer.PopTo(f.outgoing)
	}

	// Queue up the rest of the buffer to go out
	f.buffer.FlushTo(f.outgoing)

	if f.outgoing.Len() > 0 {
		return f.outgoing.Pop(), nil
	}

	return nil, nil
}

func (f *filter) fill() error {
	// Fill the buffer until EOF, punctuation, or enough word tokens
	words := 0
	for _, token := range f.buffer.All() {
		if !token.IsPunct() && !token.IsSpace() {
			words++
		}
	}

	for words < f.maxWords {
		token, err := f.incoming.Next()
		if err != nil {
			return err
		}
		if token == nil {
			// EOF
			break
		}

		f.buffer.Push(token)

		if token.IsPunct() {
			break
		}

		if token.IsSpace() {
			continue
		}

		// It's a word
		words++
	}

	return nil
}

func (f *filter) wordrun() []*jargon.Token {
	spaces := 0
	for _, token := range f.buffer.All() {
		if !token.IsSpace() {
			break
		}
		f.outgoing.Push(token)
		spaces++
	}
	f.buffer.Drop(spaces)

	var (
		run      []*jargon.Token
		words    int
		consumed int
	)

	for _, token := range f.buffer.All() {
		if token.IsPunct() {
			// fall through and send back word run we've gotten so far (if any)
			// don't consume this punct, leave it in the buffer
			break
		}

		// It's a word or space
		run = append(run, token)
		consumed++

		if !token.IsSpace() {
			// It's a word
			words++
		}

		if words >= f.maxWords {
			break
		}
	}

	return run
}
