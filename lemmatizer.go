package techlemm

import (
	"strings"
)

// Lemmatizer is the main structure for looking up canonical tags
type Lemmatizer struct {
	values map[string]string
}

// NewLemmatizer creates and populates a new Lemmatizer for the purpose of looking up canonical tags
func NewLemmatizer(d *Dictionary) *Lemmatizer {
	result := &Lemmatizer{
		values: make(map[string]string),
	}
	for _, tag := range d.Tags {
		key := normalize(tag)
		result.values[key] = tag
	}
	for synonym, canonical := range d.Synonyms {
		key := normalize(synonym)
		result.values[key] = canonical
	}
	return result
}

// Get attempts to canonicalize a given input.
// Returned string is the canonical, if found; returned bool indicates whether found
func (lem *Lemmatizer) Get(s string) (string, bool) {
	key := normalize(s)
	canonical, found := lem.values[key]
	return canonical, found
}

// normalize returns a string suitable as a key for tag lookup, removing dots and dashes and converting to lowercase
func normalize(s string) string {
	result := make([]rune, 0)

	for index, value := range s {
		if index == 0 {
			// Leading dots are meaningful and should not be removed, for example ".net"
			result = append(result, value)
			continue
		}
		if value == '.' || value == '-' {
			continue
		}
		result = append(result, value)
	}
	return strings.ToLower(string(result))
}
