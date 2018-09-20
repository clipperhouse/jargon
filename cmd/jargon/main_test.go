package main

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/numbers"
	"github.com/clipperhouse/jargon/stackexchange"
)

func TestDetermineLemmatizers(t *testing.T) {
	defaults := determineLemmatizers(false, false)
	if len(defaults) != 1 {
		t.Errorf("expected 1 lemmatizer, got %d", len(defaults))
	}
	if defaults[0].Dictionary != stackexchange.Dictionary {
		t.Errorf("expected default to be the stackexchange.Dictionary, got %T", defaults[0].Dictionary)
	}

	tech := determineLemmatizers(true, false)
	if len(tech) != 1 {
		t.Errorf("expected 1 lemmatizer, got %d", len(tech))
	}
	if tech[0].Dictionary != stackexchange.Dictionary {
		t.Errorf("expected tech lemmatizer to include stackexchange.Dictionary, got a %T", tech[0].Dictionary)
	}

	num := determineLemmatizers(false, true)
	if len(num) != 1 {
		t.Errorf("expected 1 lemmatizer when num is specified, got %d", len(num))
	}
	if num[0].Dictionary != numbers.Dictionary {
		t.Errorf("expected num lemmatizer to include numbers.Dictionary, got a %T", num[0].Dictionary)
	}

	both := determineLemmatizers(true, true)
	if len(both) != 2 {
		t.Errorf("expected 2 lemmatizer when tech and num are specified, got %d", len(both))
	}
	if both[0].Dictionary != stackexchange.Dictionary {
		t.Errorf("expected first lemmatizer to be stackexchange.Dictionary, got a %T", both[0].Dictionary)
	}
	if both[1].Dictionary != numbers.Dictionary {
		t.Errorf("expected second lemmatizer to be numbers.Dictionary, got a %T", both[1].Dictionary)
	}
}

func TestLemAll(t *testing.T) {
	s := "I luv Rails times three hundred"
	r := strings.NewReader(s)
	tokens := jargon.Tokenize(r)
	lems := []*jargon.Lemmatizer{
		jargon.NewLemmatizer(stackexchange.Dictionary, 3),
		jargon.NewLemmatizer(numbers.Dictionary, 3),
	}

	lemmatized := lemAll(tokens, lems)
	got := ""

	for {
		t := lemmatized.Next()
		if t == nil {
			break
		}
		got += t.String()
	}

	expected := "I luv ruby-on-rails times 300"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
