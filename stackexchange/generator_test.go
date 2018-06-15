package stackexchange

import (
	"testing"
)

// Run this test to do the codegen: go test -run ^TestWriteDictionary$
func TestWriteDictionary(t *testing.T) {
	err := writeDictionary()

	if err != nil {
		t.Error(err)
	}
}

func TestTrailingVersion(t *testing.T) {
	tests := map[string]string{
		"ruby-on-3-rails-4": "ruby-on-3-rails",
		"python-2.7":        "python",
		"html5":             "html5", // considered part of the name, not a trailing version per se
	}
	for given, expected := range tests {
		got := trailingVersion.ReplaceAllString(given, "")
		if got != expected {
			t.Errorf("Given %q, expected trim trailing version number to be %q, got %q", given, expected, got)
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := map[string]string{
		"foo.js":      "foojs",
		".net":        ".net",
		"asp.net-mvc": "aspnetmvc",
		"os/2":        "os2",
	}

	for given, expected := range tests {
		got := normalize(given)
		if got != expected {
			t.Errorf("Given %q, expected %q, but got %q", given, expected, got)
		}
	}
}
