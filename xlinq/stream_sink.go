package xlinq

import (
	"iter"
	"maps"
	"slices"

	std "github.com/SlamJam/go-libs"
	"github.com/SlamJam/go-libs/xiter"
)

func (s stream[V]) Iter() iter.Seq[V] {
	return s.iterator
}

func (s stream[V]) Enumerate() iter.Seq2[int, V] {
	return xiter.Enumerate(s.iterator)
}

func (s stream[V]) ToSlice() []V {
	res := make([]V, 0, s.sizeInfo.GetPreAllocSize())

	return slices.AppendSeq(res, s.iterator)
}

func (s stream[V]) ForEach(f func(V)) {
	for v := range s.iterator {
		f(v)
	}
}

func (s stream[V]) ForEachUntil(f func(V) bool) {
	for v := range s.iterator {
		if !f(v) {
			return
		}
	}
}

func (s stream[V]) UntilFirstErr(f func(V) error) (V, error) {
	for v := range s.iterator {
		if err := f(v); err != nil {
			return v, err
		}
	}

	return std.ZeroErr[V](nil)
}

func (s stream[V]) MapToErr(f func(V) error) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		for v := range s.iterator {
			if !yield(v, f(v)) {
				return
			}
		}
	}
}

func (s stream[V]) ToGroups(mapper func(V) any, capacity int) iter.Seq[[]V] {
	groups := xiter.GroupByKey(s.iterator, mapper, capacity)
	return maps.Values(groups)
}

func _() {
	// ss := FromSlice([]int{1, 2, 3, 5})
	ss := FromItems(1, 2, 3, 5)

	check := func(int) error { return nil }

	for item := range ss.Iter() {
		if err := check(item); err != nil {

		}
	}

	for item, err := range ss.MapToErr(check) {
		if err != nil {

		}
		_ = item
	}
}
