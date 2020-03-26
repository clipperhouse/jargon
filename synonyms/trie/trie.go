package trie

import (
	"bytes"
	"fmt"
	"go/format"
	"unicode"

	"github.com/clipperhouse/jargon"
)

type RuneTrie struct {
	Root       *Node
	Ignore     map[rune]bool
	IgnoreCase bool
}

func (n *RuneTrie) Decl() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "&trie.RuneTrie{\n")
	if n.Ignore != nil {
		fmt.Fprintf(&b, "Ignore: map[rune]bool{\n")
		for k, v := range n.Ignore {
			fmt.Fprintf(&b, "'%s': %t,\n", string(k), v)
		}
		fmt.Fprintf(&b, "},\n")
	}
	if n.IgnoreCase {
		// default value does not need to be declared
		fmt.Fprintf(&b, "IgnoreCase: %t,\n", n.IgnoreCase)
	}
	fmt.Fprintf(&b, "Root: %s,\n", n.Root.Decl())
	fmt.Fprintf(&b, "}")

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}

	return string(formatted)
}

type Node struct {
	Children     map[rune]*Node
	HasCanonical bool
	Canonical    string
}

func New(ignoreCase bool, ignore []rune) *RuneTrie {
	set := map[rune]bool{}
	for _, r := range ignore {
		set[r] = true
	}
	return &RuneTrie{
		Root:       &Node{},
		IgnoreCase: ignoreCase,
		Ignore:     set,
	}
}

func (trie *RuneTrie) Add(tokens []*jargon.Token, Canonical string) {
	n := trie.Root
	for _, token := range tokens {
		for _, r := range token.String() {
			if trie.IgnoreCase {
				r = unicode.ToLower(r)
			}

			if trie.Ignore[r] {
				continue
			}

			child := n.Children[r]
			if child == nil {
				if n.Children == nil {
					n.Children = map[rune]*Node{}
				}
				child = &Node{}
				n.Children[r] = child
			}
			n = child
		}
	}

	n.HasCanonical = true
	n.Canonical = Canonical
}

func (trie *RuneTrie) SearchCanonical(tokens ...*jargon.Token) (found bool, Canonical string, consumed int) {
	var result *Node
	Node := trie.Root

outer:
	for i, token := range tokens {
		for _, r := range token.String() {
			if trie.IgnoreCase {
				r = unicode.ToLower(r)
			}

			if trie.Ignore[r] {
				continue
			}

			Node = Node.Children[r]
			if Node == nil {
				break outer
			}
		}

		if Node.HasCanonical && Node != result {
			// only capture results if it's a different Node
			result = Node
			found = true
			Canonical = Node.Canonical
			consumed = i + 1
		}
	}

	return found, Canonical, consumed
}

func (n *Node) Decl() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "&trie.Node{\n")
	if n.HasCanonical {
		// default value does not need to be declared
		fmt.Fprintf(&b, "\tHasCanonical: %t,\n", n.HasCanonical)
	}
	if n.Canonical != "" {
		// default value does not need to be declared
		fmt.Fprintf(&b, "\tCanonical: %q,\n", n.Canonical)
	}
	if n.Children != nil {
		// default value does not need to be declared
		fmt.Fprintf(&b, "\tChildren: map[rune]*trie.Node{\n")
		for k, v := range n.Children {
			fmt.Fprintf(&b, "\t'%s': %s,\n", string(k), v.Decl())
		}
		fmt.Fprintf(&b, "},\n")
	}
	fmt.Fprintf(&b, "}")

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}

	return string(formatted)
}

func (n *Node) String() string {
	return n.Decl()
}
