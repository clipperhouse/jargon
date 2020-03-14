package stack

import (
	"strings"

	"github.com/clipperhouse/jargon"
)

func Filter(incoming *jargon.Tokens) *jargon.Tokens {
	f := newFilter(incoming)
	return &jargon.Tokens{
		Next: f.next,
	}
}

func newFilter(incoming *jargon.Tokens) *filter {
	return &filter{
		incoming:      incoming,
		buffer:        &jargon.TokenQueue{},
		outgoing:      &jargon.TokenQueue{},
		maxGramLength: 5,
	}
}

type filter struct {
	incoming      *jargon.Tokens
	buffer        *jargon.TokenQueue // for incoming tokens; no guarantee they will be emitted
	outgoing      *jargon.TokenQueue
	maxGramLength int
}

// next returns the next token; nil indicates end of data
func (f *filter) next() (*jargon.Token, error) {
	if f.outgoing.Len() > 0 {
		return f.outgoing.Pop(), nil
	}

	for {
		// Fill the buffer
		err := f.fill()
		if err != nil {
			return nil, err
		}

		// Nothing was filled, must be EOF
		if f.buffer.Len() == 0 {
			return nil, nil
		}

		// Consume all the words

		for {
			run := f.wordrun()

			if len(run) == 0 {
				// No more words
				break
			}

			// Try to lemmatize
			//node, consumed := trie.SearchCanonical(run...)
		}

		// Queue up the rest of the buffer to go out
		for range f.buffer.All() {
			token := f.buffer.Pop()
			f.outgoing.Push(token)
		}
	}
}

func (f *filter) wordrun() []*jargon.Token {
	var (
		run       []*jargon.Token
		wordcount int
		consumed  int // tokens consumed or 'seen', not necessarily equal to desired
	)

	for _, token := range f.buffer.All() {
		if token.IsPunct() {
			// send back word run we've gotten so far (if any)
			// don't consume this punct, leave it in the buffer
			break
		}

		// It's a word or space
		run = append(run, token)
		consumed++

		if !token.IsSpace() {
			wordcount++
		}

		if wordcount >= f.maxGramLength {
			break
		}
	}

	return run
}

// fill the buffer until EOF, punctuation, or enough word tokens to lemmatize
func (f *filter) fill() error {
	wordcount := 0

	for wordcount < f.maxGramLength {
		token, err := f.incoming.Next()
		if err != nil {
			return err
		}
		if token == nil {
			// EOF
			return nil
		}

		f.buffer.Push(token)

		if !token.IsPunct() && !token.IsSpace() {
			// it's a word
			wordcount++
		}

		if token.IsPunct() {
			break
		}
	}

	return nil
}

var trie = &TokenTrie{}

func init() {
	for k, v := range mappings {
		synonyms := strings.Split(k, ",")
		for _, synonym := range synonyms {
			tokens, err := jargon.Tokenize(strings.NewReader(synonym)).ToSlice()
			if err != nil {
				panic(err)
			}
			trie.Add(tokens, v)
		}
	}
}
