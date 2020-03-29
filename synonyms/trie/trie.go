package trie

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/clipperhouse/jargon"
)

// RuneTrie is a structure for searching strings (tokens)
type RuneTrie struct {
	root       *node
	ignore     map[rune]bool
	ignoreCase bool
}

// New creates a new RuneTrie
func New(ignoreCase bool, ignore []rune) *RuneTrie {
	set := map[rune]bool{}
	for _, r := range ignore {
		set[r] = true
	}
	return &RuneTrie{
		root:       &node{},
		ignoreCase: ignoreCase,
		ignore:     set,
	}
}

// String returns a representation of the trie as a Go source declaration. It can be large, use sparingly.
func (t *RuneTrie) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "&trie.RuneTrie{\n")
	if t.ignore != nil {
		fmt.Fprintf(&b, "Ignore: map[rune]bool{\n")
		for k, v := range t.ignore {
			fmt.Fprintf(&b, "'%s': %t,\n", string(k), v)
		}
		fmt.Fprintf(&b, "},\n")
	}
	if t.ignoreCase {
		// default value does not need to be declared
		fmt.Fprintf(&b, "IgnoreCase: %t,\n", t.ignoreCase)
	}
	fmt.Fprintf(&b, "Root: %s,\n", t.root.String())
	fmt.Fprintf(&b, "}")

	result := b.Bytes()
	// result, err := format.Source(result)
	// if err != nil {
	// 	panic(err)
	// }

	return string(result)
}

type node struct {
	children     map[rune]*node
	hasCanonical bool
	canonical    string
}

// Add adds tokens and their canonicals to the trie
func (t *RuneTrie) Add(tokens []*jargon.Token, canonical string) {
	n := t.root
	for _, token := range tokens {
		for _, r := range token.String() {
			if t.ignoreCase {
				r = unicode.ToLower(r)
			}

			if t.ignore[r] {
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

// SearchCanonical walks the trie to find a canonical matching the tokens, preferring longer (greedy) matches, i.e. 'ruby on rails' vs 'ruby'
func (t *RuneTrie) SearchCanonical(tokens ...*jargon.Token) (found bool, canonical string, consumed int) {
	var result *node
	n := t.root

outer:
	for i, token := range tokens {
		for _, r := range token.String() {
			if t.ignoreCase {
				r = unicode.ToLower(r)
			}

			if t.ignore[r] {
				continue
			}

			n = n.children[r]
			if n == nil {
				break outer
			}
		}

		if n.hasCanonical && n != result {
			// only capture results if it's a different node
			result = n
			found = true
			canonical = n.canonical
			consumed = i + 1
		}
	}

	return found, canonical, consumed
}

// String returns a representation of a node as a Go source declaration. It can be very large.
func (n *node) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "&trie.Node{\n")
	if n.hasCanonical {
		// default value does not need to be declared
		fmt.Fprintf(&b, "\tHasCanonical: %t,\n", n.hasCanonical)
	}
	if n.canonical != "" {
		// default value does not need to be declared
		fmt.Fprintf(&b, "Canonical: %q,\n", n.canonical)
	}
	if n.children != nil {
		// default value does not need to be declared
		fmt.Fprintf(&b, "Children: map[rune]*trie.Node{\n")
		for r, child := range n.children {
			fmt.Fprintf(&b, "'%s': %s,\n", string(r), child.String())
		}
		fmt.Fprintf(&b, "},\n")
	}
	fmt.Fprintf(&b, "}")

	result := b.Bytes()
	// result, err := format.Source(result)
	// if err != nil {
	// 	panic(err)
	// }

	return string(result)
}
