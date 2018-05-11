package stackexchange

import (
	"testing"
)

// Run this test to do the codegen: go test -run ^TestWriteDictionary$
func TestWriteDictionary(t *testing.T) {
	err := writeDictionary()

	if err != nil {
		t.Error(err)
	}
}
