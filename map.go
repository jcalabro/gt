package gt

import (
	"sync"
)

// SafeMap is a thread-safe wrapper over Go's built in `map` type. This struct
// should not be instantiated directly (i.e. SafeMap[int, int]{} is incorrect),
// and instead `NewSafeMap` should always be used.
type SafeMap[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// Initializes a new SafeMap that's ready for use
func NewSafeMap[K comparable, V any]() SafeMap[K, V] {
	return SafeMap[K, V]{
		items: map[K]V{},
	}
}

// Adds an item to the map, overwriting any existing value, if any
func (m *SafeMap[K, V]) Put(key K, val V) {
	m.mu.Lock()
	m.items[key] = val
	m.mu.Unlock()
}

// Retrives the value from the map for the entry with the given key,
// if any. Returns `None` if the item does not exist in the map.
func (m *SafeMap[K, V]) Get(key K) Option[V] {
	m.mu.RLock()
	item, ok := m.items[key]
	m.mu.RUnlock()

	if ok {
		return Some(item)
	}
	return None[V]()
}
