package synonyms

import (
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestFilter(t *testing.T) {
	mappings := map[string]string{
		"developer, engineer, programmer,": "boffin",
		"rock star, 10x developer":         "cliché",
		"Ruby on Rails, rails":             "ruby-on-rails",
		"nodeJS, iojs":                     "node.js",
	}

	ignore := []rune{'-', ' ', '.', '/'}
	synonyms, err := NewFilter(mappings, true, ignore)
	if err != nil {
		t.Error(err)
	}

	original := `we are looking for a rockstar, 10x developer, or engineer, for ruby on rails and Nodejs`
	tokens := jargon.TokenizeString(original)

	expected := `we are looking for a cliché, cliché, or boffin, for ruby-on-rails and node.js`

	got, err := tokens.Filter(synonyms).String()
	if err != nil {
		t.Error(err)
	}

	if expected != got {
		t.Errorf("given %q, expected %q, got %q", original, expected, got)
	}
}

func BenchmarkFilter(b *testing.B) {

	mappings := map[string]string{
		"developer, engineer, programmer,": "boffin",
		"rock star, 10x developer":         "cliché",
		"Ruby on Rails, rails":             "ruby-on-rails",
		"nodeJS, iojs":                     "node.js",
	}

	ignore := []rune{'-', ' ', '.', '/'}
	filter, err := NewFilter(mappings, true, ignore)
	if err != nil {
		b.Error(err)
	}

	original := `we are looking for a rockstar 10x developer or engineer for ruby on rails and Nodejs`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokens := jargon.TokenizeString(original)
		_, err := tokens.Filter(filter).Count()
		if err != nil {
			b.Error(err)
		}
	}
}
