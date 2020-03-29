package jargon_test

import (
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
)

// TODO: test ordering

func TestLeading(t *testing.T) {
	text := `Hi. This is a test of .net, and #hashtag and @handle, and React.js and .123.`

	r := strings.NewReader(text)
	tokens := jargon.Tokenize(r)

	type test struct {
		value string
		found bool
	}

	expecteds := []test{
		{"Hi", true},
		{".", true},
		{"Hi.", false},
		{".net", true},
		{"net", false},
		{"#hashtag", true},
		{"hashtag", false},
		{"@handle", true},
		{"handle", false},
		{"React.js", true},
		{"React.js.", false},
		{".123", true},
	}

	got := map[string]bool{}

	for {
		token, err := tokens.Next()
		if err != nil {
			t.Error(err)
		}
		if token == nil {
			break
		}

		got[token.String()] = true
	}

	for _, expected := range expecteds {
		if got[expected.value] != expected.found {
			t.Errorf("expected finding %q to be %t", expected.value, expected.found)
		}
	}
}

func TestMiddle(t *testing.T) {
	text := `Hi. This is a test of asp.net, TCP/IP, and O'Brien's and possessives’ first_last and wishy-washy and email@domain.com.`

	r := strings.NewReader(text)
	tokens := jargon.Tokenize(r)

	type test struct {
		value string
		found bool
	}

	expecteds := []test{
		{"asp.net", true},
		{"asp", false},
		{"net", false},
		{"TCP/IP", true},
		{"TCP", false},
		{"/", false},
		{"IP", false},
		{"O'Brien's", true},
		{"O", false},
		{"Brien", false},
		{"'s", false},
		{"possessives", true},
		{"’", true},
		{"possessives’", false},
		{"first_last", true},
		{"first", false},
		{"last", false},
		{"wishy-washy", false},
		{"wishy", true},
		{"-", true},
		{"washy", true},
		{"email", true},
		{"email@", false},
		{"@", true},
		{"domain.com", true},
		{"@domain.com", false},
	}

	got := map[string]bool{}

	for {
		token, err := tokens.Next()
		if err != nil {
			t.Error(err)
		}
		if token == nil {
			break
		}

		got[token.String()] = true
	}

	for _, expected := range expecteds {
		if got[expected.value] != expected.found {
			t.Errorf("expected finding %q to be %t", expected.value, expected.found)
		}
	}
}

func TestTrailing(t *testing.T) {
	text := `Hi. This is a test of F# and C++.`

	r := strings.NewReader(text)
	tokens := jargon.Tokenize(r)

	type test struct {
		value string
		found bool
	}

	expecteds := []test{
		{"Hi", true},
		{".", true},
		{"Hi.", false},
		{"F#", true},
		{"F", false},
		{"#", false},
		{"C++", true},
		{"C", false},
		{"+", false},
		{"++", false},
	}

	got := map[string]bool{}

	for {
		token, err := tokens.Next()
		if err != nil {
			t.Error(err)
		}
		if token == nil {
			break
		}

		s := token.String()
		got[s] = true
	}

	for _, expected := range expecteds {
		if got[expected.value] != expected.found {
			t.Errorf("expected finding %q to be %t", expected.value, expected.found)
		}
	}
}
