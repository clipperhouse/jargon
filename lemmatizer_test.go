package jargon

import (
	"reflect"
	"testing"

	"github.com/clipperhouse/jargon/stackexchange"
)

func TestNewLemmatizer(t *testing.T) {
	// Intended to narrowly test that the values have been added to the data structure

	dict := stackexchange.Dictionary
	lem := NewLemmatizer(dict)

	for _, value := range dict.Lemmas() {
		key := dict.Normalize(value)
		_, exists := lem.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", value, exists)
		}
	}

	for synonym, canonical := range dict.Synonyms() {
		key := dict.Normalize(synonym)
		_, exists := lem.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists to be true, but got %t", canonical, exists)
		}
	}
}

func TestLemmatizeTokens(t *testing.T) {
	tok := TechProse.Tokenize
	dict := stackexchange.Dictionary
	lem := NewLemmatizer(dict)

	original := `Here is the story of Ruby on Rails nodeJS, "Java Script" and ASPNET mvc plus TCP/IP.`
	tokens := tok(original)
	got := lem.LemmatizeTokens(tokens)
	expected := tok(`Here is the story of ruby-on-rails node.js, "javascript" and asp.net-mvc plus tcp.`)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func TestCSV(t *testing.T) {
	tok := TechProse.Tokenize
	dict := stackexchange.Dictionary
	lem := NewLemmatizer(dict)

	original := `"Ruby on Rails", 3.4, "foo"
"bar",42, "java script"`
	tokens := tok(original)
	got := lem.LemmatizeTokens(tokens)
	expected := tok(`"ruby-on-rails", 3.4, "foo"
"bar",42, "javascript"`)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func TestTSV(t *testing.T) {
	tok := TechProse.Tokenize
	dict := stackexchange.Dictionary
	lem := NewLemmatizer(dict)

	original := `Ruby on Rails	3.4	foo
ASPNET	MVC
bar	42	java script`

	tokens := tok(original)
	got := lem.LemmatizeTokens(tokens)
	expected := tok(`ruby-on-rails	3.4	foo
asp.net	model-view-controller
bar	42	javascript`)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func TestWordrun(t *testing.T) {
	tok := TechProse.Tokenize
	original := `Things and "java script"`
	tokens := tok(original)

	type result struct {
		expect   []string
		consumed int
		ok       bool
	}

	takes := []int{3, 2, 1}

	expecteds := map[int]result{
		3: {[]string{}, 0, false},                // attempting to get 3 should fail
		2: {[]string{"java", "script"}, 3, true}, // attempting to get 2 should work, consuming 3 tokens (incl the space)
		1: {[]string{"java"}, 1, true},           // attempting to get 1 should work, and consume only that token
	}

	for _, take := range takes {
		taken, consumed, ok := wordrun(tokens, 5, take) // 5 = start at the 'j' of 'java'
		got := result{Strings(taken), consumed, ok}
		expected, _ := expecteds[take]

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Attempting to take %d words, expected %v but got %v", take, expected, got)
		}
	}
}
