package stackexchange

// dictionary satisfies the jargon.Dictionary interface
// Used in generated.go
type dictionary struct {
	tags     []string
	synonyms map[string]string
}

func (d *dictionary) GetTags() []string {
	return d.tags
}

func (d *dictionary) GetSynonyms() map[string]string {
	return d.synonyms
}

func (d *dictionary) MaxGramLength() int {
	return 3
}
