package stopwords

import (
	"strings"

	"github.com/clipperhouse/jargon"
)

// NewFilter creates a token filter for the supplied stop words
func NewFilter(stopwords []string, ignoreCase bool) *filter {
	includes := make(map[string]bool)
	for _, s := range stopwords {
		var key string
		if ignoreCase {
			key = strings.ToLower(s)
		} else {
			key = s
		}
		includes[key] = true
	}

	return &filter{
		includes:   includes,
		ignoreCase: ignoreCase,
	}
}

type filter struct {
	incoming   *jargon.TokenStream
	includes   map[string]bool
	ignoreCase bool
}

func (f *filter) Filter(incoming *jargon.TokenStream) *jargon.TokenStream {
	t := tokens{
		filter:   f,
		incoming: incoming,
	}
	return jargon.NewTokenStream(t.next)
}

type tokens struct {
	filter   *filter
	incoming *jargon.TokenStream
}

func (t *tokens) next() (*jargon.Token, error) {
	for {
		token, err := t.incoming.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			return nil, nil
		}

		key := token.String()
		if t.filter.ignoreCase {
			key = strings.ToLower(token.String())
		}

		if t.filter.includes[key] {
			// Word is stopped
			continue
		}

		return token, nil
	}
}
