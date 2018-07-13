package jargon

import "unicode"

func isPunct(r rune) bool {
	return unicode.IsPunct(r) && !isPunctException(r)
}

var ok = struct{}{} // like a bool for maps, but with no allocation
var punctExceptions = map[rune]struct{}{
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

func isPunctException(r rune) bool {
	_, ok := punctExceptions[r]
	return ok
}

var leadingPunct = map[rune]struct{}{
	// Punctuation that can lead a word, like .Net
	'.': ok,
}

func mightBeLeadingPunct(r rune) bool {
	_, ok := leadingPunct[r]
	return ok
}

var midPunct = map[rune]struct{}{
	// Punctuation that can appear mid-word
	'.':  ok,
	'\'': ok,
	'’':  ok,
	':':  ok,
	'?':  ok,
	'&':  ok,
}

func mightBeMidPunct(r rune) bool {
	_, ok := midPunct[r]
	return ok
}
