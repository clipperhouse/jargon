package norm

import (
	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/mapper"
	"golang.org/x/text/unicode/norm"
)

// NFC normalizes tokens into Unicode Normalization Form C
var NFC = newFilter(norm.NFC)

// NFD normalizes tokens into Unicode Normalization Form D
var NFD = newFilter(norm.NFD)

// NFKC normalizes tokens into Unicode Normalization Form KC
var NFKC = newFilter(norm.NFKC)

// NFKD normalizes tokens into Unicode Normalization Form KD
var NFKD = newFilter(norm.NFKD)

func newFilter(form norm.Form) jargon.Filter {
	f := func(token *jargon.Token) *jargon.Token {
		if form.IsNormalString(token.String()) {
			return token
		}

		s := form.String(token.String())
		return jargon.NewToken(s, true)
	}

	return mapper.NewFilter(f)
}
