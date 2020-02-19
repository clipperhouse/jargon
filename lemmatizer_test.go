package jargon

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/numbers"
	"github.com/clipperhouse/jargon/stackexchange"
)

func TestLemmatize(t *testing.T) {
	dict := stackexchange.Dictionary

	original := `Here is the story of Ruby on Rails nodeJS, "Java Script", html5 and ASPNET mvc plus TCP/IP.`
	r1 := strings.NewReader(original)
	tokens := Tokenize(r1)

	got := Lemmatize(tokens, dict).ToSlice()
	r2 := strings.NewReader(`Here is the story of ruby-on-rails node.js, "javascript", html5 and asp.net-mvc plus tcpip.`)
	expected := Tokenize(r2).ToSlice()

	if !equals(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}

	lemmas := []string{"ruby-on-rails", "node.js", "javascript", "html5", "asp.net-mvc"}

	lookup := make(map[string]*Token)
	for _, g := range got {
		lookup[g.String()] = g
	}

	for _, lemma := range lemmas {
		if !contains(lemma, got) {
			t.Errorf("Expected to find lemma %q, but did not", lemma)
		}
		if l, ok := lookup[lemma]; !ok || !l.IsLemma() {
			t.Errorf("Expected %q to be identified as a lemma, but it was not", lemma)
		}
	}
}

func TestRetokenize(t *testing.T) {
	dict := contractions.Dictionary

	original := `Would've but also won't`
	r1 := strings.NewReader(original)
	tokens := Tokenize(r1)

	got := Lemmatize(tokens, dict).ToSlice()
	r2 := strings.NewReader(`Would have but also will not`)
	expected := Tokenize(r2).ToSlice()

	if !equals(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func BenchmarkLemmatizer(b *testing.B) {
	dict := stackexchange.Dictionary

	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		tokens := Tokenize(r)
		consume(Lemmatize(tokens, dict))
	}
}

func TestCSV(t *testing.T) {
	dict := stackexchange.Dictionary

	original := `"Ruby on Rails", 3.4, "foo"
"bar",42, "java script"`
	r1 := strings.NewReader(original)
	tokens := Tokenize(r1)

	got := Lemmatize(tokens, dict).ToSlice()
	r2 := strings.NewReader(`"ruby-on-rails", 3.4, "foo"
"bar",42, "javascript"`)
	expected := Tokenize(r2).ToSlice()

	if !equals(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func TestTSV(t *testing.T) {
	dict := stackexchange.Dictionary

	original := `Ruby on Rails	3.4	foo
ASPNET	MVC
bar	42	java script`
	r1 := strings.NewReader(original)
	tokens := Tokenize(r1)

	got := Lemmatize(tokens, dict).ToSlice()
	r2 := strings.NewReader(`ruby-on-rails	3.4	foo
asp.net	model-view-controller
bar	42	javascript`)
	expected := Tokenize(r2).ToSlice()

	if !equals(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func TestMultiple(t *testing.T) {
	s := `Here is the story of five and Rails and ASPNET and couldn't three hundred thousand.`
	r := strings.NewReader(s)

	var tokens Tokens
	tokens = Tokenize(r)
	tokens = Lemmatize(tokens, stackexchange.Dictionary, contractions.Dictionary, numbers.Dictionary)

	expected := `Here is the story of 5 and ruby-on-rails and asp.net and could not 300000.`
	var got string

	for {
		t := tokens.Next()
		if t == nil {
			break
		}
		got += t.String()
	}

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestWordrun(t *testing.T) {
	original := `java script and `
	r := strings.NewReader(original)
	tokens := Tokenize(r)

	type result struct {
		expect   []string
		consumed int
		ok       bool
	}

	var none []string // DeepEqual doesn't see zero-length slices as equal; need 'nil'
	expecteds := map[int]result{
		4: {none, 0, false},                             // attempting to get 4 should fail
		3: {[]string{"java", "script", "and"}, 5, true}, // attempting to get 3 should work, consuming 5
		2: {[]string{"java", "script"}, 3, true},        // attempting to get 2 should work, consuming 3 tokens (incl the space)
		1: {[]string{"java"}, 1, true},                  // attempting to get 1 should work, and consume only that token
	}

	lem := &lemmatizer{
		incoming: tokens,
	}

	for take, expected := range expecteds {
		taken, consumed, ok := lem.wordrun(take)
		got := result{taken, consumed, ok}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Attempting to take %d words, expected %v but got %v", take, expected, got)
		}
	}
}

func consume(tokens Tokens) {
	for {
		t := tokens.Next()
		if t == nil {
			break
		}
	}
}
