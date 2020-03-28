package jargon

// TokenQueue is a FIFO queue
type TokenQueue struct {
	tokens []*Token
}

// All returns a slice of all tokens in the queue
func (q *TokenQueue) All() []*Token {
	return q.tokens
}

// Any returns whether there are any tokens in the queue
func (q *TokenQueue) Any() bool {
	return len(q.tokens) > 0
}

// Pop returns the first token (front of) the queue, and removes it from the queue
func (q *TokenQueue) Pop() *Token {
	token := q.tokens[0]
	q.Drop(1)
	return token
}

// Push appends a token to the end of the queue
func (q *TokenQueue) Push(tokens ...*Token) {
	q.tokens = append(q.tokens, tokens...)
}

// Drop removes n elements from the front of the queue
func (q *TokenQueue) Drop(n int) {
	// Optimization to avoid array resizing
	// Move the end to the beginning
	copy(q.tokens, q.tokens[n:])
	// Chop off the end
	q.tokens = q.tokens[:len(q.tokens)-n]
}

// PopTo moves a token from one queue to another
func (q *TokenQueue) PopTo(dst *TokenQueue) {
	token := q.Pop()
	dst.Push(token)
}

// FlushTo moves all tokens from one queue to another
func (q *TokenQueue) FlushTo(dst *TokenQueue) {
	for range q.All() {
		q.PopTo(dst)
	}
}
