package xiter

import "iter"

func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}

func Repeate[V any](val V, n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for range n {
			if !yield(val) {
				return
			}
		}
	}
}

func Range(from, to, step int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for curr := from; curr < to; curr += step {
			if !yield(curr) {
				return
			}

		}
	}
}

func RangeCount(from, count int) iter.Seq[int] {
	return Range(from, from+count, 1)
}
