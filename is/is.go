// Package is provide utilities for identifying Unicode categories of runes, relating to Unicode text segmentation:
// https://unicode.org/reports/tr29/
package is

import "unicode"

// https://unicode.org/reports/tr29/#ALetter
func ALetter(r rune) bool {
	if r == '_' {
		return true
	}

	if unicode.IsLetter(r) {
		if unicode.Is(unicode.Hebrew, r) {
			return false
		}
		if unicode.Is(unicode.Katakana, r) {
			return false
		}
		if unicode.Is(unicode.Ideographic, r) {
			return false
		}
		return true
	}

	switch {
	case
		0x02C2 <= r && r <= 0x02C5,
		0x02D2 <= r && r <= 0x02D7,
		r == 0x02DE,
		r == 0x02DF,
		0x02E5 <= r && r <= 0x02EB,
		r == 0x02ED,
		0x02EF <= r && r <= 0x02FF,
		r == 0x055A,
		r == 0x055B,
		r == 0x055C,
		r == 0x055E,
		r == 0x058A,
		r == 0x05F3,
		0xA708 <= r && r <= 0xA716,
		r == 0xA720,
		r == 0xA721,
		r == 0xA789,
		r == 0xA78A,
		r == 0xAB5B:
		return true
	}

	return false
}

// AHLetter is any unicode letter or number, or underscore, and not ideographic
// Working to comply with https://unicode.org/reports/tr29/#WB5 through WB13
// Current logic implements WB5, WB8, WB9, WB10, WB13
func AHLetter(r rune) bool {
	return ALetter(r) || unicode.Is(unicode.Hebrew, r)
}

// MidLetter are runes allowed mid-word
// https://unicode.org/reports/tr29/#MidLetter
func MidLetter(r rune) bool {
	switch r {
	case
		//':',	// TODO Swedish
		'·',
		'·',
		'՟',
		'״',
		'‧',
		'︓',
		'﹕',
		'：':
		return true
	}
	return false
}

// MidNumLet are non-alphanumeric characters allowed within words
// See https://unicode.org/reports/tr29/#MidNumLet
func MidNumLet(r rune) bool {
	switch r {
	case
		'.',
		'’',
		'․',
		'﹒',
		'＇',
		'．':
		return true
	}
	return false
}

// MidNumLetQ are non-alphanumeric characters allowed within words
// See https://unicode.org/reports/tr29/#MidNumLet
func MidNumLetQ(r rune) bool {
	return MidNumLet(r) || r == '\''
}

// https://unicode.org/reports/tr14/
func InfixNumeric(r rune) bool {
	switch r {
	case
		0x002C,
		0x002E,
		0x003A,
		0x003B,
		0x037E,
		0x0589,
		0x060C,
		0x060D,
		0x07F8,
		0x2044,
		0xFE10,
		0xFE13,
		0xFE14:
		return true
	}

	return false
}

// https://unicode.org/reports/tr29/#MidNum
func MidNum(r rune) bool {
	switch r {
	case
		0x003A,
		0xFE13,
		0x002E:
		return false
	case
		0x066C,
		0xFE50,
		0xFE54,
		0xFF0C,
		0xFF1B:
		return true
	default:
		return InfixNumeric(r)
	}
}

func Numeric(r rune) bool {
	switch {
	case r == '٬':
		return false
	case 0xFF10 <= r && r <= 0xFF19:
		return true
	default:
		return unicode.IsNumber(r)
	}
}

func Cr(r rune) bool {
	return r == '\r'
}

// https://unicode.org/reports/tr29/#Katakana
func Katakana(r rune) bool {
	switch r {
	case
		0x3031,
		0x3032,
		0x3033,
		0x3034,
		0x3035,
		0x309B,
		0x309C,
		0x30A0,
		0x30FC,
		0xFF70:
		return true
	default:
		return unicode.Is(unicode.Katakana, r)
	}
}

func HebrewLetter(r rune) bool {
	return unicode.Is(unicode.Hebrew, r) && unicode.IsLetter(r)
}

func Leading(r rune) bool {
	switch {
	case r == '.':
		return true
	default:
		return false
	}
}
