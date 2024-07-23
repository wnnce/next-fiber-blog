package cache

import (
	"sync"
)

type node[T any] struct {
	Val  *T
	Next *node[T]
}

func (n *node[T]) Clear() {
	n.Val = nil
	n.Next = nil
}

type queue[T any] struct {
	size int
	head *node[T]
	tail *node[T]
	pool *sync.Pool
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{
		size: 0,
		head: nil,
		tail: nil,
		pool: &sync.Pool{
			New: func() any {
				return new(node[T])
			},
		},
	}
}

func (q *queue[T]) Push(value T) {
	n := q.pool.Get().(*node[T])
	n.Val = &value
	if q.head == nil {
		q.head = n
		q.tail = n
	} else {
		q.tail.Next = n
		q.tail = n
	}
	q.size += 1
}

func (q *queue[T]) Peek() T {
	if q.size == 0 {
		panic("queue empty")
	}
	return *q.head.Val
}

func (q *queue[T]) Pop() T {
	if q.size == 0 {
		panic("queue empty")
	}
	n := q.head
	defer func() {
		n.Clear()
		q.pool.Put(n)
	}()
	if q.size == 1 {
		q.head = nil
	} else {
		q.head = q.head.Next
	}
	q.size -= 1
	return *n.Val
}

func (q *queue[T]) Size() int {
	return q.size
}

func (q *queue[T]) IsEmpty() bool {
	return q.size == 0
}
