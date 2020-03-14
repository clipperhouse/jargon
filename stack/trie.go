package stack

import (
	"fmt"
	"strings"

	"github.com/clipperhouse/jargon"
)

type TokenTrie struct {
	Children map[string]*TokenTrie
	Value    string
}

var ignore = map[string]bool{
	"-": true,
	" ": true,
	".": true,
}

func normalize(s string) string {
	for k := range ignore {
		s = strings.ReplaceAll(s, k, "")
	}
	return s
}

func (trie *TokenTrie) Add(tokens []*jargon.Token, value string) bool {
	node := trie
	for _, token := range tokens {

		key := token.String()
		if ignore[key] {
			continue
		}

		key = normalize(token.String())
		if key == "" {
			continue
		}

		child := node.Children[key]
		if child == nil {
			if node.Children == nil {
				node.Children = map[string]*TokenTrie{}
			}
			child = &TokenTrie{}
			node.Children[key] = child
		}
		node = child
	}
	// does node have an existing value?
	isNewVal := node.Value == ""
	node.Value = value
	return isNewVal
}

func (trie *TokenTrie) Search(tokens ...*jargon.Token) *TokenTrie {
	node := trie
	for _, token := range tokens {
		node = node.Children[token.String()]
		if node == nil {
			return nil
		}
	}
	return node
}

func (trie *TokenTrie) SearchCanonical(tokens ...*jargon.Token) (*TokenTrie, int) {
	var result *TokenTrie
	consumed := 0

	node := trie
	for i, token := range tokens {
		key := token.String()
		if ignore[key] {
			continue
		}

		key = normalize(key)
		if key == "" {
			continue
		}

		node = node.Children[key]
		fmt.Println(key)
		fmt.Println(node)
		if node == nil {
			break
		}

		if node.Value != "" {
			consumed = i + 1
			result = node
		}
	}

	return result, consumed
}
