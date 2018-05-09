package techlemm

import (
	"strings"
)

// TagMap is the main structure for looking up canonical tags
type TagMap struct {
	values map[string]string
}

// NewTagMap creates and populates a new TagMap for the purpose of looking up canonical tags
func NewTagMap(d *Dictionary) *TagMap {
	result := &TagMap{
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

type Dictionary struct {
	Tags     []string
	Synonyms map[string]string
}

func NewDictionary(tags []string, synonyms map[string]string) *Dictionary {
	return &Dictionary{tags, synonyms}
}

// Get attempts to canonicalize a given input.
// Returned string is the canonical, if found; returned bool indicates whether found
func (t *TagMap) Get(s string) (string, bool) {
	key := normalize(s)
	canonical, found := t.values[key]
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
