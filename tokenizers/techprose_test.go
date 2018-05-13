package tokenizers

import (
	"fmt"
	"strings"
	"testing"
)

func TestTechProse(t *testing.T) {
	text := `Hi! This is Matt's story for foo@thing.stuff—whatevs—and Hindley-Millner. 
	I like C# and C++ 12.3 am @clipperhouse. Tech- insensitive. We like .Net around here, too.`
	lex := lex(text)

	got := lex.items
	fmt.Println(strings.Join(got, "➡"))
}
