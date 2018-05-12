package tokenizers

// Delimited is a simple tokenizer for simple cases, such as whitespace
type Delimited struct {
	isDelimiter func(rune) bool
}

func NewDelimited(isDelimiter func(rune) bool) *Delimited {
	return &Delimited{
		isDelimiter: isDelimiter,
	}
}

func (d *Delimited) Tokenize(text string) []string {
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
