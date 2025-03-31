package gt

import "sync"

// Locked wraps access to underlying data in a mutex such that retrieving the
// data takes the mutex, then returning the data releases the mutex.
type Locked[T any] struct {
	mu   sync.RWMutex
	data T
}

// Returns a new Locked[T] that's ready for use
func NewLocked[T any](data T) Locked[T] {
	return Locked[T]{data: data}
}

// Takes an exclusive lock and sets the data to the given value
func (l *Locked[T]) Set(data T) {
	l.mu.Lock()
	l.data = data
	l.mu.Unlock()
}

// Acquires an exclusive lock on the data
func (l *Locked[T]) Get() (data T, unlock func()) {
	l.mu.Lock()

	data = l.data
	unlock = func() { l.mu.Unlock() }
	return
}

// Acquires an exclusive lock on the data for the duration of `f`
func (l *Locked[T]) With(f func(T)) {
	l.mu.Lock()
	defer l.mu.Unlock()

	f(l.data)
}

// Acquires a shared lock on the data for the duration of `f`
func (l *Locked[T]) RGet() (data T, unlock func()) {
	l.mu.RLock()

	data = l.data
	unlock = func() { l.mu.RUnlock() }
	return
}

// Acquires a shared lock on the data for the duration of `f`
func (l *Locked[T]) RWith(f func(T)) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f(l.data)
}
