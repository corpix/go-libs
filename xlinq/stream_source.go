package xlinq

import (
	"iter"
	"slices"

	"github.com/SlamJam/go-libs/xiter"
)

func FromIterable[V any](iterable Iterable[V]) stream[V] {
	var si sizeInfo

	if sizable, ok := iterable.(Sizable); ok {
		si.SetSize(sizable.Size())
	}

	return stream[V]{
		iterator: iterable.Iter(),
		sizeInfo: si,
	}
}

func FromItems[T any](items ...T) stream[T] {
	return FromSlice(items)
}

func FromSlice[T any](slice []T) stream[T] {
	return stream[T]{
		iterator: slices.Values(slice),
		sizeInfo: withSize(len(slice)),
	}
}

func FromSliceCopy[T any](slice []T) stream[T] {
	data := make([]T, len(slice))
	copy(data, slice)

	return FromSlice(data)
}

func FromSeqWithSize[T any](seq iter.Seq[T], size int) stream[T] {
	return stream[T]{
		iterator: seq,
		sizeInfo: withSize(size),
	}
}

func FromSeq[T any](seq iter.Seq[T]) stream[T] {
	return FromSeqWithSize(seq, 0)
}

func FromRepeate[T any](val T, count int) stream[T] {
	return FromSeqWithSize(
		xiter.Repeate(val, count),
		count,
	)
}

func FromEmpty[V any]() stream[V] {
	return stream[V]{
		iterator: xiter.Empty[V](),
	}
}

func FromRange(from, to, step int) stream[int] {
	return FromSeqWithSize(
		xiter.Range(from, to, step),
		(to-from)/step,
	)
}

func FromRangeCount(from, count int) stream[int] {
	return FromSeqWithSize(
		xiter.RangeCount(from, count),
		count,
	)
}
