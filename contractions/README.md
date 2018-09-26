## Contractions dictionary for Jargon

This package implements a Dictionary for use with the [jargon](https://github.com/clipperhouse/jargon) lemmatizer, expanding common English contractions into separate words.

Examples:

- don't → does not
- We’ve → We have
- SHE'S -> SHE IS

It handles lower, Title and UPPER case tokens, as well as straight ' and smart ’ apostrophes.

### Command line

Assuming you have installed the [Jargon CLI](https://github.com/clipperhouse/jargon#command-line), use the `-cont` flag to specify this numbers dictionary.

```bash
echo "I would've called but he's away from his phone" | jargon -cont
```

### In your code

```go
package main

import (
    "fmt"

    "github.com/clipperhouse/jargon"
    "github.com/clipperhouse/jargon/contractions"
)

var lem = jargon.NewLemmatizer(contractions.Dictionary)

func main() {
    text := "I would've called but he's away from his phone"
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

The [Lookup method](https://github.com/clipperhouse/jargon/blob/master/contractions/dictionary.go#L7) satisfies the [jargon.Dictionary interface](https://github.com/clipperhouse/jargon/blob/master/dictionary.go).

Here is the [base list of contractions](https://github.com/clipperhouse/jargon/blob/master/contractions/generator.go#L101). Variations (case, apostrophes) are code-generated.
