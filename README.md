# Jargon

Jargon is a text pipeline, focused on recognizing variations on canonical and synonymous terms.

For example, jargon lemmatizes `react`, `React.js`, `React JS` and `REACTJS` to a canonical `reactjs`.

## Install

Binaries are available on the [Releases page](https://github.com/clipperhouse/jargon/releases).

If you have [Homebrew](https://brew.sh):
```
brew install clipperhouse/tap/jargon
```

If you have a [Go installation](https://golang.org/doc/install):
```
go install github.com/clipperhouse/jargon/cmd/jargon
```

To display usage, simply type:

```bash
jargon
```

Example:

```bash
curl -s https://en.wikipedia.org/wiki/Computer_programming | jargon -html -stack -lemmas -lines
```

[CLI usage and details...](https://github.com/clipperhouse/jargon/tree/master/cmd/jargon)

## In your code

See [GoDoc](https://godoc.org/github.com/clipperhouse/jargon). Example:

```go
import (
	"fmt"
	"log"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/stackoverflow"
)
 
text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
stream := jargon.TokenizeString(text).Filter(stackoverflow.Tags)

// Loop while Scan() returns true. Scan() will return false on error or end of tokens.
for stream.Scan() {
	token := stream.Token()
	// Do stuff with token
	fmt.Print(token)
}

if err := stream.Err(); err != nil {
	// Because the source is I/O, errors are possible
	log.Fatal(err)
}

// As an iterator, a token stream is 'forward-only'; once you consume a token, you can't go back.

// See also the convenience methods String, ToSlice, WriteTo
```

## Token filters

Canonical terms (lemmas) are looked up in token filters. Several are available:

[Stack Overflow technology tags](https://pkg.go.dev/github.com/clipperhouse/jargon/filters/stackoverflow)
  - `Ruby on Rails → ruby-on-rails`
  - `ObjC → objective-c`

[Contractions](https://pkg.go.dev/github.com/clipperhouse/jargon/filters/contractions)
  - `Couldn’t → Could not`

[ASCII fold](https://pkg.go.dev/github.com/clipperhouse/jargon/filters/ascii)
  - `café → cafe`

[Stem](https://pkg.go.dev/github.com/clipperhouse/jargon/filters/stemmer)
  - `Manager|management|manages → manag`

To implement your own, see the [Filter type](https://godoc.org/github.com/clipperhouse/jargon/#Filter).

## Performance

`jargon` is designed to work in constant memory, regardless of input size. It buffers input and streams tokens.

Execution time is designed to O(n) on input size. It is I/O-bound. In your code, you control I/O and performance implications by the `Reader` you pass to Tokenize.

## Tokenizer

Jargon includes a tokenizer based partially on [Unicode text segmentation](https://unicode.org/reports/tr29/). It’s good for many common cases.

It preserves all tokens verbatim, including whitespace and punctuation, so the original text can be reconstructed with fidelity (“round tripped”).

## Background

When dealing with technical terms in text – say, a job listing or a resume – it’s easy to use different words for the same thing. This is acute for things like “react” where it’s not obvious what the canonical term is. Is it React or reactjs or react.js?

This presents a problem when **searching** for such terms. _We_ know the above terms are synonymous but databases don’t.

A further problem is that some n-grams should be understood as a single term. We know that “Objective C” represents **one** technology, but databases naively see two words.

## What’s it for?

- Recognition of domain terms in text
- NLP for unstructured data, when we wish to ensure consistency of vocabulary, for statistical analysis.
- Search applications, where searches for “Ruby on Rails” are understood as an entity, instead of three unrelated words, or to ensure that “React” and “reactjs” and “react.js” and handled synonmously.
