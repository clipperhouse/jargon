package jargon_test

// For testing internals, non-public members

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/ascii"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/numbers"
	"github.com/clipperhouse/jargon/stackexchange"
	"github.com/clipperhouse/jargon/stemmer"
	"github.com/clipperhouse/jargon/stopwords"
)

func TestLemmatize(t *testing.T) {
	dict := stackexchange.Tags

	original := `Here is the story of Ruby on Rails nodeJS, "Java Script", html5 and ASPNET mvc plus TCP/IP.`
	r1 := strings.NewReader(original)
	tokens := jargon.Tokenize(r1)

	got, err := tokens.Lemmatize(dict).ToSlice()
	if err != nil {
		t.Error(err)
	}

	r2 := strings.NewReader(`Here is the story of ruby-on-rails node.js, "javascript", html5 and asp.net-mvc plus tcpip.`)
	expected, err := jargon.Tokenize(r2).ToSlice()
	if err != nil {
		t.Error(err)
	}

	if !equals(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}

	lemmas := []string{"ruby-on-rails", "node.js", "javascript", "html5", "asp.net-mvc"}

	lookup := make(map[string]*jargon.Token)
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

func TestLemmatizeString(t *testing.T) {
	dict := stackexchange.Tags

	s := `Here is the story of Ruby on Rails.`

	r := strings.NewReader(s)
	tokens := jargon.Tokenize(r)
	lemmatized := tokens.Lemmatize(dict)

	s1, err := lemmatized.String()
	if err != nil {
		t.Error(err)
	}

	s2 := jargon.LemmatizeString(s)

	if s1 != s2 {
		t.Errorf("Lemmatize and LemmatizeString should give the same result")
	}
}

func TestRetokenize(t *testing.T) {
	dict := contractions.Expander

	original := `Would've but also won't`
	r1 := strings.NewReader(original)
	tokens := jargon.Tokenize(r1)

	got, err := tokens.Lemmatize(dict).ToSlice()
	if err != nil {
		t.Error(err)
	}

	r2 := strings.NewReader(`Would have but also will not`)
	expected, err := jargon.Tokenize(r2).ToSlice()
	if err != nil {
		t.Error(err)
	}

	if !equals(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func BenchmarkLemmatizer(b *testing.B) {
	dict := stackexchange.Tags

	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		tokens := jargon.Tokenize(r)
		consume(tokens.Lemmatize(dict))
	}
}

func TestCSV(t *testing.T) {
	dict := stackexchange.Tags

	original := `"Ruby on Rails", 3.4, "foo"
"bar",42, "java script"`
	r1 := strings.NewReader(original)
	tokens := jargon.Tokenize(r1)

	got, err := tokens.Lemmatize(dict).ToSlice()
	if err != nil {
		t.Error(err)
	}

	r2 := strings.NewReader(`"ruby-on-rails", 3.4, "foo"
"bar",42, "javascript"`)
	expected, err := jargon.Tokenize(r2).ToSlice()
	if err != nil {
		t.Error(err)
	}

	if !equals(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func TestTSV(t *testing.T) {
	dict := stackexchange.Tags

	original := `Ruby on Rails	3.4	foo
ASPNET	MVC
bar	42	java script`
	r1 := strings.NewReader(original)
	tokens := jargon.Tokenize(r1)

	got, err := tokens.Lemmatize(dict).ToSlice()
	if err != nil {
		t.Error(err)
	}

	r2 := strings.NewReader(`ruby-on-rails	3.4	foo
asp.net	model-view-controller
bar	42	javascript`)
	expected, err := jargon.Tokenize(r2).ToSlice()
	if err != nil {
		t.Error(err)
	}

	if !equals(got, expected) {
		t.Errorf("Given tokens:\n%v\nexpected\n%v\nbut got\n%v", original, expected, got)
	}
}

func TestMultiple(t *testing.T) {
	s := `Here is the story of five and Rails and ASPNET in the CAFÉS and couldn't three hundred thousand.`

	got := jargon.LemmatizeString(s,
		stackexchange.Tags,
		contractions.Expander,
		numbers.Filter,
		stemmer.English,
		ascii.Fold,
	)

	expected := `here is the stori of 5 and ruby-on-rail and asp.net in the cafe and could not 300000.`

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFill(t *testing.T) {
	original := `one two three four five `

	count, err := jargon.Tokenize(strings.NewReader(original)).TestCount()
	if err != nil {
		t.Error(err)
	}

	tokens := jargon.Tokenize(strings.NewReader(original))

	lem := jargon.TestNewLemmatizer(tokens, nil)

	for i := 0; i < count+2; i++ {
		err := jargon.TestFill(lem, i)
		if i <= count {
			if err != nil {
				t.Error(err)
			}
			if i != lem.TestBufferLen() {
				t.Errorf("i should equal len(buffer), but i == %d, len(buffer) == %d", i, lem.TestBufferLen())
			}
		} else { // i > count
			if err == nil {
				t.Errorf("for tokens of count %d, at i == %d, there should be an error on fill", i, count)
			}
		}
	}
}

func TestEmptyCanonical(t *testing.T) {
	stops := []string{
		"This",
		"a",
	}
	filter := stopwords.NewFilter(stops, true)

	input := jargon.Tokenize(strings.NewReader("This is a test."))
	inputCount := 8 // tokens

	tokens := input.Lemmatize(filter)
	outputCount, err := tokens.TestCount()

	if err != nil {
		t.Error(err)
	}

	expectedCount := inputCount - 2

	if outputCount != expectedCount {
		t.Errorf("expected output count of %d, got %d", expectedCount, outputCount)
	}
}

func TestWordrun(t *testing.T) {
	original := `java script and, foo `
	r := strings.NewReader(original)
	tokens := jargon.Tokenize(r)

	var none []string // DeepEqual doesn't see zero-length slices as equal; need 'nil'

	type expected struct {
		words    []string
		consumed int
		err      error
	}

	expecteds := map[int]expected{
		4: {none, 0, jargon.TestErrInsufficient},       // attempting to get 4 should fail
		3: {[]string{"java", "script", "and"}, 5, nil}, // attempting to get 3 should work, consuming 5
		2: {[]string{"java", "script"}, 3, nil},        // attempting to get 2 should work, consuming 3 tokens (incl the space)
		1: {[]string{"java"}, 1, nil},                  // attempting to get 1 should work, and consume only that token
	}

	lem := jargon.TestNewLemmatizer(tokens, nil)

	for desired, expected := range expecteds {
		words, consumed, err := jargon.TestWordrun(lem, desired)
		if err != expected.err {
			t.Error(err)
		}

		if !reflect.DeepEqual(expected.words, words) {
			t.Errorf("desired %d, expected to get %s, got %s", desired, expected.words, words)
		}

		if expected.consumed != consumed {
			t.Errorf("desired %d, expected to consume %d tokens, got %d", desired, expected.consumed, consumed)
		}
	}
}

func BenchmarkLemmatize(b *testing.B) {
	dict := stackexchange.Tags

	file, err := ioutil.ReadFile("testdata/wikipedia.txt")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(file)
		tokens := jargon.Tokenize(r)
		consume(tokens.Lemmatize(dict))
	}
}

func consume(tokens *jargon.Tokens) error {
	for {
		t, err := tokens.Next()
		if err != nil {
			return err
		}
		if t == nil {
			break
		}
	}
	return nil
}

func contains(value string, tokens []*jargon.Token) bool {
	for _, t := range tokens {
		if t.String() == value {
			return true
		}
	}
	return false
}

// Checks that value, punct and space are equal for two slices of token; deliberately does not check lemma
func equals(a, b []*jargon.Token) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].String() != b[i].String() {
			return false
		}
		if a[i].IsPunct() != b[i].IsPunct() {
			return false
		}
		if a[i].IsSpace() != b[i].IsSpace() {
			return false
		}
		// deliberately not checking for IsLemma(); use reflect.DeepEquals
	}

	return true
}
