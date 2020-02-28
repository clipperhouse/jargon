package stackexchange

import (
	"testing"
)

func TestFilter(t *testing.T) {
	term := "rails"
	canonical, found := Tags.Lookup(term)
	if !found {
		t.Errorf("should have found %q", term)
	}
	if !found {
		t.Errorf("should have found canonical %q", canonical)
	}
	t.Log(found)
	t.Log(canonical)
}
