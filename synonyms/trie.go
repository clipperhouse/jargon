package synonyms

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/clipperhouse/jargon"
)

type runeTrie struct {
	root       *node
	ignore     map[rune]bool
	ignoreCase bool
}

type node struct {
	children     map[rune]*node
	hasCanonical bool
	canonical    string
}

func (n *node) String() string {
	var w strings.Builder
	w.WriteString("{")
	if n.children != nil {
		for k, v := range n.children {
			w.WriteRune(k)
			w.WriteString(": ")
			w.WriteString(v.String())
		}
	}
	w.WriteString(fmt.Sprint(n.hasCanonical))
	w.WriteString(" ")
	w.WriteString(n.canonical)
	w.WriteString("} ")

	return w.String()
}

func newTrie(ignoreCase bool, ignore []rune) *runeTrie {
	set := map[rune]bool{}
	for _, r := range ignore {
		set[r] = true
	}
	return &runeTrie{
		root:       &node{},
		ignoreCase: ignoreCase,
		ignore:     set,
	}
}

func (trie *runeTrie) Add(tokens []*jargon.Token, canonical string) {
	n := trie.root
	for _, token := range tokens {
		for _, r := range token.String() {
			if trie.ignoreCase {
				r = unicode.ToLower(r)
			}

			if trie.ignore[r] {
				continue
			}

			child := n.children[r]
			if child == nil {
				if n.children == nil {
					n.children = map[rune]*node{}
				}
				child = &node{}
				n.children[r] = child
			}
			n = child
		}
	}

	n.hasCanonical = true
	n.canonical = canonical
}

func (trie *runeTrie) SearchCanonical(tokens ...*jargon.Token) (found bool, canonical string, consumed int) {
	var result *node
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
		}

		if node.hasCanonical && node != result {
			// only capture results if it's a different node
			result = node
			found = true
			canonical = node.canonical
			consumed = i + 1
		}
	}

	return found, canonical, consumed
}
