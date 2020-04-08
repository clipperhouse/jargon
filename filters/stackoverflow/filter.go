// Package stackoverflow provides a filter for identifying technical terms in jargon
package stackoverflow

import (
	"github.com/clipperhouse/jargon/filters/synonyms"
)

//go:generate go run generate/main.go

// Tags detects Stack Overflow tags and synonyms. It's indended to identify canonical tags (technologies), even in prose.
// For example, the phrase "Ruby on Rails" (3 words) will be replaced with ruby-on-rails (1 word).
// It is insensitive to spaces, hyphens, dots and forward slashes, so "react js" and "reactjs" and "react.js" are all identified as the same canonical term.
var Tags = synonyms.NewFilter(mappings, true, ignoreRunes)

var ignoreRunes = []rune{' ', '-', '.', '/'}
