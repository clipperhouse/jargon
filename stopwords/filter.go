package stopwords

import (
	"strings"
)

// NewFilter creates a token filter for the supplied stop words
func NewFilter(stopwords []string, caseSensitive bool) *filter {
	includes := make(map[string]bool)
	for _, s := range stopwords {
		var key string
		if caseSensitive {
			key = s
		} else {
			key = strings.ToLower(s)
		}
		includes[key] = true
	}

	return &filter{
		includes:      includes,
		caseSensitive: caseSensitive,
	}
}

type filter struct {
	includes      map[string]bool
	caseSensitive bool
}

func (f *filter) Lookup(s ...string) (string, bool) {
	if len(s) < 1 {
		return "", false
	}

	word := s[0] // max gram length is 1

	var key string
	if f.caseSensitive {
		key = word
	} else {
		key = strings.ToLower(word)
	}

	if f.includes[key] {
		return "", true
	}
	return word, false
}

func (f *filter) MaxGramLength() int {
	return 1
}
