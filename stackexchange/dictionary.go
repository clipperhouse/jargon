package stackexchange

import "strings"

// dictionary satisfies the jargon.Dictionary interface
// Used in generated.go
type dictionary struct {
	tags     []string
	synonyms map[string]string
}

func (d *dictionary) GetTags() []string {
	return d.tags
}

func (d *dictionary) GetSynonyms() map[string]string {
	return d.synonyms
}

func (d *dictionary) MaxGramLength() int {
	return 3
}

// Normalize returns a string suitable as a key for tag lookup, removing dots, dashes and forward slashes, and converting to lowercase
func (d *dictionary) Normalize(s string) string {
	return normalize(s)
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
