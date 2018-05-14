package jargon

// delimited is a simple tokenizer for simple cases, such as whitespace, without lookahead.
// Intended as a 'base class', so to speak.
type delimited struct {
	isDelimiter func(rune) bool
}

func (d *delimited) Tokenize(text string) []string {
	var tokens []string

	var current string
	for _, r := range text {
		if d.isDelimiter(r) {
			if len(current) > 0 {
				// Emit the previous token & reset
				tokens = append(tokens, current)
				current = ""
			}

			// Emit the delimiter
			tokens = append(tokens, string(r))
		} else {
			current += string(r)
		}
	}
	// Emit one more after falling out of the loop
	if len(current) > 0 {
		tokens = append(tokens, current)
		current = ""
	}

	return tokens
}
