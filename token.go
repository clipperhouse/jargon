package jargon

// Tokens represents an 'iterator' of Token. Call .Next() until it returns nil.
type Tokens struct {
	Next func() *Token
}

func (tokens Tokens) ToSlice() []*Token {
	var result []*Token
	for {
		t := tokens.Next()
		if t == nil {
			break
		}
		result = append(result, t)
	}
	return result
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

// IsPunct indicates that the token should be considered 'breaking' of a run of words. Mostly uses
// Unicode's definition of punctuation, with some exceptions for our purposes.
func (t *Token) IsPunct() bool {
	return t.punct
}

// IsSpace indicates that the token consists entirely of white space, as defined by the unicode package.
//
//A token can be both IsPunct and IsSpace -- for example, line breaks and tabs are punctuation for our purposes.
func (t *Token) IsSpace() bool {
	return t.space
}

// IsLemma indicates that the token is a lemma, i.e., a canonical term that replaced original token(s).
func (t Token) IsLemma() bool {
	return t.lemma
}
