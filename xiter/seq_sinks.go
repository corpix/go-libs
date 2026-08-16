package xiter

import "iter"

func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		index := 0
		for v := range seq {
			if !yield(index, v) {
				return
			}
			index += 1
		}
	}
}

func GroupByKey[T any, K comparable](s iter.Seq[T], mapper func(T) K, capacity int) map[K][]T {
	res := make(map[K][]T, capacity)

	for v := range s {
		key := mapper(v)
		res[key] = append(res[key], v)
	}

	return res
}

// Batched разбивает исходную последовательность на части размером не более n элементов
func Batched[T, V any](s iter.Seq[T], n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		batch := make([]T, 0, n)
		for item := range s {
			batch = append(batch, item)

			if len(batch) >= n {
				if !yield(batch) {
					return
				}
				batch = make([]T, 0, n)
			}
		}

		if len(batch) != 0 {
			yield(batch)
		}
	}
}

func Collect2[K any, T any](s iter.Seq2[K, T]) ([]K, []T) {
	var ks []K
	var ts []T

	for k, t := range s {
		ks = append(ks, k)
		ts = append(ts, t)
	}

	return ks, ts
}
