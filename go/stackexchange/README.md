This package generates a Dictionary for use with the [jargon](https://github.com/clipperhouse/jargon) lemmatizer, using tags & synonyms from Stack Exchange sites.

In particular, it uses the most popular tags from Stack Overflow, Server Fault, Game Dev and and Data Science. We think that provides a pretty good sample of likely term is technology text.

It does this via code generation, pulling from the Stack Exchange API.

The easiest way to generate code is to invoke the tests using `go test`.

There is a list of stop words, intended to avoid lemmatizing plain-English words that also happen to be tags, such as `this`.