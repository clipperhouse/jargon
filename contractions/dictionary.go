// Package contractions provides a jargon.Dictionary to expand English contractions, such as "don't" → "does not"
package contractions

type dictionary struct {
	variations map[string]string
}

// Lookup attempts to convert single-token contractions to non-contracted version. Examples:
// don't → does not
// We’ve → We have
// SHE'S -> SHE IS
func (d *dictionary) Lookup(s []string) (string, bool) {
	if len(s) != 1 {
		return "", false
	}

	canonical, ok := variations[s[0]]
	return canonical, ok
}

func (d *dictionary) MaxGramLength() int {
	return 1
}
