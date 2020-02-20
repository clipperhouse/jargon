package jargon

// queue is a FIFO queue
type queue struct {
	tokens []*Token
}

func (q *queue) len() int {
	return len(q.tokens)
}

func (q *queue) peek() *Token {
	return q.tokens[0]
}

func (q *queue) pop() *Token {
	token := q.peek()
	q.drop(1)
	return token
}

func (q *queue) push(token *Token) {
	q.tokens = append(q.tokens, token)
}

func (q *queue) drop(n int) {
	// Optimization to avoid array resizing
	// Move the end to the beginning
	copy(q.tokens, q.tokens[n:])
	// Chop off the end
	q.tokens = q.tokens[:len(q.tokens)-n]
}
