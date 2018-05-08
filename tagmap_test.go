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
