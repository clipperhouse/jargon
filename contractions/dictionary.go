package contractions

type dictionary struct {
	variations map[string]string
}

// Lookup attempts to convert single-token contractions to non-contracted version.  Examples:
//	don't → does not
//	we've → we have
//	she's -> she is
// Caveats:
// - only lower case right now
// - returns expanded words as a single token with a space in it; caller might wish to re-tokenize
func (d *dictionary) Lookup(s []string) (string, bool) {
	if len(s) != 1 {
		return "", false
	}

	canonical, ok := variations[s[0]]
	if !ok {
		return "", false
	}

	return canonical, true
}
