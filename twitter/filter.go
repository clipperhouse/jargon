// Package twitter provides filters to identify Twitter-style @handles and #hashtags, and coalesce them into single tokens
package twitter

import (
	"unicode"
	"unicode/utf8"

	"github.com/clipperhouse/jargon/sigil"
)

// Handles will identify Twitter-style handles, combining the @ and name into a single token
var Handles = sigil.NewFilter("@", legalHandle)

// Hashtags will identify Twitter-style hashtags, combining the # and tag into a single token
var Hashtags = sigil.NewFilter("#", legalHashtag)

// https://help.twitter.com/en/managing-your-account/twitter-username-rules
func legalHandle(s string) bool {
	length := 0
	for _, r := range s {
		length++

		switch {
		case length > 15:
			return false
		case
			'A' <= r && r <= 'Z',
			'a' <= r && r <= 'r',
			'0' <= r && r <= '9',
			r == '_':
			continue
		default:
			return false
		}
	}

	return true
}

// Determined by playing with Twitter's web UI to see what got highlighted ;)
func legalHashtag(s string) bool {
	// One-character hashtags need to be a letter
	length := utf8.RuneCountInString(s)
	if length == 1 {
		r, _ := utf8.DecodeRuneInString(s)
		return unicode.IsLetter(r)
	}

	for _, r := range s {
		switch {
		case
			unicode.IsLetter(r),
			unicode.IsNumber(r),
			r == '_':
			continue
		default:
			return false
		}
	}

	return true
}
