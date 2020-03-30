package jargon_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/clipperhouse/jargon"
)

// TODO: test ordering

func TestTokenize2(t *testing.T) {
	text := `Hi. Let's see node.js, 123.456, 1,000. ウィキペディア 象形. It includes first_last, and 123.`
	tokens := jargon.TokenizeString2(text)

	for {
		token, err := tokens.Next()
		if err != nil {
			t.Error(err)
		}
		if token == nil {
			break
		}

		t.Log(token)
	}
}

func BenchmarkTokenize2(b *testing.B) {
	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	var count int
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		c, err := jargon.Tokenize2(r).Count()
		if err != nil {
			b.Error(err)
		}
		count = c
	}
	b.Logf("%d tokens\n", count)
}
