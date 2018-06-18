package main

import (
	"testing"
)

func BenchmarkFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lemFile("~/Downloads/dict/data.verb")
	}
}
