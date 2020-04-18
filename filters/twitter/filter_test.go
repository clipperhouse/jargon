package twitter

import (
	"fmt"
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestFilter(t *testing.T) {
	test := "This is a @handle."
	tokens := jargon.TokenizeString(test)
	//	tokens = Hashtags(tokens)
	tokens = tokens.Filter(Handles)

	got, err := tokens.ToSlice()
	fmt.Printf("%q", got)
	if err != nil {
		t.Error(err)
	}
	t.Log(got)
}
