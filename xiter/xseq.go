package xiter

import (
	"iter"
	"slices"
)

type XSeq[K any] iter.Seq[K]

func FromSeq[T any](seq iter.Seq[T]) XSeq[T] {
	return XSeq[T](seq)
}

func FromSlice[T any, S ~[]T](s S) XSeq[T] {
	return XSeq[T](slices.Values(s))
}

func FromItems[T any](items ...T) XSeq[T] {
	return FromSlice(items)
}

func FromChan[T any](ch <-chan T) XSeq[T] {
	return func(yield func(T) bool) {
		for item := range ch {
			if !yield(item) {
				return
			}
		}
	}
}

func (s XSeq[V]) AsSeq() iter.Seq[V] {
	return iter.Seq[V](s)
}
