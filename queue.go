package jargon

// TokenQueue is a FIFO queue
type TokenQueue struct {
	tokens []*Token
}

func (q *TokenQueue) All() []*Token {
	return q.tokens
}

func (q *TokenQueue) At(i int) *Token {
	return q.tokens[i]
}

func (q *TokenQueue) Len() int {
	return len(q.tokens)
}

func (q *TokenQueue) First() *Token {
	return q.tokens[0]
}

func (q *TokenQueue) Pop() *Token {
	token := q.First()
	q.Drop(1)
	return token
}

func (q *TokenQueue) Push(token *Token) {
	q.tokens = append(q.tokens, token)
}

func (q *TokenQueue) Drop(n int) {
	// Optimization to avoid array resizing
	// Move the end to the beginning
	copy(q.tokens, q.tokens[n:])
	// Chop off the end
	q.tokens = q.tokens[:len(q.tokens)-n]
}

func (src *TokenQueue) PopTo(dst *TokenQueue) {
	token := src.Pop()
	dst.Push(token)
}

func (src *TokenQueue) FlushTo(dst *TokenQueue) {
	for range src.All() {
		src.PopTo(dst)
	}
}
