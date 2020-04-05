package twitter

import (
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestFilter(t *testing.T) {
	test := "This is a @handle and a #hashtag"
	tokens := jargon.TokenizeString(test)
	_, err := tokens.Filter(Filter).ToSlice()
	if err != nil {
		t.Error(err)
	}
}
