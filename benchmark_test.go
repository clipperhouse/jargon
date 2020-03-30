package jargon_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/clipperhouse/jargon"
)

func BenchmarkTokenizeUniseg(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		_, err := jargon.TokenizeUniseg(r).Count()
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTokenize(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	var count int
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		c, err := jargon.Tokenize(r).Count()
		if err != nil {
			b.Error(err)
		}
		count = c
	}
	b.Logf("%d tokens\n", count)
}

func BenchmarkTokenizeHTML(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		_, err := jargon.TokenizeHTML(r).Count()
		if err != nil {
			b.Error(err)
		}
	}
}
