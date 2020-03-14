package stack

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestTrie(t *testing.T) {
	s := "ruby on rails and asp.net and node.js"
	tokens, err := jargon.Tokenize(strings.NewReader(s)).ToSlice()

	if err != nil {
		t.Error(err)
	}

	found, consumed := trie.SearchCanonical(tokens[:9]...)
	t.Log(found)
	t.Log(consumed)
}
