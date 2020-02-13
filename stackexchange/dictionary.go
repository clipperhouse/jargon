package stackexchange

import (
	"bytes"
	"strings"
	"sync"
)

// Dictionary is the main exported Dictionary of Stack Exchange tags and synonyms, from the following Stack Exchange sites: Stack Overflow,
// Server Fault, Game Dev and Data Science. It's indended to identify canonical tags (technologies),
// e.g. Ruby on Rails (3 words) will be replaced with ruby-on-rails (1 word).
// It includes the most popular 2530 tags and 2022 synonyms
var Dictionary = &dictionary{}

// dictionary satisfies the jargon.Dictionary interface
// Used in generated.go
type dictionary struct{}

func (d *dictionary) Lookup(s []string) (string, bool) {
	gram := strings.Join(s, "")
	key := normalize(gram)
	canonical1, found1 := tags[key]

	if found1 {
		return canonical1, found1
	}

	canonical2, found2 := synonyms[key]
	return canonical2, found2
}

func (d *dictionary) MaxGramLength() int {
	return 3
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
