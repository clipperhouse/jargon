// Package contractions provides a jargon.TokenFilter to expand English contractions, such as "don't" → "does not"
package contractions

// Expander expands common contractions into distinct words. Examples:
// don't → does not
// We’ve → We have
// SHE'S -> SHE IS
var Expander = &filter{}

type filter struct{}

// Lookup attempts to convert single-token contractions to non-contracted version. Examples:
// don't → does not
// We’ve → We have
// SHE'S -> SHE IS
func (f *filter) Lookup(s []string) (string, bool) {
	if len(s) != 1 {
		return "", false
	}

	canonical, ok := variations[s[0]]
	return canonical, ok
}

func (f *filter) MaxGramLength() int {
	return 1
}
