package stemmer

import "testing"

func TestEnglish(t *testing.T) {
	// Just testing one of the filters to test the lookup logic; have to defer
	// to kljensen/snowball on correctness of other language stemmers

	type test struct {
		// input
		word string
		// expected
		stem    string
		stemmed bool
	}

	tests := []test{
		{"Accumulations", "accumul", true},
		{"hello", "hello", false},
	}

	for _, test := range tests {
		stem, stemmed := English.Lookup(test.word)

		if test.stemmed != stemmed {
			t.Errorf("expected stemmed %t, got %t", test.stemmed, stemmed)
		}

		if test.stem != stem {
			t.Errorf("expected stem of %s to be %s, got %s", test.word, test.stem, stem)
		}
	}
}
