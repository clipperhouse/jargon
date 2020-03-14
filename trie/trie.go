package trie

type Trie struct {
	Children map[string]*Trie
	Value    string
}

func (trie *Trie) Add(ss []string, value string) bool {
	node := trie
	for _, s := range ss {
		child := node.Children[s]
		if child == nil {
			if node.Children == nil {
				node.Children = map[string]*Trie{}
			}
			child = &Trie{}
			node.Children[s] = child
		}
		node = child
	}
	// does node have an existing value?
	isNewVal := node.Value == ""
	node.Value = value
	return isNewVal
}

func (trie *Trie) Search(ss ...string) *Trie {
	node := trie
	for _, s := range ss {
		node = node.Children[s]
		if node == nil {
			return nil
		}
	}
	return node
}

func (trie *Trie) SearchValue(ss ...string) *Trie {
	var result *Trie

	node := trie
	for _, s := range ss {
		node = node.Children[s]
		if node == nil {
			break
		}
		if node.Value != "" {
			result = node
		}
	}
	return result
}
