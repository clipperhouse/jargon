package stemmer

import (
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestEnglish(t *testing.T) {
	// Just testing one of the filters to test the lookup logic; have to defer
	// to kljensen/snowball on correctness of language stemmers

	type test struct {
		// input
		input string
		// expected
		output string
	}

	tests := []test{
		{"Accumulations are expected", "accumul are expect"},
		{"hello it's me", "hello it me"},
	}

	for _, test := range tests {
		tokens := jargon.TokenizeString(test.input)
		got, err := English(tokens).String()
		if err != nil {
			t.Error(err)
		}

		if test.output != got {
			t.Errorf("expected stem of %q to be %q, got %q", test.input, test.output, got)
		}
	}
}
