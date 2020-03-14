package jargon

// Bridge to expose some internals to jargon_test package

var (
	// TestNewLemmatizer is exported only for testing
	TestNewLemmatizer = newLemmatizer
	// TestFill is exported only for testing
	TestFill = (*lemmatizer).fill
	// TestWordrun is exported only for testing
	TestWordrun = (*lemmatizer).wordrun
	// TestErrInsufficient is exported only for testing
	TestErrInsufficient = errInsufficient
)

// TestBufferLen is exported only for testing
func (lem *lemmatizer) TestBufferLen() int {
	return lem.buffer.Len()
}

// TestCount is exported only for testing
func (tokens Tokens) TestCount() (int, error) {
	count := 0
	for {
		t, err := tokens.Next()
		if err != nil {
			return count, err
		}
		if t == nil {
			break
		}
		count++
	}
	return count, nil
}
