package jargon

import (
	"testing"
	"unicode"
)

func TestPunct(t *testing.T) {
	// ensure that punct exceptions are actually punct
	for r := range punctExceptions {
		if !unicode.IsPunct(r) {
			t.Errorf("%q is included in punctExceptions but it's not defined as unicode.IsPunct, and therefore is redundant", r)
		}
	}
	for r := range leadingPunct {
		if !unicode.IsPunct(r) {
			t.Errorf("%q is included in leadingPunct, but it's not defined as unicode.IsPunct, and therefore is redundant", r)
		}
	}
	for r := range midPunct {
		if !unicode.IsPunct(r) {
			t.Errorf("%q is included in midPunct, but it's not defined as unicode.IsPunct, and therefore is redundant", r)
		}
	}
}
