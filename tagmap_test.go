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

func TestAddTag(t *testing.T) {
	tagmap := NewTagMap()
	tag := "Node.js"
	tagmap.AddTag(tag)

	key := normalize(tag)
	got, exists := tagmap.values[key]
	if !exists {
		t.Errorf("Given added tag %q, expected exists true, but got %t", tag, exists)
	}
	if got != tag {
		t.Errorf("Given added tag %q, expected to retrieve same, but got %q", tag, got)
	}
}
