package jargon_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/clipperhouse/jargon"
)

// TODO: test ordering

func TestTokenize2(t *testing.T) {
	text := `Hi. 
	node.js, first_last, my.name@domain.com
	123.456, 789, 1,000, a16z, 3G and $200.13.
	wishy-washy and C++ and F#
	Let’s Let's possessive' possessive’
	Then ウィキペディア and 象形.`
	tokens := jargon.TokenizeString2(text)

	type test struct {
		value string
		found bool
	}

	expecteds := []test{
		{"Hi", true},
		{".", true},
		{"Hi.", false},

		{"node.js", true},
		{"node", false},
		{"js", false},

		{"first_last", true},
		{"first", false},
		{"_", false},
		{"last", false},

		{"my.name", true},
		{"my.name@", false},
		{"@", true},
		{"domain.com", true},
		{"@domain.com", false},

		{"123.456", true},
		{"123,", false},
		{"456", false},
		{"123.456,", false},

		{"wishy-washy", false},
		{"wishy", true},
		{"-", true},
		{"washy", true},

		{"C++", false},
		{"C", true},
		{"+", true},

		{"F#", false},
		{"F", true},
		{"#", true},

		{"Let's", true},
		{"Let’s", true},
		{"Let", false},
		{"s", false},

		{"possessive", true},
		{"'", true},
		{"’", true},
		{"possessive'", false},
		{"possessive’", false},

		{"789", true},
		{"789,", false},

		{"1,000", true},
		{"1,000,", false},

		{"a16z", true},

		{"3G", true},

		{"$", true},
		{"200.13", true},

		{"ウィキペディア", true},
		{"ウ", false},

		{"象", true},
		{"形", true},
		{"象形", false},
	}

	got := map[string]bool{}

	for {
		token, err := tokens.Next()
		if err != nil {
			t.Error(err)
		}
		if token == nil {
			break
		}

		got[token.String()] = true
	}

	for _, expected := range expecteds {
		if got[expected.value] != expected.found {
			t.Errorf("expected %q to be %t", expected.value, expected.found)
		}
	}
}

func BenchmarkTokenize2(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	var count int
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		c, err := jargon.Tokenize2(r).Count()
		if err != nil {
			b.Error(err)
		}
		count = c
	}
	b.Logf("%d tokens\n", count)
}
