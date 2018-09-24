package contractions

import (
	"strings"
	"testing"
)

func TestExhaustive(t *testing.T) {
	given := "i'll she’d they'll wouldn’t should've"
	expected := "i will she would they will would not should have"

	var lookups []string

	for _, word := range strings.Split(given, " ") {
		canonical, ok := Dictionary.Lookup([]string{word})
		if ok {
			lookups = append(lookups, canonical)
		}
	}

	got := strings.Join(lookups, " ")

	if got != expected {
		t.Errorf("given %q, expected %q, got %q", given, expected, got)
	}
}
