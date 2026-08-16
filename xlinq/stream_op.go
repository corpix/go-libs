package xlinq

import (
	"iter"
	"slices"

	"github.com/SlamJam/go-libs/xiter"
)

func (s stream[T]) Where(pred func(T) bool) stream[T] {
	return stream[T]{
		iterator: xiter.Filter(s.iterator, pred),
		sizeInfo: s.sizeInfo,
	}
}

func (s stream[T]) Take(n int) stream[T] {
	return stream[T]{
		iterator: xiter.Take(s.iterator, n),
		sizeInfo: s.sizeInfo.CopyWithSize(min(n, s.Size())),
	}
}

func (s stream[T]) TakeWhile(pred func(T) bool) stream[T] {
	return stream[T]{
		iterator: xiter.TakeWhile(s.iterator, pred),
		sizeInfo: s.sizeInfo.CopyWithUnknownSize(),
	}
}

func (s stream[T]) Skip(n int) stream[T] {
	return stream[T]{
		iterator: xiter.Skip(s.iterator, n),
		sizeInfo: s.sizeInfo.CopyWithSize(max(0, s.Size()-n)),
	}
}

func (s stream[T]) SkipWhile(pred func(T) bool) stream[T] {
	return stream[T]{
		iterator: xiter.SkipWhile(s.iterator, pred),
		sizeInfo: s.sizeInfo.CopyWithUnknownSize(),
	}
}

func (s stream[T]) Transform(transformer func(T) T) stream[T] {
	return stream[T]{
		iterator: xiter.Map(s.iterator, transformer),
		sizeInfo: s.sizeInfo,
	}
}

func (s stream[T]) Map[V any](f func(T) V) stream[V] {
	return stream[V]{
		iterator: xiter.Map(s.iterator, f),
		sizeInfo: s.sizeInfo,
	}
}

func (s stream[T]) FlatMap[V any](transformer func(T) iter.Seq[V]) stream[V] {
	return stream[V]{
		iterator: xiter.FlatMap(s.iterator, transformer),
		// Это MinSize
		sizeInfo: s.sizeInfo,
	}
}

func (s stream[T]) FlatTransform(transformer func(T) iter.Seq[T]) stream[T] {
	return stream[T]{
		iterator: xiter.FlatMap(s.iterator, transformer),
		// Это MinSize
		sizeInfo: s.sizeInfo,
	}
}

func (s stream[T]) OrderBy(cmp func(T, T) int) stream[T] {
	items := slices.SortedFunc(s.iterator, cmp)
	s.iterator = slices.Values(items)
	return s
}

func (s stream[T]) Reverse() stream[T] {
	items := s.ToSlice()
	slices.Reverse(items)
	s.iterator = slices.Values(items)
	return s
}

func (s stream[T]) UniqueByKeyPrior[V comparable](mapper func(T) V, cmp func(T, T) int) stream[T] {
	return stream[T]{
		iterator: xiter.UniqueByKeyPrior(
			s.iterator,
			mapper,
			cmp,
			s.sizeInfo.GetPreAllocSize(),
		),
		sizeInfo: s.sizeInfo.CopyWithUnknownSize(),
	}
}

func (s stream[T]) UniqueByKey[V comparable](mapper func(T) V) stream[T] {
	return stream[T]{
		iterator: xiter.UniqueByKey(
			s.iterator,
			mapper,
			s.sizeInfo.GetPreAllocSize(),
		),
		sizeInfo: s.sizeInfo.CopyWithUnknownSize(),
	}
}

// Если T не comparable, то получим runtime panic
func (s stream[T]) Unique() stream[T] {
	return stream[T]{
		iterator: xiter.UniqueByKey(
			s.iterator,
			func(v T) any { return v },
			s.sizeInfo.GetPreAllocSize(),
		),
		sizeInfo: s.sizeInfo.CopyWithUnknownSize(),
	}
}
