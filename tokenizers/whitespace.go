package tokenizers

import "unicode"

// WhiteSpace is a tokenizer with whitespace as a delimiter (as defined in Go's `unicode` package)
// It uses the definition of whitespace from the unicode package: https://golang.org/pkg/unicode/#IsSpace
// Space are tokens and will be returned alongside words; it's up to the caller to ignore them if desired
var WhiteSpace = &delimited{isDelimiter: unicode.IsSpace}
