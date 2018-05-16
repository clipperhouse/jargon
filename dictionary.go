package jargon

// Dictionary is a structure for containing tags and synonyms, for easy passing around
type Dictionary interface {
	GetTags() []string
	GetSynonyms() map[string]string
	// What is the longest n-gram (word run) to try to canonicalize
	MaxGramLength() int
}
