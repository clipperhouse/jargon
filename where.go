package jargon

type where struct {
	stream    *TokenStream
	predicate func(*Token) bool
}

// Where filters a stream of Tokens that match a predicate
func (stream *TokenStream) Where(predicate func(*Token) bool) *TokenStream {
	w := &where{
		stream:    stream,
		predicate: predicate,
	}
	return NewTokenStream(w.next)
}

func (w *where) next() (*Token, error) {
	for {
		token, err := w.stream.Next()
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
