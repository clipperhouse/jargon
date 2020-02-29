// Package synonyms enables mapping synonyms to canonical terms
package synonyms

import (
	"fmt"
	"strings"

	"github.com/clipperhouse/jargon"
)

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
func NewFilter(mappings map[string]string, ignoreFuncs ...IgnoreFunc) (*Filter, error) {
	lookup := make(map[string]string)
	var maxGramLength int = 1

	for synonyms, canonical := range mappings {
		for _, synonym := range strings.Split(synonyms, ",") {

			// Need to count word tokens, not naively split on space.
			tokens := jargon.Tokenize(strings.NewReader(synonym))
			count, err := tokens.Words().Count()
			if err != nil {
				return nil, err
			}
			if count > maxGramLength {
				maxGramLength = count
			}

			key := strings.TrimSpace(synonym)
			key = normalize(ignoreFuncs, key)
			if key == "" {
				err := fmt.Errorf("the synonym %q, from the {%q: %q} mapping, results in an empty string when normalized", synonym, synonyms, canonical)
				return nil, err
			}

			// The same key should not point to multiple different synonyms
			existing, found := lookup[key]
			if found && existing != canonical {
				err := fmt.Errorf("the synonym %q (normalized to %q) from the {%q: %q} mapping, would overwrite an earlier mapping to %q. choose one or the other", synonym, key, synonyms, canonical, existing)
				return nil, err
			}

			lookup[key] = canonical
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
