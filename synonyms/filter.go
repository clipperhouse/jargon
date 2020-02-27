package synonyms

import (
	"fmt"
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
func Ignore(runes ...rune) IgnoreFunc {
	return func(s string) string {
		for _, r := range runes {
			s = strings.ReplaceAll(s, string(r), "")
		}
		return s
	}
}

// Filter implements the jargon.Filter interface
type Filter struct {
	ignoreFuncs   []IgnoreFunc
	lookup        map[string]string
	maxGramLength int
}

// NewFilter creates a new synonyms filter based on a set of Mappings and IgnoreFuncs. The latter are used to specify insensitivity to case or spaces, for example.
func NewFilter(mappings []Mapping, ignoreFuncs ...IgnoreFunc) (*Filter, error) {
	lookup := make(map[string]string)
	maxGramLength := 1

	for _, m := range mappings {
		synonyms := strings.Split(m.Synonyms, ",")

		for _, synonym := range synonyms {
			grams := len(strings.Fields(synonym))
			if grams > maxGramLength {
				maxGramLength = grams
			}

			key := strings.TrimSpace(synonym)
			key = normalize(ignoreFuncs, key)
			if key == "" {
				err := fmt.Errorf("the synonym %q, from the {%q: %q} mapping, results in an empty string when normalized", synonym, m.Synonyms, m.Canonical)
				return nil, err
			}

			// The same key should not point to multiple different synonyms
			existing, found := lookup[key]
			if found && existing != m.Canonical {
				err := fmt.Errorf("the synonym %q (normalized to %q) from the {%q: %q} mapping, would overwrite an earlier mapping to %q", synonym, key, m.Synonyms, m.Canonical, existing)
				return nil, err
			}

			lookup[key] = m.Canonical
		}
	}

	return &Filter{
		ignoreFuncs:   ignoreFuncs,
		lookup:        lookup,
		maxGramLength: maxGramLength,
	}, nil
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
