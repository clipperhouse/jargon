package stackexchange

import "testing"

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

func BenchmarkNormalize(b *testing.B) {
	strs := []string{
		"foo.js",
		".net",
		"asp.net-mvc",
		"os/2",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range strs {
			normalize(s)
		}
	}
}
