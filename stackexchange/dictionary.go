package stackexchange

import "strings"

// dictionary satisfies the jargon.Dictionary interface
// Used in generated.go
type dictionary struct {
	tags     map[string]string
	synonyms map[string]string
}

func (d *dictionary) Lookup(s string) (string, bool) {
	key := normalize(s)
	canonical1, found1 := d.tags[key]

	if found1 {
		return canonical1, found1
	}

	canonical2, found2 := d.synonyms[key]
	return canonical2, found2
}

func normalize(s string) string {
	result := make([]rune, 0)

	for index, value := range s {
		if index == 0 {
			// Leading dots are meaningful and should not be removed, for example ".net"
			result = append(result, value)
			continue
		}
		if value == '.' || value == '-' || value == '/' {
			continue
		}
		result = append(result, value)
	}
	return strings.ToLower(string(result))
}
