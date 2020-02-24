package jargon_test

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestTokenize2(t *testing.T) {
	text := `Hi. This wishy-washy is .net and -123 and 12.34 — F# and C++, and TCP/IP and 
#hashtag and @handle and me_you-us+@email.com and http://foo.com/thing-stuff.`

	r := strings.NewReader(text)
	tokens := jargon.Tokenize2(r)

	for {
		token, err := tokens.Next()
		if err != nil {
			t.Error(err)
		}
		if token == nil {
			break
		}
		//		t.Log(token)
	}
}
