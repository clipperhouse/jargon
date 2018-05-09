package techlemm

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := map[string]string{
		"foo.js":      "foojs",
		".net":        ".net",
		"asp.net-mvc": "aspnetmvc",
	}

	for given, expected := range tests {
		got := normalize(given)
		if got != expected {
			t.Errorf("Given %q, expected %q, but got %q", given, expected, got)
		}
	}
}

func TestNewTagMap(t *testing.T) {
	// Intended to narrowly test that the values have been added to the data structure

	tags := []string{"Node.js", "ASP.net"}
	synonyms := map[string]string{
		"io.js":      "Node.js",
		"ECMAScript": "JavaScript",
	}
	tagmap := NewTagMap(tags, synonyms)

	for _, value := range tags {
		key := normalize(value)
		got, exists := tagmap.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", value, exists)
		}
		if got != value {
			t.Errorf("Given added tag %q, expected to retrieve same, but got %q", value, got)
		}
	}

	for synonym, canonical := range synonyms {
		key := normalize(synonym)
		got, exists := tagmap.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", canonical, exists)
		}
		if got != canonical {
			t.Errorf("Given added tag %q, expected to retrieve same, but got %q", canonical, got)
		}
	}
}

func TestGet(t *testing.T) {
	tags := []string{"Node.js", "ASP.net"}
	synonyms := map[string]string{
		"io.js":      "Node.js",
		"ECMAScript": "JavaScript",
	}
	tagmap := NewTagMap(tags, synonyms)

	type test struct {
		input, expected string
		found           bool
	}

	tests := []test{
		{"nodejs", "Node.js", true},
		{"IOjs", "Node.js", true},
		{"foo", "", false},
	}

	for _, test := range tests {
		got, found := tagmap.Get(test.input)
		if found != test.found {
			t.Errorf("Given input %q, expected found to be true, but got %t", test.input, found)
		}
		if got != test.expected {
			t.Errorf("Given input %q, expected get %q, but got %q", test.input, test.expected, got)
		}
	}
}
