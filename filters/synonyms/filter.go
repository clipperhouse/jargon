// Package synonyms provides a builder for filtering and replacing synonyms in a token stream
package synonyms

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/synonyms/trie"
	"github.com/clipperhouse/jargon/tokenqueue"
)

// filter is the data structure of a synonyms filter. Use NewFilter to create.
type filter struct {
	// For lazy loading, see build() below
	config *config
	once   sync.Once

	trie     *trie.RuneTrie
	maxWords int
}

type config struct {
	mappings    map[string]string
	ignoreCase  bool
	ignoreRunes []rune
}

// NewFilter creates a new synonyms Filter
func NewFilter(mappings map[string]string, ignoreCase bool, ignoreRunes []rune) jargon.Filter {
	// Save the parameters for lazy loading (below)
	f := &filter{
		config: &config{
			mappings:    mappings,
			ignoreCase:  ignoreCase,
			ignoreRunes: ignoreRunes,
		},
	}
	return f.Filter
}

func (f *filter) build() error {
	trie := trie.New(f.config.ignoreCase, f.config.ignoreRunes)
	maxWords := 1
	for synonyms, canonical := range f.config.mappings {
		tokens, err := jargon.TokenizeString(synonyms).ToSlice()
		if err != nil {
			return err
		}

		start := 0
		skipSpaces := true
		for i, token := range tokens {
			// Leading spaces, and spaces following commas, should be ignored
			if skipSpaces && token.IsSpace() {
				start++
				continue
			}
			skipSpaces = false

			if token.String() == "," {
				slice := tokens[start:i]
				trie.Add(slice, canonical)
				updateMaxWords(slice, &maxWords)

				start = i + 1 // ignore the comma
				skipSpaces = true
				continue
			}
		}

		// Remaining after the last comma
		slice := tokens[start:]
		trie.Add(slice, canonical)
		updateMaxWords(slice, &maxWords)
	}

	// Populate with new values
	f.trie = trie
	f.maxWords = maxWords

	// Kill the config
	f.config = nil

	return nil
}

func updateMaxWords(tokens []*jargon.Token, maxWords *int) {
	words := 0
	for _, token := range tokens {
		if !token.IsSpace() && !token.IsPunct() {
			words++
		}
	}
	if words > *maxWords {
		*maxWords = words
	}
}

// Filter replaces tokens with their canonical terms, based on Stack Overflow tags & synonyms
func (f *filter) Filter(incoming *jargon.TokenStream) *jargon.TokenStream {
	// Lazily build the trie on first call, i.e. don't pay for the construction
	// unless we use it
	var err error
	f.once.Do(func() {
		err = f.build()
	})

	t := &tokens{
		incoming: incoming,
		buffer:   tokenqueue.New(),
		outgoing: tokenqueue.New(),
		filter:   f,
	}

	// Catch the error that may have resulted from lazy construction above
	next := func() (*jargon.Token, error) {
		if err != nil {
			return nil, err
		}
		return t.next()
	}

	return jargon.NewTokenStream(next)
}

type tokens struct {
	// incoming stream of tokens from another source, such as a tokenizer
	incoming *jargon.TokenStream
	// a 'lookahead' buffer for incoming tokens
	buffer *tokenqueue.TokenQueue
	// outgoing queue of filtered tokens
	outgoing *tokenqueue.TokenQueue
	filter   *filter
}

// next returns the next token; nil indicates end of data
func (t *tokens) next() (*jargon.Token, error) {
	// Buffer should be clear after every call to next()
	if t.buffer.Len() != 0 {
		return nil, fmt.Errorf("expected buffer to be empty")
	}

	// Clear out any outgoing
	if t.outgoing.Any() {
		return t.outgoing.Pop(), nil
	}

	// Consume all the words
	for {
		err := t.fill()
		if err != nil {
			return nil, err
		}

		if t.buffer.Len() == 0 {
			// Nothing left
			break
		}

		run, err := t.wordrun()
		if err != nil {
			return nil, err
		}

		if len(run) == 0 {
			// No more words
			break
		}

		// Try to lemmatize
		found, canonical, consumed := t.filter.trie.SearchCanonical(run...)
		if found {
			if canonical != "" {
				token := jargon.NewToken(canonical, true)
				t.outgoing.Push(token)
			}
			t.buffer.Drop(consumed)
			continue
		}

		// The word didn't lemmatize, pass it along verbatim
		t.buffer.PopTo(t.outgoing)
	}

	// Queue up the rest of the buffer to go out
	t.buffer.FlushTo(t.outgoing)

	if t.outgoing.Any() {
		return t.outgoing.Pop(), nil
	}

	return nil, nil
}

// fill the buffer until EOF, punctuation, or enough word tokens
func (t *tokens) fill() error {
	// Leading buffered space & punct should go straight out
	drop := 0
	for _, token := range t.buffer.Tokens {
		if token.IsSpace() || token.IsPunct() {
			t.outgoing.Push(token)
			drop++
			continue
		}
		// Else, it's a word
		break
	}
	t.buffer.Drop(drop)

	if t.buffer.Len() == 0 {
		// Leading incoming space & punct should go straight out, don't even buffer
		for t.incoming.Scan() {
			token := t.incoming.Token()
			if token.IsSpace() || token.IsPunct() {
				t.outgoing.Push(token)
				continue
			}
			// It's a word
			t.buffer.Push(token)
			break
		}
		// Gotta check err after Scan
		if err := t.incoming.Err(); err != nil {
			return err
		}
	}

	// Count the words we have
	wordcount := 0
	for _, token := range t.buffer.Tokens {
		if !token.IsPunct() && !token.IsSpace() {
			wordcount++
		}
	}

	// Fill until we have enough words, hit a punct, or EOF
	for wordcount < t.filter.maxWords {
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
		wordcount++
	}

	return nil
}

// wordrun pulls the longest series of tokens comprised of words
func (t *tokens) wordrun() ([]*jargon.Token, error) {
	// Requires a previous call to fill()
	if t.buffer.Len() == 0 {
		return nil, fmt.Errorf("expected buffer to have tokens")
	}
	head := t.buffer.Tokens[0]
	if head.IsSpace() || head.IsPunct() {
		return nil, fmt.Errorf("expected buffer to have word as first token, got %q", head)
	}

	var (
		end      int
		words    int
		consumed int
	)

	for i, token := range t.buffer.Tokens {
		if token.IsPunct() {
			// fall through and send back word run we've gotten so far (if any)
			// don't consume this punct, leave it in the buffer
			break
		}

		// It's a word or space
		end = i + 1
		consumed++

		if !token.IsSpace() {
			// It's a word
			words++
		}

		if words >= t.filter.maxWords {
			break
		}
	}

	return t.buffer.Tokens[:end], nil
}

// String returns a Go source declaration of the Filter
func (f *filter) String() string {
	return f.Decl()
}

// Decl returns a Go source declaration of the Filter
func (f *filter) Decl() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "&synonyms.Filter{\n")
	if f.trie != nil {
		fmt.Fprintf(&b, "Trie: %s,\n", f.trie.String())
	}
	if f.maxWords > 0 {
		// default value does not need to be declared
		fmt.Fprintf(&b, "MaxWords: %d,\n", f.maxWords)
	}
	fmt.Fprintf(&b, "}")

	result := b.Bytes()
	// result, err := format.Source(result)
	// if err != nil {
	// 	panic(err)
	// }

	return string(result)
}
