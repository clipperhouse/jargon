// Package contractions provides a jargon.TokenFilter to expand English contractions, such as "don't" → "does not"
package contractions

import (
	"strings"

	"github.com/clipperhouse/jargon"
)

// Expand converts single-token contractions to non-contracted version. Examples:
// don't → does not
// We’ve → We have
// SHE'S -> SHE IS
var Expander = &filter{}

type filter struct{}

// Filter converts single-token contractions to non-contracted version. Examples:
// don't → does not
// We’ve → We have
// SHE'S -> SHE IS
func (f *filter) Filter(incoming *jargon.Tokens) *jargon.Tokens {
	t := &tokens{
		incoming: incoming,
		outgoing: &jargon.TokenQueue{},
	}
	return &jargon.Tokens{
		Next: t.next,
	}
}

type tokens struct {
	incoming *jargon.Tokens
	outgoing *jargon.TokenQueue
}

func (t *tokens) next() (*jargon.Token, error) {
	if t.outgoing.Len() > 0 {
		return t.outgoing.Pop(), nil
	}

	token, err := t.incoming.Next()
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, nil
	}

	// Try case-sensitive
	found, err := t.tryExpansion(token, false)
	if err != nil {
		return nil, err
	}
	if !found {
		// Try case-insensitive
		_, err := t.tryExpansion(token, true)
		if err != nil {
			return nil, err
		}
	}

	if t.outgoing.Len() > 0 {
		return t.outgoing.Pop(), nil
	}

	return token, nil
}

func (t *tokens) tryExpansion(token *jargon.Token, ignoreCase bool) (bool, error) {
	key := token.String()
	if ignoreCase {
		key = strings.ToLower(key)
	}

	expansion, found := variations[key]

	if found {
		tokens, err := jargon.TokenizeString(expansion).ToSlice()
		if err != nil {
			return found, err
		}
		t.outgoing.Push(tokens...)
	}

	return found, nil
}
