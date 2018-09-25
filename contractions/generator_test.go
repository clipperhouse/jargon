package contractions

import (
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	if err := write(); err != nil {
		t.Error(err)
	}
}

func TestVariations(t *testing.T) {
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

	got := len(variations)

	if got != expected {
		t.Errorf("generated variations should have %d items, but got %d", expected, got)
	}
}
