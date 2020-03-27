package stackoverflow

//go:generate go run generate/main.go

// Tags is the main exported Tags of Stack Overflow tags and synonyms. It's indended to identify canonical tags (technologies), even in prose.
// For example, the phrase "Ruby on Rails" (3 words) will be replaced with ruby-on-rails (1 word).
// It is insensitive to spaces, hyphens, dots and forward slashes, so "react js" and "reactjs" and "react.js" are all identified as the same canonical term.
var Tags = filter
