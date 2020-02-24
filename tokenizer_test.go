package jargon_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/clipperhouse/jargon"
)

func TestTokenize(t *testing.T) {
	text := `Hi! This is a test of tech terms.
It should consider F#, C++, .net, Node.JS and 3.141592 to be their own tokens. 
Similarly, #hashtag and @handle should work, as should an first.last+@example.com.
It should—wait for it—break on things like em-dashes and "quotes" and it ends.
It'd be great it it’ll handle apostrophes.
`
	r := strings.NewReader(text)
	tokens := jargon.Tokenize2(r)
	got, err := tokens.ToSlice()
	if err != nil {
		t.Error(err)
	}

	expected := []string{
		"Hi", "!", "a",
		"F#", "C++", ".net", "Node.JS", "3.141592",
		"#hashtag", "@handle", "first.last+@example.com",
		"should", "—", "wait", "it", "break", "em-dashes", "\"", "quotes", "ends",
		"It'd", "it’ll", "apostrophes",
	}

	for _, e := range expected {
		if !contains(e, got) {
			t.Errorf("Expected to find token %q, but did not.", e)
		}
	}

	// Check that last .
	nextToLast := got[len(got)-2]
	if nextToLast.String() != "." {
		fmt.Println(got)
		t.Errorf("next-to-last token should be %q, got %q.", ".", nextToLast)
	}

	// Check that last \
	last := got[len(got)-1]
	if last.String() != "\n" {
		t.Errorf("The last token should be %q, got %q.", "\n", last)
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

func TestTokenizeHTML(t *testing.T) {
	h := `<html>
<p foo="bar">
Hi! Let's talk Ruby on Rails.
<!-- Ignore ASPNET MVC in comments -->
</p>
</html>
`
	r := strings.NewReader(h)
	got, err := jargon.TokenizeHTML(r).ToSlice()

	if err != nil {
		t.Error(err)
	}

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

func TestURLs(t *testing.T) {
	// We mostly get lucky on URLs due to punct rules

	tests := map[string]string{
		`http://www.google.com`:                     `http://www.google.com`,                     // as-is
		`http://www.google.com.`:                    `http://www.google.com`,                     // "."" should be considered trailing punct
		`http://www.google.com/`:                    `http://www.google.com/`,                    // trailing slash OK
		`http://www.google.com/?`:                   `http://www.google.com/`,                    // "?" should be considered trailing punct
		`http://www.google.com/?foo=bar`:            `http://www.google.com/?foo=bar`,            // "?" is querystring
		`http://www.google.com/?foo=bar.`:           `http://www.google.com/?foo=bar`,            // trailing "."
		`http://www.google.com/?foo=bar&qaz=qux`:    `http://www.google.com/?foo=bar&qaz=qux`,    // "?" with &
		`http://www.google.com/?foo=bar&qaz=q%20ux`: `http://www.google.com/?foo=bar&qaz=q%20ux`, // with encoding
		`//www.google.com`:                          `//www.google.com`,                          // scheme-relative
		`/usr/local/bin/foo.bar`:                    `/usr/local/bin/foo.bar`,
		`c:\windows\notepad.exe`:                    `c:\windows\notepad.exe`,
	}

	for input, expected := range tests {
		r := strings.NewReader(input)
		got, err := jargon.Tokenize(r).Next() // just take the first token
		if err != nil {
			t.Error(err)
		}

		if got.String() != expected {
			t.Errorf("Expected URL %s to result in %s, but got %s", input, expected, got)
		}
	}
}

func BenchmarkTokenize(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		consume(jargon.Tokenize(r))
	}
}
