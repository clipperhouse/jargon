package jargon

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTechProse(t *testing.T) {
	text := `Hi! This is a test of tech terms.
It should consider F#, C++, .net, Node.JS and 3.141592 to be their own tokens. 
Similarly, #hashtag and @handle should work, as should an first.last+@example.com.
It should—wait for it—break on things like em-dashes and "quotes" and it ends.
It'd be great it it’ll handle apostrophes.`
	got := TechProse.Tokenize(text)

	expected := []string{
		"Hi", "!",
		"F#", "C++", ".net", "Node.JS", "3.141592",
		"#hashtag", "@handle", "first.last+@example.com",
		"should", "—", "wait", "it", "break", "em-dashes", "quotes",
		"It'd", "it’ll",
	}

	for _, e := range expected {
		if !contains(e, got) {
			t.Errorf("Expected to find token %q, but did not.", e)
		}
	}

	// Check that last .
	last := got[len(got)-1]
	if last.String() != "." {
		t.Errorf("The last token should be %q, got %q.", ".", last.String())
	}

	// No trailing punctuation
	for _, token := range got {
		if utf8.RuneCountInString(token.String()) == 1 {
			// Skip actual (not trailing) punctuation
			continue
		}

		if strings.HasSuffix(token.String(), ",") || strings.HasSuffix(token.String(), ".") {
			t.Errorf("Found trailing punctuation in %q", token.String())
		}
	}
}

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

func contains(value string, tokens []Token) bool {
	for _, t := range tokens {
		if t.String() == value {
			return true
		}
	}
	return false
}
