package contractions

import (
	"strings"
	"testing"
)

func TestSome(t *testing.T) {
	given := "i'll SHE’D they're Can’t should've GOTTA Wanna"
	expected := "i will SHE WOULD they are Can not should have GOT TO Want to"

	var lookups []string

	for _, word := range strings.Split(given, " ") {
		canonical, ok := Expander.Lookup([]string{word})
		if ok {
			lookups = append(lookups, canonical)
		}
	}

	got := strings.Join(lookups, " ")

	if got != expected {
		t.Errorf("given %q, expected %q, got %q", given, expected, got)
	}
}
