# Jargon

Jargon offers a **tokenizer** for Go, with an emphasis on handling technology terms correctly:

- C++, ASP.net, and other non-alphanumeric terms are recognized as single tokens
- #hashtags and @handles
- Simple URLs and email address are handled _pretty well_, though can be notoriously hard to get right

There is also an HTML tokenizer, which applies the above to text nodes in markup.

The tokenizer preserves all tokens verbatim, so that the original text can be reconstructed with fidelity (“round tripped”).

In turn, Jargon offers a **lemmatizer**, for recognizing canonical and synonymous terms. For example the n-gram “Ruby on Rails” becomes ruby-on-rails. It implements “insensitivity” to spaces, dots and dashes.

(It turns out™️ that the above rules apply well to structured text such as CSV and JSON.)

### Command line

```bash
go install github.com/clipperhouse/jargon/cmd/jargon
```

To display usage, simply type:

```bash
jargon
```

Use `-f` to lemmatize a file and pipe to stdout:

```bash
jargon -f file.txt
```

If you’re dealing with large files, you might wish to pipe the results into another file

```bash
jargon -f file.txt > result.txt
```

Use `-s` to lemmatize a string and pipe to stdout

```bash
jargon -s "Here is a string with Ruby and SQL"
```

Use `-u` to fetch a URL and lemmatize, and pipe to stdout

```bash
jargon -u https://en.wikipedia.org/wiki/Programming_language
```

### Online demo

[https://clipperhouse.com/jargon](https://clipperhouse.com/jargon)

### In your code

[GoDoc](https://godoc.org/github.com/clipperhouse/jargon)

```go
package main

import (
    "fmt"

    "github.com/clipperhouse/jargon"
    "github.com/clipperhouse/jargon/stackexchange"
)

var lem = jargon.NewLemmatizer(stackexchange.Dictionary)

func main() {
    text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
    r := strings.NewReader(text)
    tokens := jargon.Tokenize(r)

    // iterate over the resulting tokens, or pass on to the lemmatizer...

    lemmatized := lem.Lemmatize(tokens)
    for t := range lemmatized {
        fmt.Print(t)
    }
}
```

Jargon uses a streaming API – reader in, channel out. This is good for avoiding blowing out memory on large files.

## Background

When dealing with technology terms in text – say, a job listing or a resume –
it’s easy to use different words for the same thing. This is acute for things like “react” where it’s not obvious
what the canonical term is. Is it React or reactjs or react.js?

This presents a problem when **searching** for such terms. _We_ know the above terms are synonymous but databases don’t.

A further problem is that some n-grams should be understood as a single term. We know that “Objective C” represents
**one** technology, but databases naively see two words.

## Prior art

Existing tokenizers (such as Treebank), appear not to be round-trippable, i.e., are destructive. They also take a hard line on punctuation, so “ASP.net” would come out as two tokens instead of one. Of course I’d like to be corrected or pointed to other implementations.

Search-oriented databases like Elastic handle synonyms with [analyzers](https://www.elastic.co/guide/en/elasticsearch/reference/current/analysis-analyzers.html).

In NLP, it’s handled by [stemmers](https://en.wikipedia.org/wiki/Stemming) or [lemmatizers](https://en.wikipedia.org/wiki/Lemmatisation). There, the goal is to replace variations of a term (manager, management, managing) with a single canonical version.

Recognizing mutli-words-as-a-single-term (“Ruby on Rails”) is [named-entity recognition](https://en.wikipedia.org/wiki/Named-entity_recognition).

## Who’s it for?

Dunno yet, some ideas…

- Recognition of domain terms appearing in text
- NLP on unstructured data, when we wish to ensure consistency of vocabulary, for statistical analysis.
- Search applications, where searches for “Ruby on Rails” are understood as an entity, instead of three unrelated words, or to ensure that “React” and “reactjs” and “react.js” and handled synonmously.
