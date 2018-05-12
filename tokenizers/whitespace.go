package tokenizers

import "unicode"

// WhiteSpaceDelimited is a tokenizer with whitespace as a delimiter
// It uses the definition of whitespace from the unicode package: https://golang.org/pkg/unicode/#IsSpace
// Space are tokens and will be returned alongside words; it's up to the caller to ignore them if desired
var WhiteSpaceDelimited = NewDelimited(unicode.IsSpace)
