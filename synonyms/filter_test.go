package synonyms_test

import (
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/synonyms"
)

func TestFilter(t *testing.T) {

	mappings := map[string]string{
		"developer, engineer, programmer": "boffin",
		"rock star, 10x developer":        "cliché",
		"Ruby on Rails, rails":            "ruby-on-rails",
		"nodeJS, iojs":                    "node.js",
	}

	ignore := []rune{'-', ' ', '.', '/'}
	filter, err := synonyms.NewFilter(mappings, true, ignore)
	if err != nil {
		t.Error(err)
	}

	// t.Log(filter.Decl())
	//t.Logf("%#v", syns.Trie)
	original := `we are looking for a rockstar 10x developer or engineer for ruby on rails and Nodejs`
	tokens := jargon.TokenizeString(original)

	expected := `we are looking for a cliché cliché or boffin for ruby-on-rails and node.js`

	got, err := filter.Filter(tokens).String()
	if err != nil {
		t.Error(err)
	}

	if expected != got {
		t.Errorf("given %q, expected %q, got %q", original, expected, got)
	}
}

func BenchmarkFilter(b *testing.B) {
	// file, err := ioutil.ReadFile("../testdata/wikipedia.txt")

	// if err != nil {
	// 	b.Error(err)
	// }

	// ignore := []rune{'-', ' ', '.', '/'}
	// syns, err := synonyms.NewFilter(mappings, true, ignore)
	// if err != nil {
	// 	b.Error(err)
	// }

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	r := bytes.NewReader(file)
	// 	tokens := jargon.Tokenize(r)
	// 	syns.Filter(tokens).Count() // consume them
	// }
}
