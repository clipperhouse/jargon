package jargon_test

import (
	"strings"
	"testing"

	"github.com/blevesearch/segment"
	"github.com/clipperhouse/jargon"
)

// TODO: test ordering

func TestSegmenter(t *testing.T) {
	text := `Hi. This is, a very basic test of the segmenter—with node.js and first_last.
`

	r := strings.NewReader(text)
	segmenter := segment.NewSegmenter(r)

	got := map[string]bool{}

	for segmenter.Segment() {
		s := string(segmenter.Bytes())
		got[s] = true
	}

	type test struct {
		value string
		found bool
	}

	expecteds := []test{
		{"Hi", true},
		{".", true},
		{"Hi.", false},
		{"is", true},
		{",", true},
		{"—", true},
		{"node.js", true},
		{"node", false},
		{"js", false},
		{"first_last", true},
		{"first", false},
		{"_", false},
		{"last", false},
		{"\n", true},
	}

	for _, expected := range expecteds {
		if got[expected.value] != expected.found {
			t.Errorf("expected finding %q to be %t", expected.value, expected.found)
		}
	}
}

func TestLeading(t *testing.T) {
	text := `Hi. This is a test of .net, and #hashtag and @handle, and React.js.`

	r := strings.NewReader(text)
	tokens := jargon.Tokenize2(r)

	type test struct {
		value string
		found bool
	}

	expecteds := []test{
		{"Hi", true},
		{".", true},
		{"Hi.", false},
		{".net", true},
		{"net", false},
		{"#hashtag", true},
		{"hashtag", false},
		{"@handle", true},
		{"handle", false},
		{"React.js", true},
		{"React.js.", false},
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
			t.Errorf("expected finding %q to be %t", expected.value, expected.found)
		}
	}
}

func TestMiddle(t *testing.T) {
	text := `Hi. This is a test of https://foo.com:443/thing-stuff.html/ https://foo.com/?query=bar+baz&qux and some_name+spam@example.com and wishy-washy.`

	r := strings.NewReader(text)
	tokens := jargon.Tokenize2(r)

	type test struct {
		value string
		found bool
	}

	expecteds := []test{
		{"Hi", true},
		{".", true},
		{"Hi.", false},
		{"some_name+spam@example.com", true},
		{"_", false},
		{"+", false},
		{"@", false},
		{"example.com", false},
		{"wishy-washy", true},
		{"wishy", false},
		{"-", false},
		{"washy", false},
		{"wishy-washy", true},
		{"https://foo.com:443/thing-stuff.html/", true},
		{"https", false},
		{":", false},
		{"/", false},
		{"thing-stuff", false},
		{"html", false},
		{"https://foo.com/?query=bar+baz&qux", true},
		{"?", false},
		{"query", false},
		{"=", false},
		{"bar", false},
		{"baz", false},
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
			t.Errorf("expected finding %q to be %t", expected.value, expected.found)
		}
	}
}

func TestTrailing(t *testing.T) {
	text := `Hi. This is a test of F# and C++.`

	r := strings.NewReader(text)
	tokens := jargon.Tokenize2(r)

	type test struct {
		value string
		found bool
	}

	expecteds := []test{
		{"Hi", true},
		{".", true},
		{"Hi.", false},
		{"F#", true},
		{"F", false},
		{"#", false},
		{"C++", true},
		{"C", false},
		{"+", false},
		{"++", false},
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
			t.Errorf("expected finding %q to be %t", expected.value, expected.found)
		}
	}
}
