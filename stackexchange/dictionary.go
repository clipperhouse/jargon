package stackexchange

import (
	"bytes"
	"strings"
	"sync"
)

// dictionary satisfies the jargon.Dictionary interface
// Used in generated.go
type dictionary struct {
	tags     map[string]string
	synonyms map[string]string
}

func (d *dictionary) Lookup(s []string) (string, bool) {
	gram := strings.Join(s, "")
	key := normalize(gram)
	canonical1, found1 := d.tags[key]

	if found1 {
		return canonical1, found1
	}

	canonical2, found2 := d.synonyms[key]
	return canonical2, found2
}

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
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
		b := bufPool.Get().(*bytes.Buffer)
		b.Reset()
		for i, r := range s {
			if i > 0 {
				// Leading dots are meaningful and should not be removed, for example ".net"
				switch r {
				case '.', '-', '/':
					continue
				}
			}
			b.WriteRune(r)
		}
		s = b.String()
		bufPool.Put(b)
	}
	return strings.ToLower(s)
}
