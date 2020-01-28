package jargon

import "unicode"

type runeMap map[rune]struct{}

func (rm runeMap) includes(r rune) bool {
	_, ok := rm[r]
	return ok
}

func isPunct(r rune) bool {
	return (unicode.IsPunct(r) || spaceAsPunct.includes(r)) && !punctAsSymbol.includes(r)
}

var ok = struct{}{} // like a bool for maps, but with no allocation
var punctAsSymbol = runeMap{
	// In some cases, we want to consider a rune a symbol, even though Unicode defines it as punctuation
	// See http://www.unicode.org/faq/punctuation_symbols.html
	'-':  ok,
	'#':  ok,
	'@':  ok,
	'*':  ok,
	'%':  ok,
	'_':  ok,
	'/':  ok,
	'\\': ok,
}

var spaceAsPunct = runeMap{
	'\n': ok,
	'\r': ok,
	'\t': ok,
}

var leadingPunct = runeMap{
	// Punctuation that can lead a word, like .Net
	'.': ok,
}

var midPunct = runeMap{
	// Punctuation that can appear mid-word
	'.':  ok,
	'\'': ok,
	'’':  ok,
	':':  ok,
	'?':  ok,
	'&':  ok,
}
