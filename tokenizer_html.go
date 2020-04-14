package jargon

import (
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// TokenizeHTML tokenizes HTML. Text nodes are tokenized using jargon.Tokenize; everything else (tags, comments) are left verbatim.
// It returns a Tokens, intended to be iterated over by calling Next(), until nil.
// It returns all tokens (including white space), so text can be reconstructed with fidelity. Ignoring (say) whitespace is a decision for the caller.
func TokenizeHTML(r io.Reader) *TokenStream {
	t := &htokenizer{
		htokenizer: html.NewTokenizer(r),
	}
	return NewTokenStream(t.next)
}

type htokenizer struct {
	htokenizer *html.Tokenizer
	ttokens    *TokenStream
	parent     atom.Atom
}

// next is the implementation of the Tokens interface. To iterate, call until it returns nil
func (t *htokenizer) next() (*Token, error) {
	// Are we "inside" a text node?
	if t.ttokens != nil {
		ttoken, err := t.ttokens.Next()
		if err != nil {
			return nil, err
		}
		if ttoken != nil {
			return ttoken, nil
		}

		// Done with text node
		t.ttokens = nil
	}

	htype := t.htokenizer.Next()

	if htype == html.ErrorToken {
		err := t.htokenizer.Err()
		if err == io.EOF {
			// No problem
			return nil, nil
		}
		return nil, err
	}

	htoken := t.htokenizer.Token()

	switch htoken.Type {
	case html.StartTagToken:
		// Record that we are entering script or style blocks; don't tokenize text
		if htoken.DataAtom == atom.Script || htoken.DataAtom == atom.Style {
			t.parent = htoken.DataAtom
		}
	case html.TextToken:
		switch t.parent {
		case atom.Script, atom.Style:
			// Don't tokenize script and style blocks, just return as one big string
			token := &Token{
				value: htoken.String(),
				punct: false,
				space: false,
			}
			return token, nil
		default:
			t.ttokens = TokenizeString(htoken.String())
			return t.ttokens.Next()
		}
	case html.EndTagToken:
		if htoken.DataAtom == t.parent {
			htoken.DataAtom = 0
		}
	}

	// Everything else is punct for our purposes
	token := &Token{
		value: htoken.String(),
		punct: true,
		space: false,
	}
	return token, nil
}
