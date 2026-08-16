package xiter

import (
	"iter"
	"maps"

	std "github.com/SlamJam/go-libs"
)

// Или интерфейс с WithCapacity(8)???
func Unique[T comparable](s iter.Seq[T], capacity int) iter.Seq[T] {
	return UniqueByKey(s, func(t T) T { return t }, capacity)
}

func UniqueByKey[T any, V comparable](s iter.Seq[T], mapper func(T) V, capacity int) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[V]std.Void, capacity)
		for v := range s {
			key := mapper(v)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = std.Empty()
			if !yield(v) {
				return
			}
		}
	}
}

func UniqueByKeyPrior[T any, V comparable](s iter.Seq[T], mapper func(T) V, cmp func(T, T) int, capacity int) iter.Seq[T] {
	seen := make(map[V]T, capacity)

	for item := range s {
		key := mapper(item)
		prev, ok := seen[key]

		if !ok || cmp(item, prev) > 0 {
			seen[key] = item
		}
	}

	return maps.Values(seen)
}

func Take[T any](s iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		// Если запросили Take(0), то не прочитаем ни одного элемента из s
		if n < 1 {
			return
		}

		for v := range s {
			n -= 1
			// Такая форма позволяет избежать чтения лишнего элемента из s
			if !yield(v) || n < 1 {
				return
			}
		}
	}
}

func TakeWhile[T any](s iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !pred(v) || !yield(v) {
				return
			}
		}
	}
}

func Skip[T any](s iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		n = max(0, n)
		for v := range s {
			if n > 0 {
				n -= 1
				continue
			}

			if !yield(v) {
				return
			}
		}
	}
}

func SkipWhile[T any](s iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		skip := true
		for v := range s {
			skip = skip && pred(v)

			if skip {
				continue
			}

			if !yield(v) {
				return
			}
		}
	}
}

func Filter[T any](s iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !pred(v) {
				continue
			}

			if !yield(v) {
				return
			}
		}
	}
}

func Map[T, V any](s iter.Seq[T], mapper func(T) V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			if !yield(mapper(v)) {
				return
			}
		}
	}
}

func Map2[K, T, V any](s iter.Seq2[K, T], mapper func(K, T) V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, t := range s {
			v := mapper(k, t)
			if !yield(k, v) {
				return
			}
		}
	}
}

func FlatMap[T, V any](s iter.Seq[T], mapper func(T) iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			seq := mapper(v)
			for t := range seq {
				if !yield(t) {
					return
				}
			}
		}
	}
}
