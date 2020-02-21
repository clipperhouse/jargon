package stopwords

import (
	"testing"
)

func TestCaseSensitive(t *testing.T) {
	stopwords := []string{
		"This",
		"and",
		"the",
	}
	filter := NewFilter(stopwords, true)

	type test struct {
		input   string
		output  string
		stopped bool
	}
	tests := []test{
		{"This", "", true},
		{"this", "this", false},
		{"The", "The", false},
		{"the", "", true},
		{"Foo", "Foo", false},
		{"foo", "foo", false},
	}

	for _, test := range tests {
		output, stopped := filter.Lookup([]string{test.input})
		if output != test.output {
			t.Errorf("output should have been %q, got %q", test.output, output)
		}
		if stopped != test.stopped {
			t.Errorf("stopped should have been %t, got %t", test.stopped, stopped)
		}
	}
}

func TestCaseInsensitive(t *testing.T) {
	stopwords := []string{
		"This",
		"and",
		"the",
	}
	filter := NewFilter(stopwords, false)

	type test struct {
		input   string
		output  string
		stopped bool
	}
	tests := []test{
		{"This", "", true},
		{"this", "", true},
		{"The", "", true},
		{"the", "", true},
		{"Foo", "Foo", false},
		{"foo", "foo", false},
	}

	for _, test := range tests {
		output, stopped := filter.Lookup([]string{test.input})
		if output != test.output {
			t.Errorf("output should have been %q, got %q", test.output, output)
		}
		if stopped != test.stopped {
			t.Errorf("stopped should have been %t, got %t", test.stopped, stopped)
		}
	}
}
