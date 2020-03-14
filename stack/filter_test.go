package stack

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestFill(t *testing.T) {
	type test struct {
		input string
		count int
	}

	expecteds := []test{
		{`one two three four five six`, 9},  // 5 words, 4 spaces
		{`one two three four five, six`, 9}, // same
		{`one two three four, five six`, 8}, // 4 words, 3 spaces, one comma
		{`one two three`, 5},                // 3 words, 2 spaces
	}

	for _, expected := range expecteds {
		tokens := jargon.Tokenize(strings.NewReader(expected.input))

		f := newFilter(tokens)

		err := f.fill()
		if err != nil {
			t.Error(err)
		}

		count := f.buffer.Len()
		if count != expected.count {
			t.Errorf("for input %q, got %q; expected count of %d, got %d", expected.input, f.buffer.All(), expected.count, count)
		}
	}
}

func TestWordrun(t *testing.T) {
	type test struct {
		input string
		count int
	}

	expecteds := []test{
		{`one two three four five six`, 9},  // 5 words, 4 spaces
		{`one two three four, five six`, 7}, // 4 words, 3 spaces, no comma
		{`one two three`, 5},                // 3 words, 2 spaces
	}

	for _, expected := range expecteds {
		tokens := jargon.Tokenize(strings.NewReader(expected.input))

		f := newFilter(tokens)

		err := f.fill()
		if err != nil {
			t.Error(err)
		}

		run := f.wordrun()
		count := len(run)
		if count != expected.count {
			t.Errorf("for input %q, got %q; expected count of %d, got %d", expected.input, run, expected.count, count)
		}
	}
}
