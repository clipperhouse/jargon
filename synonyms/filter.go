package synonyms

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/synonyms/trie"
)

// Filter is the data structure of a synonyms filter. Use NewFilter to create.
type Filter struct {
	Trie     *trie.RuneTrie
	MaxWords int
}

// NewFilter creates a new synonyms Filter
func NewFilter(mappings map[string]string, ignoreCase bool, ignoreRunes []rune) (*Filter, error) {
	trie := trie.New(ignoreCase, ignoreRunes)
	maxWords := 1
	for k, v := range mappings {
		synonyms := strings.Split(k, ",")
		for _, synonym := range synonyms {

			synonym = strings.TrimSpace(synonym)
			tokens, err := jargon.TokenizeString(synonym).ToSlice()
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
		Trie:     trie,
		MaxWords: maxWords,
	}, nil
}

// Filter replaces tokens with their canonical terms, based on Stack Overflow tags & synonyms
func (f *Filter) Filter(incoming *jargon.Tokens) *jargon.Tokens {
	t := newTokens(incoming, f.Trie, f.MaxWords)
	return &jargon.Tokens{
		Next: t.next,
	}
}

func newTokens(incoming *jargon.Tokens, trie *trie.RuneTrie, maxWords int) *tokens {
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
	trie     *trie.RuneTrie
	maxWords int
}

// next returns the next token; nil indicates end of data
func (t *tokens) next() (*jargon.Token, error) {
	// Clear out any outgoing
	if t.outgoing.Len() > 0 {
		return t.outgoing.Pop(), nil
	}

	// Consume all the words
	for {
		err := t.fill()
		if err != nil {
			return nil, err
		}

		run := t.wordrun()
		if len(run) == 0 {
			// No more words
			break
		}

		// Try to lemmatize
		found, canonical, consumed := t.trie.SearchCanonical(run...)
		if found {
			if canonical != "" {
				token := jargon.NewToken(canonical, true)
				t.outgoing.Push(token)
			}
			t.buffer.Drop(consumed)
			continue
		}

		t.buffer.PopTo(t.outgoing)
	}

	// Queue up the rest of the buffer to go out
	t.buffer.FlushTo(t.outgoing)

	if t.outgoing.Len() > 0 {
		return t.outgoing.Pop(), nil
	}

	return nil, nil
}

// fill the buffer until EOF, punctuation, or enough word tokens
func (t *tokens) fill() error {
	words := 0
	for _, token := range t.buffer.All() {
		if !token.IsPunct() && !token.IsSpace() {
			words++
		}
	}

	for words < t.maxWords {
		token, err := t.incoming.Next()
		if err != nil {
			return err
		}
		if token == nil {
			// EOF
			break
		}

		t.buffer.Push(token)

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
func (t *tokens) wordrun() []*jargon.Token {
	spaces := 0
	for _, token := range t.buffer.All() {
		if !token.IsSpace() {
			break
		}
		t.outgoing.Push(token)
		spaces++
	}
	t.buffer.Drop(spaces)

	var (
		run      []*jargon.Token
		words    int
		consumed int
	)

	for _, token := range t.buffer.All() {
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

		if words >= t.maxWords {
			break
		}
	}

	return run
}

// String returns a Go source declaration of the Filter
func (f *Filter) String() string {
	return f.Decl()
}

// Decl returns a Go source declaration of the Filter
func (f *Filter) Decl() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "&synonyms.Filter{\n")
	if f.Trie != nil {
		fmt.Fprintf(&b, "Trie: %s,\n", f.Trie.Decl())
	}
	if f.MaxWords > 0 {
		// default value does not need to be declared
		fmt.Fprintf(&b, "MaxWords: %d,\n", f.MaxWords)
	}
	fmt.Fprintf(&b, "}")

	result := b.Bytes()
	// result, err := format.Source(result)
	// if err != nil {
	// 	panic(err)
	// }

	return string(result)
}
