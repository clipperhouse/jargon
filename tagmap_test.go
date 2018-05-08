package techlemm

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := map[string]string{
		"foo.js":      "foojs",
		".net":        ".net",
		"asp.net-mvc": "aspnetmvc",
	}

	for given, expected := range tests {
		got := normalize(given)
		if got != expected {
			t.Errorf("Given %q, expected %q, but got %q", given, expected, got)
		}
	}
}

func TestNewTagMap(t *testing.T) {
	tags := []string{"Node.js", "ASP.net"}
	tagmap := NewTagMap(tags)

	for _, value := range tags {
		key := normalize(value)
		got, exists := tagmap.values[key]
		if !exists {
			t.Errorf("Given added tag %q, expected exists true, but got %t", value, exists)
		}
		if got != value {
			t.Errorf("Given added tag %q, expected to retrieve same, but got %q", value, got)
		}
	}
}
