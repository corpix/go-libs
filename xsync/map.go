package xsync

import (
	"sync"
)

type Map[K comparable, V any] struct{ m sync.Map }

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	return v.(V), ok
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// Clear deletes all the entries, resulting in an empty Map.
func (m *Map[K, V]) Clear() {
	m.m.Clear()
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, ok := m.m.LoadOrStore(key, value)
	return v.(V), ok
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, ok := m.m.LoadAndDelete(key)
	return v.(V), ok
}

// Delete deletes the value for a key.
// If the key is not in the map, Delete does nothing.
func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, ok := m.m.Swap(key, value)
	return v.(V), ok
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (m *Map[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	return m.m.CompareAndSwap(key, old, new)
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool { return f(key.(K), value.(V)) })
}

func (m *Map[K, V]) Keys(f func(key K) bool) {
	m.Range(func(key K, _ V) bool { return f(key) })
}

func (m *Map[K, V]) Values(f func(values V) bool) {
	m.Range(func(_ K, value V) bool { return f(value) })
}
