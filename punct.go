package jargon

import "unicode"

func isPunct(r rune) bool {
	return (unicode.IsPunct(r) || spaceIsPunct(r)) && !punctIsSymbol(r)
}

func punctIsSymbol(r rune) bool {
	switch r {
	case
		'-',
		'+',
		'#',
		'@',
		'*',
		'%',
		'/',
		'\\',
		':':
		return true
	}
	return false
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

func isLeadingPunct(r rune) bool {
	switch r {
	case
		'.',
		'-':
		return true
	}
	return false
}

func isMidPunct(r rune) bool {
	switch r {
	case
		'.',
		'\'',
		'’',
		':',
		'?',
		'&':
		return true
	}
	return false
}
