This package generates a Dictionary for use with the [jargon](https://github.com/clipperhouse/jargon) lemmatizer, using technology tags & synonyms from Stack Exchange sites.

Examples:

- "Ruby on Rails" → "ruby-on-rails"
- "ObjC" → "objective-c"

It includes the most popular tags from Stack Overflow, Server Fault, Game Dev and and Data Science. We think that provides a good set of likely terms in technological text.

### Implementation

The [Lookup](https://github.com/clipperhouse/jargon/blob/master/stackexchange/dictionary.go#L16) method satisfies the [jargon.Dictionary interface](https://github.com/clipperhouse/jargon/blob/master/dictionary.go).

The dictionary is code-generated, pulling from the Stack Exchange API. Have a look at the [`writeDictionary` method](https://github.com/clipperhouse/jargon/blob/master/stackexchange/generator.go#L24).

There is a list of [stop words](https://github.com/clipperhouse/jargon/blob/master/stackexchange/stopwords.go), intended to avoid lemmatizing plain-English words that also happen to be tags, such as `this`.
