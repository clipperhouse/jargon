package stopwords_test

import (
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/stopwords"
)

func TestCaseSensitive(t *testing.T) {
	words := []string{
		"One",
		"four",
		"five",
	}
	stop := stopwords.NewFilter(words, false)

	type test struct {
		input  string
		output string
	}
	tests := []test{
		{"One two three four five", " two three  "},
		{"one two three Four five", "one two three Four "},
	}

	for _, test := range tests {
		tokens := jargon.TokenizeString(test.input)
		output, err := stop(tokens).String()
		if err != nil {
			t.Error(err)
		}
		if output != test.output {
			t.Errorf("output should have been %q, got %q", test.output, output)
		}
	}
}

func TestCaseInsensitive(t *testing.T) {
	words := []string{
		"One",
		"four",
		"five",
	}
	stop := stopwords.NewFilter(words, true)

	type test struct {
		input  string
		output string
	}
	tests := []test{
		{"One two three four five", " two three  "},
		{"one two three Four five", " two three  "},
	}

	for _, test := range tests {
		tokens := jargon.TokenizeString(test.input)
		output, err := stop(tokens).String()
		if err != nil {
			t.Error(err)
		}
		if output != test.output {
			t.Errorf("output should have been %q, got %q", test.output, output)
		}
	}
}
