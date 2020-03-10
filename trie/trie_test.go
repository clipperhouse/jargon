package trie

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	trie := &Trie{
		Value: "",
		Children: map[string]*Trie{
			"ruby": &Trie{
				Value: "ruby",
				Children: map[string]*Trie{
					"on": &Trie{
						Value: "",
						Children: map[string]*Trie{
							"rails": &Trie{
								Value: "ruby-on-rails",
							},
						},
					},
				},
			},
		},
	}
	t.Log(trie.Search("ruby", "on"))
	t.Log(trie.Search("ruby", "the"))
	t.Log(trie.SearchValue("ruby"))
	t.Log(trie.SearchValue("ruby", "on"))
	t.Log(trie.SearchValue("ruby", "on", "rails"))
	t.Log(trie.SearchValue("foo", "on", "rails"))
}

func TestAdd(t *testing.T) {
	trie := &Trie{}
	trie.Add([]string{"ruby", "on", "rails"}, "ruby-on-rails")
	trie.Add([]string{"react", "js"}, "react.js")
	printTrie(trie)
	t.Log(trie.Search("ruby", "on"))
	t.Log(trie.Search("ruby", "the"))
	t.Log(trie.SearchValue("ruby"))
	t.Log(trie.SearchValue("ruby", "on"))
	t.Log(trie.SearchValue("ruby", "on", "rails").Value)
	t.Log(trie.SearchValue("foo", "on", "rails"))
}

func printTrie(t *Trie) {
	for k, v := range t.Children {
		fmt.Println(k)
		fmt.Println(v.Value)
		printTrie(v)
	}
}
