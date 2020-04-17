package twitter

import (
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestFilter(t *testing.T) {
	test := "This is a @handle and a #hashtag"
	tokens := jargon.TokenizeString(test)
	tokens = Hashtags(tokens)
	tokens = Handles(tokens)

	got, err := tokens.ToSlice()
	if err != nil {
		t.Error(err)
	}
	t.Log(got)
}
