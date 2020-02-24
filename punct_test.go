package jargon

import (
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
