package techlemm

// Dictionary is a structure for containing tags and synonyms, for easy passing around
type Dictionary struct {
	Tags     []string
	Synonyms map[string]string
}

// NewDictionary instantiates a Dictionary pointer given tags and synonyms. Synonyms are a map of synonym (key) → canonical (value).
func NewDictionary(tags []string, synonyms map[string]string) *Dictionary {
	return &Dictionary{tags, synonyms}
}
