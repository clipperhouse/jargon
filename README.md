# jargon
A lemmatizer for technology terms in text or HTML, written in Go.

## Problem
When dealing with technology terms in text – say, a job listing or a resume or a document – 
it’s easy to use different words for the same thing. This is acute for things like “react” where it’s not obvious
what the canonical term is. Is it React or reactjs or react.js?

This presents a problem when **searching** for such terms. _We_ know the above terms are synonymous but databases don’t.

A further problem is that some ngrams should be understood as a single term. We know that “Ruby on Rails” represents 
**one** technology, but databases naively see three words.


## Try it

```go
package main

import (
    "fmt"
    
    "github.com/clipperhouse/jargon"
)

func main() {
    text := `Let’s talk about Ruby on Rails and ASPNET MVC.`
    result := jargon.Lemmatize(text)
    fmt.Println(result)
}
```

## Prior art
This is effectively a problem of synonyms. Search-oriented databases like Elastic handle this problem with [analyzers](https://www.elastic.co/guide/en/elasticsearch/reference/current/analysis-analyzers.html).

In NLP, it’s handled by [stemmers](https://en.wikipedia.org/wiki/Stemming) or [lemmatizers](https://en.wikipedia.org/wiki/Lemmatisation). There, the goal is to replace variations of a term (manager, management, managing) with a single canonical version.

Recognizing mutli-words-as-a-single-term (“Ruby on Rails”) is [named-entity recognition](https://en.wikipedia.org/wiki/Named-entity_recognition).

## Who’s it for?
Dunno yet, some ideas…

- Data scientists doing NLP on unstructured data, who want to ensure consistency of vocabulary, for statistical analysis.
- Search applications, where searches for “Ruby on Rails” are understood as an entity, instead of three unrelated words, or to ensure that “React” and “reactjs” and “react.js” and handled synonmously.

## How it works

### Lemmatizer
A lemmatizer is constructed using a Dictionary (below), which contains all the synonym data, as well as some rules.

### Dictionary
Dictionary is an interface with the following methods:

`GetTags()` : The list of canonical terms

`GetSynonyms()` : A map of synonyms to their canonical terms (lemmas)

`MaxGramLength()` : The maximum number of individual words that the lemmatizer will attempt to join into a single term. For example, if we want to recognize Ruby on Rails, we’d want an n-gram length of 3.

`Normalize()` : Defines ‘insensitivity’ rules when matching words against their canonical versions. In the default case, we want the lookups to be case-insensitive, as well as insensitive to dots and dashes. So `NodeJS` and `node.js` are handled identically.

### Tokenizers
Before we can lemmatize text, we need it to separated into words and punctuation, which we call tokens. Getting this right matters! There are two built-in tokenizers.

- TechProse: follows typical rules of English, where spaces and punctuation define the separation of words. It mostly relies on Unicode’s definitions. Not just prose, though: these rules should™️ work for delimited files like CSV.

- TechHTML: Tokenizes HTML, and in turn, tokenizes text nodes using TechProse above.

Importantly, these tokenizers capture all of the text, including white space, so it can be reconstructed with fidelity.
