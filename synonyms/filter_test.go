package synonyms_test

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/synonyms"
)

func TestFilter(t *testing.T) {
	original := `Here is the story of Ruby on Rails node JS, "Java Script", C++ cpp fsharp html5 and ASPNET mvc plus TCP/IP.`
	r1 := strings.NewReader(original)
	tokens := jargon.Tokenize(r1)

	ignore := []rune{'-', ' ', '.', '/'}
	syns, err := synonyms.NewFilter(mappings, true, ignore)
	if err != nil {
		t.Error(err)
	}

	got, err := syns.Filter(tokens).ToSlice()
	if err != nil {
		t.Error(err)
	}

	t.Log(got)
	return

	r2 := strings.NewReader(`Here is the story of ruby-on-rails node.js, "javascript", html and asp.net-mvc plus tcp.`)
	expected, err := jargon.Tokenize(r2).ToSlice()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}

	lemmas := []string{"ruby-on-rails", "node.js", "javascript", "html", "asp.net-mvc"}

	lookup := make(map[string]*jargon.Token)
	for _, g := range got {
		lookup[g.String()] = g
	}

	for _, lemma := range lemmas {
		l, ok := lookup[lemma]
		if !ok {
			t.Errorf("Expected to find lemma %q, but did not", lemma)
		}
		if !l.IsLemma() {
			t.Errorf("Expected %q to be identified as a lemma, but it was not", lemma)
		}
	}
}

func BenchmarkFilter(b *testing.B) {
	file, err := ioutil.ReadFile("../testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	ignore := []rune{'-', ' ', '.', '/'}
	syns, err := NewFilter(mappings, true, ignore)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		tokens := jargon.Tokenize(r)
		syns.Filter(tokens).Count() // consume them
	}
}
