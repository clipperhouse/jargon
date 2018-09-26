## Numbers dictionary for Jargon

This package implements a Dictionary for use with the [jargon](https://github.com/clipperhouse/jargon) lemmatizer, intending to canonicalize simple number phrases appearing in text.

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

### Command line

Assuming you have installed the [Jargon CLI](https://github.com/clipperhouse/jargon#command-line), use the `-num` flag to specify this numbers dictionary.

```bash
echo "The U.S. population is around 3 hundred million people" | jargon -num
```

### In your code

```go
package main

import (
    "fmt"

    "github.com/clipperhouse/jargon"
    "github.com/clipperhouse/jargon/numbers"
)

var lem = jargon.NewLemmatizer(numbers.Dictionary)

func main() {
    text := `The U.S. population is around 3 hundred million people`
    r := strings.NewReader(text)
    tokens := jargon.Tokenize(r)

    // Or! Pass tokens on to the lemmatizer
    lemmas := lem.Lemmatize(tokens)
    for {
        lemma := tokens.Next()
        if lemma == nil {
            break
        }

        fmt.Print(lemma)
    }
}
```

### Implementation

The [Lookup](https://github.com/clipperhouse/jargon/blob/master/numbers/dictionary.go#L35) method satisfies the [jargon.Dictionary interface](https://github.com/clipperhouse/jargon/blob/master/dictionary.go).
