package main

import (
	"testing"
)

func BenchmarkFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lemFile("~/Downloads/dict/data.verb")
	}
}

func TestFlush(t *testing.T) {
	s := "Hi objective c and pythOn"
	err := lemString(s)
	if err != nil {
		t.Error(err)
	}
	if w.Buffered() > 0 {
		t.Errorf("There are %d bytes left in the write buffer; should be zero (should have flushed)", w.Buffered())
	}
}
