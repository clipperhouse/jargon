package stackoverflow_test

import (
	"testing"

	"github.com/clipperhouse/jargon/stackoverflow"
)

func TestFilter(t *testing.T) {
	type test struct {
		input     string
		found     bool
		canonical string
	}
	expecteds := []test{
		{"c sharp", true, "c#"},
		{"Ruby on Rails", true, "ruby-on-rails"},
		{"Ruby", true, "ruby"},
		{"foo", false, ""},
	}

	for _, expected := range expecteds {
		canonical, found := stackoverflow.Tags.Lookup(expected.input)
		if found != expected.found {
			t.Errorf("found should be %t", expected.found)
		}
		if canonical != expected.canonical {
			t.Errorf("expcted to find canonical %q, got %q", expected.canonical, canonical)
		}
	}
}
