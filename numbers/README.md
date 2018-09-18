This package generates a Dictionary for use with the [jargon](https://github.com/clipperhouse/jargon) lemmatizer, intending to canonicalize simple number phrases appearing in text.

Examples:

- "three" → "3"
- "three thousand" → "3000"
- "thirty-five thousand" → "35000"
- "three hundred thousand" → "300000"
- "3 thousand" → "3000"
- "-3 thousand" → "-3000"
- "+3 thousand" → "3000"
- "2.54 million" → "2540000"
- "1,000,000" → "1000000"

### Implementation

The [Lookup](https://github.com/clipperhouse/jargon/blob/master/numbers/dictionary.go#L35) method satisfies the [jargon.Dictionary interface](https://github.com/clipperhouse/jargon/blob/master/dictionary.go).
