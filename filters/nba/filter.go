package nba

import "github.com/clipperhouse/jargon/filters/synonyms"

//go:generate go run generate/main.go

var ignore = []rune{' ', '.', '\'', '-', '–'}

// CurrentPlayers is a token filter for identifying current NBA players accoring to Wikipedia
// It is insensitive to spaces, dashes, apostrophes, periods and diacritics in players' names.
var CurrentPlayers = synonyms.NewFilter(mappings, true, ignore)
