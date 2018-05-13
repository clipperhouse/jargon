package tokenizers

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTechProse(t *testing.T) {
	text := `Hi! This is a test of tech terms.
It should consider F#, C++, .net, Node.JS and 3.141592 to be their own tokens. 
Similarly, #hashtag and @handle should work, as should an first.last+@example.com.
It should—wait for it—break on things like em-dashes and "quotes" and it ends.`
	got := TechProse.Tokenize(text)

	expected := []string{
		"Hi", "!",
		"F#", "C++", ".net", "Node.JS", "3.141592",
		"#hashtag", "@handle", "first.last+@example.com",
		"should", "—", "wait", "it", "break", "em-dashes", "quotes",
	}

	//	fmt.Println(strings.Join(got, "➡"))

	for _, e := range expected {
		if !contains(e, got) {
			t.Errorf("Expected to find token %q, but did not.", e)
		}
	}

	// Check that last .
	last := got[len(got)-1]
	if last != "." {
		t.Errorf("The last token should be %q, got %q.", ".", last)
	}

	// No trailing punctuation
	for _, token := range got {
		if utf8.RuneCountInString(token) == 1 {
			// Skip actual (not trailing) punctuation
			continue
		}

		if strings.HasSuffix(token, ",") || strings.HasSuffix(token, ".") {
			t.Errorf("Found trailing punctuation in %q", token)
		}
	}
}

func contains(token string, tokens []string) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}
	return false
}
