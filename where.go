package jargon

type where struct {
	incoming  *Tokens
	predicate func(*Token) bool
}

// Where filters a stream of Tokens that match a predicate
func (incoming *Tokens) Where(predicate func(*Token) bool) *Tokens {
	w := &where{
		incoming:  incoming,
		predicate: predicate,
	}
	return &Tokens{
		Next: w.next,
	}
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
