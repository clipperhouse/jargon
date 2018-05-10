package techlemm

import (
	"strings"
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

var testTags = []string{"Node.js", "ASP.net"}
var testSynonyms = map[string]string{
	"io.js":      "Node.js",
	"ECMAScript": "JavaScript",
}
var testDict = NewDictionary(testTags, testSynonyms)
var testLem = NewLemmatizer(testDict)

func TestLemmatizer(t *testing.T) {
	// Intended to narrowly test that the values have been added to the data structure
	for _, value := range testTags {
		key := normalize(value)
		got, exists := testLem.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", value, exists)
		}
		if got != value {
			t.Errorf("Given added tag %q, expected to retrieve same, but got %q", value, got)
		}
	}

	for synonym, canonical := range testSynonyms {
		key := normalize(synonym)
		got, exists := testLem.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", canonical, exists)
		}
		if got != canonical {
			t.Errorf("Given added tag %q, expected to retrieve same, but got %q", canonical, got)
		}
	}
}

func TestGetCanonical(t *testing.T) {
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
		got, found := testLem.GetCanonical(test.input)
		if found != test.found {
			t.Errorf("Given input %q, expected found to be true, but got %t", test.input, found)
		}
		if !test.found && got != test.expected { // if test doesn't expect it to be found, don't test value
			t.Errorf("Given input %q, expected get %q, but got %q", test.input, test.expected, got)
		}
	}
}

func TestLemmatize(t *testing.T) {
	tokens := strings.Split("This is the story of nodeJS and ASPNET", " ")
	lemmatized := testLem.Lemmatize(tokens)
	got := strings.Join(lemmatized, " ")
	expected := "This is the story of Node.js and ASP.net"
	if got != expected {
		t.Errorf("Given tokens %v, expected get %q, but got %q", tokens, expected, got)
	}
}
