package jargon_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/clipperhouse/jargon"
)

func BenchmarkTokenize(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	var count int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		tokens, err := jargon.TokenizeLegacy(r).ToSlice()
		if err != nil {
			b.Error(err)
		}
		if i == 0 {
			count = len(tokens)
		}
	}
	b.Logf("token count: %d", count)
}

func BenchmarkTokenizeLegacy(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	var count int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		tokens, err := jargon.TokenizeLegacy(r).ToSlice()
		if err != nil {
			b.Error(err)
		}
		if i == 0 {
			count = len(tokens)
		}
	}
	b.Logf("token count: %d", count)
}
