package jargon

type where struct {
	incoming  *TokenStream
	predicate func(*Token) bool
}

// Where filters a stream of Tokens that match a predicate
func (incoming *TokenStream) Where(predicate func(*Token) bool) *TokenStream {
	w := &where{
		incoming:  incoming,
		predicate: predicate,
	}
	return NewTokenStream(w.next)
}

func (w *where) next() (*Token, error) {
	for {
		token, err := w.incoming.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			break
		}

		if w.predicate(token) {
			return token, nil
		}
	}

	return nil, nil
}
