package tokenqueue

import (
	"fmt"

	"github.com/clipperhouse/jargon"
)

// Token is an alias to jargon.Token for convenience
type Token = jargon.Token

// New creates a new TokenQueue
func New(tokens ...*Token) *TokenQueue {
	return &TokenQueue{
		Tokens: tokens,
	}
}

// TokenQueue is a FIFO queue
type TokenQueue struct {
	Tokens []*Token
}

// Any returns whether there are any tokens in the queue
func (q *TokenQueue) Any() bool {
	return q.Len() > 0
}

// Len is len(q.Tokens)
func (q *TokenQueue) Len() int {
	return len(q.Tokens)
}

// Pop returns the first token (front of) the queue, and removes it from the queue
func (q *TokenQueue) Pop() *Token {
	token := q.Tokens[0]
	q.Drop(1)
	return token
}

// Push appends a token to the end of the queue
func (q *TokenQueue) Push(tokens ...*Token) {
	q.Tokens = append(q.Tokens, tokens...)
}

// Drop removes n elements from the front of the queue
func (q *TokenQueue) Drop(n int) {
	// Optimization to avoid array resizing
	// Move the end to the beginning
	copy(q.Tokens, q.Tokens[n:])
	// Chop off the end
	q.Tokens = q.Tokens[:len(q.Tokens)-n]
}

// Clear drops all tokens from the queue
func (q *TokenQueue) Clear() {
	q.Tokens = q.Tokens[:0]
}

// PopTo moves a token from one queue to another
func (q *TokenQueue) PopTo(dst *TokenQueue) {
	token := q.Pop()
	dst.Push(token)
}

// FlushTo moves all tokens from one queue to another
func (q *TokenQueue) FlushTo(dst *TokenQueue) {
	dst.Tokens = append(dst.Tokens, q.Tokens...)
	q.Clear()
}

func (q *TokenQueue) String() string {
	return fmt.Sprintf("%q", q.Tokens)
}
