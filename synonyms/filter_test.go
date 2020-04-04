package synonyms

import (
	"reflect"
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestFill(t *testing.T) {
	type test struct {
		// input
		input    string
		maxWords int
		previous *jargon.TokenQueue

		// expected
		buffer   *jargon.TokenQueue
		outgoing *jargon.TokenQueue
	}

	tests := []test{
		{
			input:    "test one",
			maxWords: 3,
			previous: &jargon.TokenQueue{},

			buffer: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("test", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("one", false),
				},
			},
			outgoing: &jargon.TokenQueue{},
		},
		{
			input:    "test two",
			maxWords: 1,
			previous: &jargon.TokenQueue{},

			buffer: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("test", false),
				},
			},
			outgoing: &jargon.TokenQueue{},
		},
		{
			input:    " test three",
			maxWords: 2,
			previous: &jargon.TokenQueue{},

			buffer: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("test", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("three", false),
				},
			},
			outgoing: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken(" ", false),
				},
			},
		},
		{
			input:    "test four, and five",
			maxWords: 4,
			previous: &jargon.TokenQueue{},

			buffer: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("test", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("four", false),
					jargon.NewToken(",", false),
				},
			},
			outgoing: &jargon.TokenQueue{
				Tokens: []*jargon.Token{},
			},
		},
		{
			input:    ", test six and seven",
			maxWords: 4,
			previous: &jargon.TokenQueue{},

			buffer: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("test", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("six", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("and", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("seven", false),
				},
			},
			outgoing: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken(",", false),
					jargon.NewToken(" ", false),
				},
			},
		},
		{
			input:    " test eight and nine",
			maxWords: 4,
			previous: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("previous", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("stuff", false),
				},
			},

			buffer: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("previous", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("stuff", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("test", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("eight", false),
				},
			},
			outgoing: &jargon.TokenQueue{
				Tokens: []*jargon.Token{},
			},
		},
		{
			input:    ". test ten and eleven",
			maxWords: 4,
			previous: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("previous", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("stuff", false),
				},
			},

			buffer: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("previous", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("stuff", false),
					jargon.NewToken(".", false),
				},
			},
			outgoing: &jargon.TokenQueue{
				Tokens: []*jargon.Token{},
			},
		},
		{
			input:    " test twelve and thirteen",
			maxWords: 3,
			previous: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken(".", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("leftover", false),
				},
			},

			buffer: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken("leftover", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("test", false),
					jargon.NewToken(" ", false),
					jargon.NewToken("twelve", false),
				},
			},
			outgoing: &jargon.TokenQueue{
				Tokens: []*jargon.Token{
					jargon.NewToken(".", false),
					jargon.NewToken(" ", false),
				},
			},
		},
	}

	for _, test := range tests {
		incoming := jargon.TokenizeString(test.input)
		tokens := &tokens{
			incoming: incoming,
			buffer:   test.previous,
			outgoing: &jargon.TokenQueue{},
			filter: &Filter{
				maxWords: test.maxWords,
			},
		}
		tokens.fill()

		expected := test.buffer.String()
		got := tokens.buffer.String()
		if expected != got {
			t.Errorf("expected %s, got %s", expected, got)
		}

		expected = test.outgoing.String()
		got = tokens.outgoing.String()
		if expected != got {
			t.Errorf("expected %s, got %s", expected, got)
		}
	}

}

func TestPassthrough(t *testing.T) {
	// If the filter doesn't do anything, the tokens should come back verbatim

	mappings := map[string]string{}
	ignore := []rune{}
	synonyms, err := NewFilter(mappings, false, ignore)
	if err != nil {
		t.Error(err)
	}

	text := "This is a test, with spaces and punctuation."

	original := jargon.TokenizeString(text)
	if err != nil {
		t.Error(err)
	}
	expected, err := original.ToSlice()
	if err != nil {
		t.Error(err)
	}

	filtered := jargon.TokenizeString(text)
	if err != nil {
		t.Error(err)
	}
	got, err := filtered.Filter(synonyms).ToSlice()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

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
