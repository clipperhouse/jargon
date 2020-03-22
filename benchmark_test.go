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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		_, err := jargon.Tokenize(r).Count()
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTokenizeLegacy(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		_, err := jargon.TokenizeLegacy(r).Count()
		if err != nil {
			b.Error(err)
		}
	}
}
