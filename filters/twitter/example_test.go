package twitter_test

import (
	"fmt"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/twitter"
)

func ExampleHandles() {
	text := "Here's a @username."

	before, err := jargon.TokenizeString(text).ToSlice()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Before: %q, ", before)

	after, err := jargon.TokenizeString(text).Filter(twitter.Handles).ToSlice()
	if err != nil {
		panic(err)
	}
	fmt.Printf("after: %q", after)

	// Output: Before: ["Here's" " " "a" " " "@" "username" "."], after: ["Here's" " " "a" " " "@username" "."]
}
