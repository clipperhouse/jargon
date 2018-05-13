package tokenizers

type token struct {
	value string
	punct bool
}

func (t token) Value() string {
	return t.value
}

func (t token) Punct() bool {
	return t.punct
}
