package jargon

import "unicode"

func isPunct(r rune) bool {
	return unicode.IsPunct(r) || spaceIsPunct(r)
}

func spaceIsPunct(r rune) bool {
	switch r {
	case
		'\n',
		'\r',
		'\t':
		return true
	}
	return false
}
