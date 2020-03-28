package jargon

import "unicode"

type runeSet map[rune]bool

func isPunct(r rune) bool {
	return (unicode.IsPunct(r) || spaceIsPunct(r)) && !punctIsSymbol(r)
}

func punctIsSymbol(r rune) bool {
	switch r {
	case '-',
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

var punctAsSymbol = runeSet{
	// In some cases, we want to consider a rune a symbol, even though Unicode defines it as punctuation
	// See http://www.unicode.org/faq/punctuation_symbols.html
	'-':  true,
	'+':  true,
	'#':  true,
	'@':  true,
	'*':  true,
	'%':  true,
	'/':  true,
	'\\': true,
	':':  true,
}

func spaceIsPunct(r rune) bool {
	switch r {
	case '\n', '\r', '\t':
		return true
	}
	return false
}

var spaceAsPunct = runeSet{
	'\n': true,
	'\r': true,
	'\t': true,
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

var leadingPunct = runeSet{
	// Punctuation that can lead a word, like .Net
	'.': true,
	'-': true,
}

func isMidPunct(r rune) bool {
	switch r {
	case '.',
		'\'',
		'’',
		':',
		'?',
		'&':
		return true
	}
	return false
}

var midPunct = runeSet{
	// Punctuation that can appear mid-word
	'.':  true,
	'\'': true,
	'’':  true,
	':':  true,
	'?':  true,
	'&':  true,
}
