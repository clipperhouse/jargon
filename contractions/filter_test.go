package contractions_test

import (
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/contractions"
)

func TestContractions(t *testing.T) {
	given := "i'll SHE’D they're Can’t should've GOTTA Wanna"
	expected := "i will SHE WOULD they are Can not should have GOT TO Want to"

	tokens := jargon.TokenizeString(given)
	got, err := contractions.Expander.Filter(tokens).String()
	if err != nil {
		t.Error(err)
	}

	if got != expected {
		t.Errorf("given %q, expected %q, got %q", given, expected, got)
	}
}
