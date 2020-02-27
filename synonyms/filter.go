package synonyms

import (
	"strings"
)

// Mapping is a tuple for mapping Synonyms to a Canonical form, thesaurus-like. The Synonyms property can be a comma separated string, indicating that each variation will map to the Canonical.
// For 'rules-based' synonyms, such as insensitivity to spaces or hyphens, pass an IgnoreFunc rather than trying to enumerate all the variations.
type Mapping struct {
	Synonyms, Canonical string
}

// IgnoreFunc is a function type specifying 'what to ignore' when looking up synonyms.
type IgnoreFunc func(string) string

// IgnoreCase instructs the sysnonyms filter to ignore (be insensitive to) case when looking up synonyms.
var IgnoreCase IgnoreFunc = strings.ToLower

// IgnoreSpace instructs the synonyms filter to ignore (be insensitive to) spaces when looking up synonyms. Not all whitespace, mind you, but precisely ascii 32.
var IgnoreSpace IgnoreFunc = Ignore(' ')

// Ignore instructs a synonyms filter to ignore (be insensitive to) specific characters when looking up synonyms.
func Ignore(r ...rune) IgnoreFunc {
	return func(s string) string {
		return strings.ReplaceAll(s, string(r), "")
	}
}

type Filter struct {
	ignoreFuncs   []IgnoreFunc
	lookup        map[string]string
	maxGramLength int
}

// New creates a new synonyms filter based on a set of Mappings and IgnoreFuncs. The latter are used to specify insensitivity to case or spaces, for example.
func New(mappings []Mapping, ignoreFuncs ...IgnoreFunc) *Filter {
	lookup := make(map[string]string)
	maxGramLength := 1

	for _, m := range mappings {
		synonyms := strings.Split(m.Synonyms, ",")
		len := len(synonyms)
		if len > maxGramLength {
			maxGramLength = len
		}
		for _, s := range synonyms {
			key := strings.TrimSpace(s)
			key = normalize(ignoreFuncs, key)
			if key != "" {
				lookup[key] = m.Canonical
			}
		}
	}

	return &Filter{
		ignoreFuncs:   ignoreFuncs,
		lookup:        lookup,
		maxGramLength: maxGramLength,
	}
}

func normalize(ignoreFuncs []IgnoreFunc, ss ...string) string {
	var b strings.Builder

	for _, s := range ss {
		for _, fn := range ignoreFuncs {
			s = fn(s)
		}
		b.WriteString(s)
	}

	return b.String()
}

func (f *Filter) Lookup(ss ...string) (string, bool) {
	key := normalize(f.ignoreFuncs, ss...)
	if key == "" {
		return "", false
	}

	canonical, found := f.lookup[key]

	return canonical, found
}

func (f *Filter) MaxGramLength() int {
	return f.maxGramLength
}
