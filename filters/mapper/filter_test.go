package mapper

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestNewFilter(t *testing.T) {
	fn := func(token *jargon.Token) *jargon.Token {
		s := strings.ToLower(token.String())
		if s != token.String() {
			return jargon.NewToken(s, true)
		}
		return token
	}

	filter := NewFilter(fn)

	text := "Here IS a nEw tesT."
	expected := "here is a new test."
	got, err := jargon.TokenizeString(text).Filter(filter).String()
	if err != nil {
		t.Error(err)
	}

	if expected != got {
		t.Errorf("given %q, expected %q, got %s", text, expected, got)
	}
}
