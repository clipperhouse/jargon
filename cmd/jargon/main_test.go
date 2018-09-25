package main

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/numbers"
	"github.com/clipperhouse/jargon/stackexchange"
)

func TestDefaultLemmatizer(t *testing.T) {
	defaults := determineLemmatizers(false, false, false)

	if len(defaults) != 1 {
		t.Errorf("expected 1 lemmatizer, got %d", len(defaults))
	}

	if defaults[0].Dictionary != stackexchange.Dictionary {
		t.Errorf("expected default to be the stackexchange.Dictionary, got %T", defaults[0].Dictionary)
	}
}

func TestTechLemmatizer(t *testing.T) {
	tech := determineLemmatizers(true, false, false)

	if len(tech) != 1 {
		t.Errorf("expected 1 lemmatizer, got %d", len(tech))
	}

	if tech[0].Dictionary != stackexchange.Dictionary {
		t.Errorf("expected tech lemmatizer to include stackexchange.Dictionary, got a %T", tech[0].Dictionary)
	}
}

func TestNumbersLemmatizer(t *testing.T) {
	num := determineLemmatizers(false, true, false)

	if len(num) != 1 {
		t.Errorf("expected 1 lemmatizer when num is specified, got %d", len(num))
	}

	if num[0].Dictionary != numbers.Dictionary {
		t.Errorf("expected num lemmatizer to include numbers.Dictionary, got a %T", num[0].Dictionary)
	}
}

func TestContractionsLemmatizer(t *testing.T) {
	cont := determineLemmatizers(false, false, true)

	if len(cont) != 1 {
		t.Errorf("expected 1 lemmatizer when cont is specified, got %d", len(cont))
	}

	if cont[0].Dictionary != contractions.Dictionary {
		t.Errorf("expected cont lemmatizer to include contractions.Dictionary, got a %T", cont[0].Dictionary)
	}
}

func TestAllLemmatizers(t *testing.T) {
	all := determineLemmatizers(true, true, true)

	if len(all) != 3 {
		t.Errorf("expected 3 lemmatizers when tech and num and cont are specified, got %d", len(all))
	}

	if all[0].Dictionary != stackexchange.Dictionary {
		t.Errorf("expected first lemmatizer to be stackexchange.Dictionary, got a %T", all[0].Dictionary)
	}

	if all[1].Dictionary != numbers.Dictionary {
		t.Errorf("expected second lemmatizer to be numbers.Dictionary, got a %T", all[1].Dictionary)
	}

	if all[2].Dictionary != contractions.Dictionary {
		t.Errorf("expected second lemmatizer to be contractions.Dictionary, got a %T", all[1].Dictionary)
	}
}

func TestLemAll(t *testing.T) {
	s := "I can't luv Rails times three hundred"
	r := strings.NewReader(s)
	tokens := jargon.Tokenize(r)
	lems := []*jargon.Lemmatizer{
		jargon.NewLemmatizer(stackexchange.Dictionary, 3),
		jargon.NewLemmatizer(numbers.Dictionary, 3),
		jargon.NewLemmatizer(contractions.Dictionary, 3),
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

	expected := "I can not luv ruby-on-rails times 300"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
