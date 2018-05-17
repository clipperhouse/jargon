# jargon
A lemmatizer for technology terms

## Problem
When dealing with technology terms in text – say, a job listing or a resume or a document – 
it’s easy to use different words for the same thing. This is acute for things like “react” where it’s not obvious
what the canonical term is. Is it React or reactjs or react.js?

This presents a problem when **searching** for such terms. _We_ know the above terms are synonymous but databases don’t.

A further problem is that some ngrams should be understood as a single term. We know that “Ruby on Rails” represents 
**one** technology, but databases naively see three words.

## Prior art
This is effectively a problem of synonyms. Search-oriented databases like Elastic handle this problem with [analyzers](https://www.elastic.co/guide/en/elasticsearch/reference/current/analysis-analyzers.html).

In NLP, it’s handled by [stemmers](https://en.wikipedia.org/wiki/Stemming) or [lemmatizers](https://en.wikipedia.org/wiki/Lemmatisation). There, the goal is to replace variations of a term (manager, management, managing) with a single canonical version.

Recognizing mutli-words-as-a-single-term (“Ruby on Rails”) is [named-entity recognition](https://en.wikipedia.org/wiki/Named-entity_recognition).

## Who’s it for?
Dunno yet, some ideas…

- Data scientists doing NLP on unstructured data, who want to ensure consistency of vocabulary, for statistical analysis.
- Search applications, where searches for “Ruby on Rails” are understood as an entity, instead of three unrelated words, or to ensure that “React” and “reactjs” and “react.js” and handled synonmously.