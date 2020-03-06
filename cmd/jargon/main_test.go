package main

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/numbers"
	"github.com/clipperhouse/jargon/stackoverflow"
)

func TestDefaultLemmatizer(t *testing.T) {
	defaults := determineFilters(false, false, false)

	if len(defaults) != 1 {
		t.Errorf("expected 1 lemmatizer, got %d", len(defaults))
	}

	if defaults[0] != stackoverflow.Tags {
		t.Errorf("expected default to be the stackoverflow.Tags, got %T", defaults[0])
	}
}

func TestTechLemmatizer(t *testing.T) {
	tech := determineFilters(true, false, false)

	if len(tech) != 1 {
		t.Errorf("expected 1 lemmatizer, got %d", len(tech))
	}

	if tech[0] != stackoverflow.Tags {
		t.Errorf("expected tech lemmatizer to include stackoverflow.Tags, got a %T", tech[0])
	}
}

func TestNumbersLemmatizer(t *testing.T) {
	num := determineFilters(false, true, false)

	if len(num) != 1 {
		t.Errorf("expected 1 lemmatizer when num is specified, got %d", len(num))
	}

	if num[0] != numbers.Filter {
		t.Errorf("expected num lemmatizer to include numbers.Filter, got a %T", num[0])
	}
}

func TestContractionsLemmatizer(t *testing.T) {
	cont := determineFilters(false, false, true)

	if len(cont) != 1 {
		t.Errorf("expected 1 lemmatizer when cont is specified, got %d", len(cont))
	}

	if cont[0] != contractions.Expander {
		t.Errorf("expected cont lemmatizer to include contractions.Expand, got a %T", cont[0])
	}
}

func TestAllLemmatizers(t *testing.T) {
	all := determineFilters(true, true, true)

	if len(all) != 3 {
		t.Errorf("expected 3 lemmatizers when tech and num and cont are specified, got %d", len(all))
	}

	if all[0] != stackoverflow.Tags {
		t.Errorf("expected first lemmatizer to be stackoverflow.Tags, got a %T", all[0])
	}

	if all[1] != numbers.Filter {
		t.Errorf("expected second lemmatizer to be numbers.Filter, got a %T", all[1])
	}

	if all[2] != contractions.Expander {
		t.Errorf("expected second lemmatizer to be contractions.Expander, got a %T", all[1])
	}
}

func TestLemAll(t *testing.T) {
	s := "I can't luv Rails times three hundred"
	r := strings.NewReader(s)
	tokens := jargon.Tokenize(r)

	filters := []jargon.TokenFilter{
		stackoverflow.Tags,
		numbers.Filter,
		contractions.Expander,
	}

	lemmatized := lemAll(tokens, filters)
	got, err := lemmatized.String()

	if err != nil {
		t.Error(err)
	}

	expected := "I can not luv ruby-on-rails times 300"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
