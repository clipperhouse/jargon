package jargon

import (
	"testing"

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

func TestLemmatizeTokens(t *testing.T) {
	text := "This is the story of Ruby on Rails nodeJS and ASPNET mvc"
	tokens := TechProse.Tokenize(text)
	lemmatized := testLem.LemmatizeTokens(tokens)
	got := Join(lemmatized, Token.Value)
	expected := "This is the story of ruby-on-rails node.js and asp.net-mvc"
	if got != expected {
		t.Errorf("Given tokens %v, expected %q, but got %q", text, expected, got)
	}
}
