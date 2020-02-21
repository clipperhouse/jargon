package stackexchange

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := map[string]string{
		"foo.js":      "foojs",
		".Net":        ".net",
		"ASP.net-mvc": "aspnetmvc",
		"os/2":        "os2",
	}

	for given, expected := range tests {
		got := normalize(given)
		if got != expected {
			t.Errorf("Given %q, expected %q, but got %q", given, expected, got)
		}
	}
}

func BenchmarkNormalize(b *testing.B) {
	wikipedia, err := ioutil.ReadFile("../testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	words := strings.Fields(string(wikipedia)) // good enough for this test

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range words {
			normalize(s)
		}
	}
}
