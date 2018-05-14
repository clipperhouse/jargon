package jargon

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon/tokenizers"

	"github.com/clipperhouse/jargon/stackexchange"
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

var testDict = stackexchange.Dictionary
var testLem = NewLemmatizer(testDict)

func TestLemmatizer(t *testing.T) {
	// Intended to narrowly test that the values have been added to the data structure
	tags := testDict.GetTags()
	for _, value := range tags {
		key := normalize(value)
		_, exists := testLem.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", value, exists)
		}
	}

	synonyms := testDict.GetSynonyms()
	for synonym, canonical := range synonyms {
		key := normalize(synonym)
		_, exists := testLem.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", canonical, exists)
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
	tokens := strings.Split("This is the story of Ruby on Rails nodeJS and ASPNET mvc", " ")
	lemmatized := testLem.Lemmatize(tokens)
	got := strings.Join(lemmatized, " ")
	expected := "This is the story of ruby-on-rails node.js and asp.net-mvc"
	if got != expected {
		t.Errorf("Given tokens %v, expected %q, but got %q", tokens, expected, got)
	}
}

func TestLemmatizeTokens(t *testing.T) {
	text := "This is the story of Ruby on Rails nodeJS and ASPNET mvc"
	tokens := tokenizers.TechProse.Tokenize(text)
	lemmatized := testLem.LemmatizeTokens(tokens)
	got := tokenizers.Join(lemmatized, tokenizers.Token.Value)
	expected := "This is the story of ruby-on-rails node.js and asp.net-mvc"
	if got != expected {
		t.Errorf("Given tokens %v, expected %q, but got %q", text, expected, got)
	}
}
