// Package numbers provides a jargon.Dictionary to lemmatize numbers expressed as words, such as "three hundred" => "300"
package numbers

import (
	"strconv"
	"strings"
)

// The form is exactly one leading number followed by zero or more magnitudes
// In the above example, "three" is the number, "hundred" and "thousand" are magnitudes
// All passed tokens must contribute to the number for the lookup to succeed

type dictionary struct{}

var Dictionary = &dictionary{}

// Lookup attempts to turn a slice of token strings into a canonical number string.
// If successful, Lookup will return the number, and a bool indicating success or failure. Examples:
//	["three"] => "3", true
//	["three","thousand"] => "3000", true
//	["thirty-five","thousand"] => "35000", true
//	["three","hundred","thousand"] => "300000", true
//	["3","thousand"] => "3000", true
//	["-3","thousand"] => "3000", true
//	["+3","thousand"] => "3000", true
//	["2.54","million"] => "2540000", true
//	["1,000,000"] => "1000000", true
//	["hundred"] => "", false
//	["foo","3","hundred"] => "", false
// All tokens need to contribute to the number. Any token that is not a number results in Lookup returning false.
// Lookup works on short 'multiplicative' number phrases, where each number is multiplied into a total value, as in the examples above.
// Lookup does not handle 'compound' or 'additive' phrases like "one thousand five hundred twenty".
// In the above example, 'thirty-five' only works when it is a single token, no spaces, hyphen optional.
// Commas are ignored, so long as there are no spaces.
func (d *dictionary) Lookup(s []string) (string, bool) {
	if len(s) == 0 {
		return "", false
	}
	p := parser{tokens: s}
	return p.Parse()
}

type parser struct {
	tokens []string
	ints   []int64
	floats []float64
	pos    int
}

func (p *parser) current() string {
	return normalize(p.tokens[p.pos])
}

func (p *parser) Parse() (string, bool) {
	p.pos = 0
	first := p.current()

	// Try parsing first token as int
	i, err := strconv.ParseInt(first, 10, 64)
	if err == nil {
		switch {
		case hasLeadingZero(first):
			// Special case for leading zeros on an integer: assume it is intentional,
			// such as a zip code, serial number, phone number
			return "", false
		default:
			p.ints = append(p.ints, i)
			p.pos++
			return p.parseMagnitudesInt()
		}
	}

	// Try parsing first token as float
	f, err := strconv.ParseFloat(first, 64)
	if err == nil {
		p.floats = append(p.floats, f)
		p.pos++
		return p.parseMagnitudesFloat()
	}

	// Try parsing first token as number word (int)
	num, ok := numbers[first]
	if ok {
		p.ints = append(p.ints, num)
		p.pos++
		return p.parseMagnitudesInt()
	}

	return "", false
}

func normalize(s string) string {
	result := s
	result = strings.Replace(result, ",", "", -1)
	result = strings.Replace(result, "-", "", -1)
	if s[0] == '-' { // leading hyphen is ok, bring it back
		return "-" + result
	}
	return result
}

func hasLeadingZero(s string) bool {
	if s[0] == '0' {
		return true
	}

	if len(s) > 1 {
		switch s[:2] {
		case "+0", "-0":
			return true
		}
	}

	return false
}

func (p *parser) parseMagnitudesInt() (string, bool) {
	for p.pos < len(p.tokens) {
		m, ok := magnitudes[p.current()]
		if !ok {
			return "", false
		}
		p.ints = append(p.ints, m)
		p.pos++
	}

	result := int64(1) // identity
	for i := 0; i < len(p.ints); i++ {
		result = result * p.ints[i]
	}

	return strconv.FormatInt(result, 10), true
}

func (p *parser) parseMagnitudesFloat() (string, bool) {
	for p.pos < len(p.tokens) {
		m, ok := magnitudes[p.current()]
		if !ok {
			return "", false
		}
		p.floats = append(p.floats, float64(m))
		p.pos++
	}

	result := float64(1) // identity
	for i := 0; i < len(p.floats); i++ {
		result = result * p.floats[i]
	}

	return strconv.FormatFloat(result, 'f', -1, 64), true
}

var numbers = map[string]int64{
	"one":          1,
	"two":          2,
	"three":        3,
	"four":         4,
	"five":         5,
	"six":          6,
	"seven":        7,
	"eight":        8,
	"nine":         9,
	"ten":          10,
	"eleven":       11,
	"twelve":       12,
	"thirteen":     13,
	"fourteen":     14,
	"fifteen":      15,
	"sixteen":      16,
	"seventeen":    17,
	"eighteen":     18,
	"nineteen":     19,
	"twenty":       20,
	"twentyone":    21,
	"twentytwo":    22,
	"twentythree":  23,
	"twentyfour":   24,
	"twentyfive":   25,
	"twentysix":    26,
	"twentyseven":  27,
	"twentyeight":  28,
	"twentynine":   29,
	"thirty":       30,
	"thirtyone":    31,
	"thirtytwo":    32,
	"thirtythree":  33,
	"thirtyfour":   34,
	"thirtyfive":   35,
	"thirtysix":    36,
	"thirtyseven":  37,
	"thirtyeight":  38,
	"thirtynine":   39,
	"forty":        40,
	"fortyone":     41,
	"fortytwo":     42,
	"fortythree":   43,
	"fortyfour":    44,
	"fortyfive":    45,
	"fortysix":     46,
	"fortyseven":   47,
	"fortyeight":   48,
	"fortynine":    49,
	"fifty":        50,
	"fiftyone":     51,
	"fiftytwo":     52,
	"fiftythree":   53,
	"fiftyfour":    54,
	"fiftyfive":    55,
	"fiftysix":     56,
	"fiftyseven":   57,
	"fiftyeight":   58,
	"fiftynine":    59,
	"sixty":        60,
	"sixtyone":     61,
	"sixtytwo":     62,
	"sixtythree":   63,
	"sixtyfour":    64,
	"sixtyfive":    65,
	"sixtysix":     66,
	"sixtyseven":   67,
	"sixtyeight":   68,
	"sixtynine":    69,
	"seventy":      70,
	"seventyone":   71,
	"seventytwo":   72,
	"seventythree": 73,
	"seventyfour":  74,
	"seventyfive":  75,
	"seventysix":   76,
	"seventyseven": 77,
	"seventyeight": 78,
	"seventynine":  79,
	"eighty":       80,
	"eightyone":    81,
	"eightytwo":    82,
	"eightythree":  83,
	"eightyfour":   84,
	"eightyfive":   85,
	"eightysix":    86,
	"eightyseven":  87,
	"eightyeight":  88,
	"eightynine":   89,
	"ninety":       90,
	"ninetyone":    91,
	"ninetytwo":    92,
	"ninetythree":  93,
	"ninetyfour":   94,
	"ninetyfive":   95,
	"ninetysix":    96,
	"ninetyseven":  97,
	"ninetyeight":  98,
	"ninetynine":   99,
}

var magnitudes = map[string]int64{
	"hundred":     1e2,
	"thousand":    1e3,
	"million":     1e6,
	"billion":     1e9,
	"trillion":    1e12,
	"quadrillion": 1e15,
	"quintillion": 1e15,
}
