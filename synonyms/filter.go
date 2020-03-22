package synonyms

import (
	"strings"

	"github.com/clipperhouse/jargon"
)

type Filter struct {
	trie     *runeTrie
	maxWords int
}

func NewFilter(mappings map[string]string, ignoreCase bool, ignoreRunes []rune) (*Filter, error) {
	trie := newTrie(ignoreCase, ignoreRunes)
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

	return &Filter{
		trie:     trie,
		maxWords: maxWords,
	}, nil
}

func (syns *Filter) Filter(incoming *jargon.Tokens) *jargon.Tokens {
	f := newTokens(incoming, syns.trie, syns.maxWords)
	return &jargon.Tokens{
		Next: f.next,
	}
}

func newTokens(incoming *jargon.Tokens, trie *runeTrie, maxWords int) *tokens {
	return &tokens{
		incoming: incoming,
		buffer:   &jargon.TokenQueue{},
		outgoing: &jargon.TokenQueue{},
		trie:     trie,
		maxWords: maxWords,
	}
}

type tokens struct {
	// incoming stream of tokens from another source, such as a tokenizer
	incoming *jargon.Tokens
	// a 'lookahead' buffer for incoming tokens
	buffer *jargon.TokenQueue
	// outgoing queue of filtered tokens
	outgoing *jargon.TokenQueue
	trie     *runeTrie
	maxWords int
}

// next returns the next token; nil indicates end of data
func (f *tokens) next() (*jargon.Token, error) {
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

// fill the buffer until EOF, punctuation, or enough word tokens
func (f *tokens) fill() error {
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

// wordrun pulls the longest series of tokens comprised of words
func (f *tokens) wordrun() []*jargon.Token {
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
