This package generates a Dictionary for use with the [jargon](https://github.com/clipperhouse/jargon) lemmatizer, using technology tags & synonyms from Stack Exchange sites.

Examples:

- "Ruby on Rails" → "ruby-on-rails"
- "ObjC" → "objective-c"

It includes the most popular tags from Stack Overflow, Server Fault, Game Dev and and Data Science. We think that provides a good set of likely terms in technological text.

### Implementation

The [Lookup](https://github.com/clipperhouse/jargon/blob/master/stackexchange/filter.go#L16) method satisfies the [jargon.TokenFilter interface](https://github.com/clipperhouse/jargon/blob/master/tokenfilter.go).

The dictionary is code-generated, pulling from the Stack Exchange API. Have a look at the [`writeDictionary` method](https://github.com/clipperhouse/jargon/blob/master/stackexchange/generator.go#L24).
