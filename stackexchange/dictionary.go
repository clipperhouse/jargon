package stackexchange

import (
	"bytes"
	"strings"
)

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
	needsRewrite := false
loop:
	for i, r := range s {
		if i > 0 {
			switch r {
			case '.', '-', '/':
				needsRewrite = true
				break loop
			}
		}
	}

	if needsRewrite {
		buf := new(bytes.Buffer)
		for i, r := range s {
			if i > 0 {
				// Leading dots are meaningful and should not be removed, for example ".net"
				switch r {
				case '.', '-', '/':
					continue
				}
			}
			buf.WriteRune(r)
		}
		s = buf.String()
	}
	return strings.ToLower(s)
}
