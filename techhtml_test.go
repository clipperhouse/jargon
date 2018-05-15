package jargon

import (
	"testing"
)

func TestTechHTML(t *testing.T) {
	h := `
<html>
<p foo="bar">
Hi! Let's talk Ruby on Rails.
<!-- Ignore ASPNET MVC in comments -->
</p>
</html>
`
	got := TechHTML.Tokenize(h)

	expected := []string{
		`<p foo="bar">`, // tags kept whole
		"\n",            // whitespace preserved
		"Hi", "!",
		"Ruby", "on", "Rails", // make sure text node got tokenized
		"<!-- Ignore ASPNET MVC in comments -->", // make sure comment kept whole
		"</p>",
	}

	for _, e := range expected {
		if !contains(e, got) {
			t.Errorf("Expected to find token %q, but did not.", e)
		}
	}
}
