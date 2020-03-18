package stack

import (
	"unicode"

	"github.com/clipperhouse/jargon"
)

type TokenTrie struct {
	root       *Node
	ignore     map[rune]bool
	ignoreCase bool
}

type Node struct {
	children     map[rune]*Node
	hasCanonical bool
	canonical    string
}

func newTokenTrie(ignoreCase bool, ignore []rune) *TokenTrie {
	ignoreset := map[rune]bool{}
	for _, r := range ignore {
		ignoreset[r] = true
	}
	return &TokenTrie{
		root:       &Node{},
		ignoreCase: ignoreCase,
		ignore:     ignoreset,
	}
}

func (trie *TokenTrie) Add(tokens []*jargon.Token, canonical string) {
	node := trie.root
	for _, token := range tokens {
		for _, r := range token.String() {
			if trie.ignoreCase {
				r = unicode.ToLower(r)
			}

			if trie.ignore[r] {
				continue
			}

			child := node.children[r]
			if child == nil {
				if node.children == nil {
					node.children = map[rune]*Node{}
				}
				child = &Node{}
				node.children[r] = child
			}
			node = child
		}
	}

	node.hasCanonical = true
	node.canonical = canonical
}

func (trie *TokenTrie) SearchCanonical(tokens ...*jargon.Token) (found bool, canonical string, consumed int) {
	node := trie.root

outer:
	for i, token := range tokens {
		for _, r := range token.String() {
			if trie.ignoreCase {
				r = unicode.ToLower(r)
			}

			if trie.ignore[r] {
				continue
			}

			node = node.children[r]
			if node == nil {
				break outer
			}

			if node.hasCanonical {
				found = true
				canonical = node.canonical
				consumed = i + 1
			}
		}
	}

	return found, canonical, consumed
}
