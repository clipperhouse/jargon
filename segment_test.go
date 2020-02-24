package jargon_test

import (
	"strings"
	"testing"

	"github.com/blevesearch/segment"
)

func TestSegments(t *testing.T) {
	s := "This is .net -123 60% thing_stuff foo-bar e*trade 1+2 @handle @hashtag asp.net and tcp/ip."
	r := strings.NewReader(s)
	segmenter := segment.NewSegmenter(r)

	for segmenter.Segment() {
		//		fmt.Printf("%q\n", segmenter.Bytes())
	}
}
