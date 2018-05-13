package tokenizers

import (
	"reflect"
	"testing"
)

func TestWhiteSpaceDelimited(t *testing.T) {
	text := `
This thing has  spaces and line breaks
and a	tab
`

	got := WhiteSpace.Tokenize(text)
	expected := []string{
		"\n",
		"This", " ", "thing", " ", "has", " ", " ", "spaces", " ", "and", " ", "line", " ", "breaks", "\n",
		"and", " ", "a", "\t", "tab", "\n",
	}

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected\n%v,\ngot\n%v\n", expected, got)
	}
}
