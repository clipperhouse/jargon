package jargon

// Tokens represents an "iterator" interface for the results of tokenization or lemmatization
// Callers should call Next() until it returns nil, indicating the end of data
type Tokens interface {
	Next() *Token
}

// Token represents a piece of text with metadata.
type Token struct {
	value               string
	punct, space, lemma bool
}

// String is the string value of the token
func (t *Token) String() string {
	return t.value
}

// IsPunct indicates that the token should be considered 'breaking' of a run of words; a delimiter. Mostly determined
// by the `unicode` package's definition, with some exceptions for our purposes.
func (t *Token) IsPunct() bool {
	return t.punct
}

// IsSpace indicates that the token consists entirely of white space (as defined by the `unicode `package).
//
//A token can be both IsPunct and IsSpace -- line breaks and tabs to be punctuation for our purposes.
func (t *Token) IsSpace() bool {
	return t.space
}

// IsLemma indicates that the token is a lemma, i.e., a canonical term that that replaced the original token(s).
func (t Token) IsLemma() bool {
	return t.lemma
}
