package jargon_test

import (
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
