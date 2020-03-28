package jargon

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestPunct(t *testing.T) {
	// Likely obviated by new tokenizer

	// ensure that punct exceptions are actually punct
	// for r := range punctAsSymbol {
	// 	if !unicode.IsPunct(r) {
	// 		t.Errorf("%q is included in punctExceptions but it's not defined as unicode.IsPunct, and therefore is redundant", r)
	// 	}
	// }
	// for r := range leadingPunct {
	// 	if !unicode.IsPunct(r) {
	// 		t.Errorf("%q is included in leadingPunct, but it's not defined as unicode.IsPunct, and therefore is redundant", r)
	// 	}
	// }
	// for r := range midPunct {
	// 	if !unicode.IsPunct(r) {
	// 		t.Errorf("%q is included in midPunct, but it's not defined as unicode.IsPunct, and therefore is redundant", r)
	// 	}
	// }
}

func BenchmarkPunctSwitch(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	dummy := false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br := bytes.NewReader(file)
		r, _, err := br.ReadRune()
		if err != nil {
			b.Error(err)
		}

		dummy = spaceIsPunct(r)
		dummy = punctIsSymbol(r)
		dummy = isLeadingPunct(r)
		dummy = isMidPunct(r)
	}

	b.Log(dummy)
}
