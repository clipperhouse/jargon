package jargon

import "unicode"

type runeSet map[rune]bool

func isPunct(r rune) bool {
	return (unicode.IsPunct(r) || spaceAsPunct[r]) && !punctAsSymbol[r]
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
}

var spaceAsPunct = runeSet{
	'\n': true,
	'\r': true,
	'\t': true,
}

var leadingPunct = runeSet{
	// Punctuation that can lead a word, like .Net
	'.': true,
	'-': true,
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
