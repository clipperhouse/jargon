package jargon

type where struct {
	incoming  *Tokens
	predicate func(*Token) bool
}

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
		t, err := w.incoming.Next()
		if err != nil {
			return nil, err
		}
		if t == nil {
			break
		}

		if w.predicate(t) {
			return t, nil
		}
	}

	return nil, nil
}
