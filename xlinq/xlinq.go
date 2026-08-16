package xlinq

import (
	"iter"
)

type Iterable[V any] interface {
	Iter() iter.Seq[V]
}

type Sizable interface {
	Size() int
}

// Stream асбтакция потока элеметов типа V
type Stream[V any] interface {
	Iterable[V]
	Sizable

	Where(func(V) bool) Stream[V]
	Take(int) Stream[V]
	TakeWhile(func(V) bool) Stream[V]
	Skip(int) Stream[V]
	SkipWhile(func(V) bool) Stream[V]

	Transform(func(V) V) Stream[V]
	FlatTransform(func(V) iter.Seq[V]) Stream[V]
	// SelectWithIndex

	// Map[T any](func(V) V) Stream[V]

	// Ordering
	OrderBy(cmp func(V, V) int) Stream[V]
	Reverse() Stream[V]

	UniqueByKeyPrior(mapper func(V) any, cmp func(V, V) int) Stream[V]
	UniqueByKey(mapper func(V) any) Stream[V]
	Unique() Stream[V]

	// Grouping: GroupBy

	// Combining: Zip - combine two sequences using result selector

	// Terminal: ToSlice, First, Last, ElementAt, ElementAtOrDefault, Contains, ContainsBy, Count, Any, AnyMatch, All, Aggregate, ForEach
	Iter() iter.Seq[V]
	Enumerate() iter.Seq2[int, V]
	ToSlice() []V
	ForEach(f func(V))
	ForEachUntil(f func(V) bool)
	UntilFirstErr(f func(V) error) (V, error)
	MapToErr(f func(V) error) iter.Seq2[V, error]
	ToGroups(mapper func(V) any, capacity int) iter.Seq[[]V]

	// Size Information
	WithSizeHintExact(sizeHint int) Stream[V]
	WithSizeHintRelative(ratio int) Stream[V]
	WithSizeHint(f SizeHint) Stream[V]
}

type stream[V any] struct {
	iterator iter.Seq[V]
	sizeInfo sizeInfo
}

func (s stream[V]) Size() int {
	return s.sizeInfo.GetSize()
}

// ???
func (s stream[V]) WithSizeHintExact(sizeHint int) stream[V] {
	s.sizeInfo.SetSizeHint(SizeHintExact(sizeHint))
	return s
}

// ???
func (s stream[V]) WithSizeHintRelative(ratio int) stream[V] {
	s.sizeInfo.SetSizeHint(SizeHintRelative(ratio))
	return s
}

func (s stream[V]) WithSizeHint(f SizeHint) stream[V] {
	s.sizeInfo.SetSizeHint(f)
	return s
}

// var _ Stream[int] = stream[int]{}

/*
// StreamKV асбтакция потока пар K и V типа, где K является сравнимым.
// Такая пара олицетворяет собой ключ и соответствющее значение из словаря
type StreamKV[K comparable, V any] interface {
	Stream[xiter.KV[K, V]]

	ToMap() map[K]V
}

type streamKV[K comparable, V any] struct {
	stream[xiter.KV[K, V]]
}

var _ StreamKV[int, int] = streamKV[int, int]{}

// StreamComparable поток сравнимых значений типа V
type StreamComparable[V comparable] interface {
	Stream[V]

	Unique(int) StreamComparable[V]
}

type streamComparable[V comparable] struct {
	stream[V]
}

var _ StreamComparable[int] = streamComparable[int]{}
*/

// func foo() {
// 	s := Comparable(FromEmpty[int]())
// 	s.
// }
