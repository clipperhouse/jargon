package synonyms_test

import (
	"testing"

	"github.com/clipperhouse/jargon/synonyms"
)

func TestBasics(t *testing.T) {
	expecteds := []test{
		{"engineer", true, "developer"},
		{" engineer", false, ""},
		{"synonym", true, "synonym"},
		{"the same", true, "synonym"},
		{"thesame", false, ""},
	}

	mappings := []synonyms.Mapping{
		{"developer, engineer, programmer, some word", "developer"},
		{"synonym, equivalent, the same", "synonym"},
	}
	filter := synonyms.New(mappings)
	testSynonyms(t, filter, expecteds)
}

func TestSpaceAndCase(t *testing.T) {
	expecteds := []test{
		{"ecmascript", true, "javascript"},
		{"ecma script", true, "javascript"},
		{"In my humble opinion", true, "imo"},
		{"imho", true, "imo"},
	}

	mappings := []synonyms.Mapping{
		{"ecma script, java script, js", "javascript"},
		{"in my opinion, in my humble  opinion, IMHO", "imo"},
	}
	filter := synonyms.New(mappings, synonyms.IgnoreSpace, synonyms.IgnoreCase)

	testSynonyms(t, filter, expecteds)
}

type test struct {
	input     string
	found     bool
	canonical string
}

func testSynonyms(t *testing.T, filter *synonyms.Filter, expecteds []test) {
	for _, expected := range expecteds {
		canonical, found := filter.Lookup(expected.input)
		if expected.found != found {
			t.Errorf("given input %q, expected found to be %t", expected.input, expected.found)
		}
		if expected.canonical != canonical {
			t.Errorf("given input %q, expected canonical %q, got %q", expected.input, expected.canonical, canonical)
		}
	}
}