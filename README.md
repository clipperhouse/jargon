# jargon
A lemmatizer for technology terms

## Problem
When dealing with technology terms in text – say, a job listing or a resume or structured tags – 
it’s easy to use the different words for the same thing. This is acute for things like “react” where it’s not obvious
what the canonical term is. Is it React or reactjs or react.js?

This presents a problem when **searching** for such terms. _We_ know the above terms are synonymous but databases don’t.

A further problem is that some ngrams should be understood as a single term. We know that “Ruby on Rails” represents 
**one** technology, but databases naively see three words.

## Prior art
This is effectively a problem of synonyms. Search-oriented databases like Elastic handle this problem with analyzers.

In NLP, it’s handled by stemmers or lemmatizers.

##Who’s it for?
Dunno yet, but some ideas…

- Data scientists doing NLP on unstructured data, who want to ensure consistency of terms
- Search applications
