package jargon

import (
	"testing"

	"github.com/clipperhouse/jargon/stackexchange"
)

var testDict = stackexchange.Dictionary
var testLem = NewLemmatizer(testDict)

func TestLemmatizer(t *testing.T) {
	// Intended to narrowly test that the values have been added to the data structure
	tags := testDict.Lemmas()
	for _, value := range tags {
		key := testDict.Normalize(value)
		_, exists := testLem.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", value, exists)
		}
	}

	synonyms := testDict.Synonyms()
	for synonym, canonical := range synonyms {
		key := testDict.Normalize(synonym)
		_, exists := testLem.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", canonical, exists)
		}
	}
}

func TestLemmatizeTokens(t *testing.T) {
	text := "Here is the story of Ruby on Rails nodeJS and ASPNET mvc plus TCP/IP."
	tokens := TechProse.Tokenize(text)
	lemmatized := testLem.LemmatizeTokens(tokens)
	got := Join(lemmatized)
	expected := "Here is the story of ruby-on-rails node.js and asp.net-mvc plus tcp."
	if got != expected {
		t.Errorf("Given tokens %v, expected %q, but got %q", text, expected, got)
	}
}

func TestCSV(t *testing.T) {
	text := `"Ruby on Rails", 3.4, "foo"
"bar",42, "java script"
`
	tokens := TechProse.Tokenize(text)
	lemmatized := testLem.LemmatizeTokens(tokens)
	got := Join(lemmatized)
	expected := `"ruby-on-rails", 3.4, "foo"
"bar",42, "javascript"
`

	if got != expected {
		t.Errorf("Given tokens %v, expected %q, but got %q", text, expected, got)
	}
}

func TestTabs(t *testing.T) {
	text := `Ruby on Rails	3.4	foo
ASPNET	MVC
bar	42	java script`

	tokens := TechProse.Tokenize(text)
	lemmatized := testLem.LemmatizeTokens(tokens)
	got := Join(lemmatized)
	expected := `ruby-on-rails	3.4	foo
asp.net	model-view-controller
bar	42	javascript`

	if got != expected {
		t.Errorf("Given tokens %v, expected %q, but got %q", text, expected, got)
	}
}
