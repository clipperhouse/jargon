package jargon_test

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestTokenizeHTML(t *testing.T) {
	h := `<html>
<p foo="bar">
Hi! Let's talk Ruby on Rails.
<!-- Ignore ASPNET MVC in comments -->
<script src="foo">var Nodejs = Reactjs;</script>
<style>p { margin-bottom:20px; } </style>
</p>
</html>
`
	r := strings.NewReader(h)
	tokens, err := jargon.TokenizeHTML(r).ToSlice()
	if err != nil {
		t.Error(err)
	}

	got := map[string]bool{}
	for _, token := range tokens {
		got[token.String()] = true
	}

	expected := []string{
		// tags kept whole
		`<p foo="bar">`,
		"</p>",
		// whitespace preserved
		"\n",
		// text nodes got tokenized
		"Hi", "!",
		"Ruby", "on", "Rails",
		// comment kept whole
		"<!-- Ignore ASPNET MVC in comments -->",
		// contents of script not tokenized
		`<script src="foo">`,
		"var Nodejs = Reactjs;",
		"</script>",
		// contents of style not tokenized
		"<style>",
		"p { margin-bottom:20px; } ",
		"</style>",
	}

	for _, e := range expected {
		if !got[e] {
			t.Errorf("Expected to find token %q, but did not.", e)
		}
	}
}
