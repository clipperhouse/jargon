package main

import (
	"strings"
	"testing"
)

func TestSanity(t *testing.T) {
	for contraction, expansion := range contractions {
		// first letters should match, intended to catch dumb typos
		if contraction[0] != expansion[0] {
			t.Errorf("the first character of the mapping %q → %q should match", contraction, expansion)
		}
	}
}

func TestMappings(t *testing.T) {
	var n1, n2 int

	for key := range contractions {
		if strings.Contains(key, "'") {
			n1++
		} else {
			n2++
		}
	}

	expected := 3 * 2 * n1 // three cases (lower, upper, title) times two apostrophe variations (', ’)
	expected += 3 * n2     // three cases (lower, upper, title), apostrophes irrelevant

	mappings, err := getMappings()
	if err != nil {
		t.Error(err)
	}

	got := len(mappings)

	if got != expected {
		t.Errorf("generated variations should have %d items, but got %d", expected, got)
	}
}
