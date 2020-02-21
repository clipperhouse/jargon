# Jargon

Jargon is a **lemmatizer**, useful for recognizing variations on canonical and synonymous terms.

For example, jargon lemmatizes `react`, `React.js`, `React JS` and `REACTJS` to a canonical `reactjs`.

Jargon uses Stack Overflow tags & synonyms, and implements “insensitivity” to spaces, dots and dashes.

### Online demo

[Give it a try](https://clipperhouse.com/jargon/)

### Command line

```bash
go install github.com/clipperhouse/jargon/cmd/jargon
```

(Assumes a [Go installation](https://golang.org/dl/).)

To display usage, simply type:

```bash
jargon
```

```
jargon accepts piped UTF8 text from Stdin and pipes lemmatized text to Stdout

  Example: echo "I luv Rails" | jargon

Alternatively, use jargon 'standalone' by passing flags for inputs and outputs:

  -f string
    	Input file path
  -o string
    	Output file path
  -s string
    	A (quoted) string to lemmatize
  -u string
    	A URL to fetch and lemmatize

  Example: jargon -f /path/to/original.txt -o /path/to/lemmatized.txt
```

### In your code

See [GoDoc](https://godoc.org/github.com/clipperhouse/jargon).

## Token filters

Canonical terms (lemmas) are looked up in token filters. Three are available:

- [Stack Exchange technology tags](https://github.com/clipperhouse/jargon/stackexchange)
  - `Ruby on Rails → ruby-on-rails`
  - `ObjC → objective-c`
- [Contractions](https://github.com/clipperhouse/jargon/contractions)
  - `Couldn‘t → Could not`
- [Simple numbers](https://github.com/clipperhouse/jargon/numbers)
  - `Thirty-five hundred → 3500`

To implement your own, see the [jargon.TokenFilter interface](https://godoc.org/github.com/clipperhouse/jargon/#TokenFilter)

## Tokenizer

Jargon includes its own tokenizer, with an emphasis on handling technology terms correctly:

- C++, ASP.net, and other non-alphanumeric terms are recognized as single tokens
- #hashtags and @handles
- Simple URLs and email address are handled _pretty well_, though can be notoriously hard to get right

The tokenizer preserves all tokens verbatim, including whitespace and punctuation, so the original text can be reconstructed with fidelity (“round tripped”).

(It turns out that the above rules work well in structured text such as CSV and JSON.)

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

## What’s it for?

- Recognition of domain terms in text
- NLP for unstructured data, when we wish to ensure consistency of vocabulary, for statistical analysis.
- Search applications, where searches for “Ruby on Rails” are understood as an entity, instead of three unrelated words, or to ensure that “React” and “reactjs” and “react.js” and handled synonmously.
