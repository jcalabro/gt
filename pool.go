package gt

import (
	"sync"
)

// A thin wrapper around the stdlib's `sync.Pool`, but with type safety
type Pool[T any] struct {
	pool sync.Pool
}

// Creates a new `Pool` that's ready for use
func NewPool[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}
}

// Adds an item from the pool
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}

// Returns an item from the pool
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T) // nolint:errcheck
}
