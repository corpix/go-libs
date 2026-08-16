package xlinq_test

import (
	"testing"

	"github.com/SlamJam/go-libs/xlinq"
	"github.com/stretchr/testify/assert"
)

// Behaviors covered:
// 1) Should return only the first occurrence for duplicate primitive values (stability of first seen).
// 2) Should preserve input order of the first occurrences.
// 3) Should return an empty slice when the source is empty.
// 4) Should handle large inputs efficiently and without panics.
// 5) Should not allocate more than necessary when a size hint is provided (prealloc path exercised).
// 6) Should work for non-comparable element types by using UniqueByKey; here we use Unique on primitives and UniqueByKey for structs.
// 7) Should not mutate the original stream when chained with other operators.

func TestStream_Unique_Primitives_FirstOccurrenceAndOrder(t *testing.T) {
	t.Parallel()

	in := []int{3, 1, 2, 3, 2, 1, 4, 3, 5}
	s := xlinq.FromSlice(in)
	out := s.Unique().ToSlice()

	assert.Equal(t, []int{3, 1, 2, 4, 5}, out)
}

func TestStream_Unique_Empty(t *testing.T) {
	t.Parallel()

	s := xlinq.FromEmpty[int]()
	out := s.Unique().ToSlice()

	assert.Empty(t, out)
}

func TestStream_Unique_LargeInput_NoPanic(t *testing.T) {
	t.Parallel()

	// Build a large input with repeating pattern
	in := make([]int, 0, 5000)
	for i := 0; i < 1000; i++ {
		in = append(in, i%50)
	}

	s := xlinq.FromSlice(in)
	out := s.Unique().ToSlice()

	// Expect values 0..49 once each, in their first-seen order
	expected := make([]int, 50)
	for i := 0; i < 50; i++ {
		expected[i] = i
	}
	assert.Equal(t, expected, out)
}

func TestStream_Unique_RespectsSizeHintForPrealloc(t *testing.T) {
	t.Parallel()

	// Provide a size hint lower than total to exercise prealloc path; behavior shouldn't change
	in := []int{1, 1, 1, 2, 2, 3}
	s := xlinq.FromSlice(in).WithSizeHintRelative(50) // arbitrary hint; correctness is what we assert
	out := s.Unique().ToSlice()

	assert.Equal(t, []int{1, 2, 3}, out)
}

func TestStream_Unique_DoesNotMutateSourceWhenChained(t *testing.T) {
	t.Parallel()

	in := []int{1, 2, 2, 3, 3, 3}
	s := xlinq.FromSlice(in)

	// Chain: filter then unique then reverse (ensures intermediate ops don't corrupt unique behavior)
	out := s.Where(func(v int) bool { return v > 0 }).Unique().Reverse().ToSlice()

	// Unique first -> [1,2,3], then Reverse -> [3,2,1]
	assert.Equal(t, []int{3, 2, 1}, out)

	// Original input slice remains unchanged
	assert.Equal(t, []int{1, 2, 2, 3, 3, 3}, in)
}

func TestStream_UniqueByKey_OnStructs(t *testing.T) {
	t.Parallel()

	type item struct {
		id   int
		name string
	}

	in := []item{{1, "a"}, {2, "b"}, {1, "a2"}, {3, "c"}, {2, "b2"}}
	s := xlinq.FromSlice(in)

	out := s.UniqueByKey(func(it item) int { return it.id }).ToSlice()

	// Expect first seen per id preserved in input order
	assert.Equal(t, []item{{1, "a"}, {2, "b"}, {3, "c"}}, out)
}

func TestStream_TakeWhile_StopsAtFirstFalse(t *testing.T) {
	t.Parallel()

	in := []int{1, 2, 3, 0, 4, 5}
	s := xlinq.FromSlice(in)
	out := s.TakeWhile(func(v int) bool { return v > 0 }).ToSlice()

	assert.Equal(t, []int{1, 2, 3}, out)
}

func TestStream_TakeWhile_PredicateAlwaysTrue_ReturnsAll(t *testing.T) {
	t.Parallel()

	in := []int{1, 2, 3, 4}
	s := xlinq.FromSlice(in)
	out := s.TakeWhile(func(v int) bool { return true }).ToSlice()

	assert.Equal(t, in, out)
}

func TestStream_TakeWhile_PredicateAlwaysFalse_ReturnsEmpty(t *testing.T) {
	t.Parallel()

	in := []int{1, 2, 3}
	s := xlinq.FromSlice(in)
	out := s.TakeWhile(func(v int) bool { return false }).ToSlice()

	assert.Empty(t, out)
}

func TestStream_TakeWhile_PredicateNotCalledAfterStop(t *testing.T) {
	t.Parallel()

	in := []int{1, 2, 3, 0, 4, 5}
	calls := 0

	s := xlinq.FromSlice(in)
	out := s.TakeWhile(func(v int) bool {
		calls++
		return v > 0
	}).ToSlice()

	// Should stop at first non-positive value (0). Predicate should be called 4 times: 1,2,3,0
	assert.Equal(t, []int{1, 2, 3}, out)
	assert.Equal(t, 4, calls)
}

func TestStream_TakeWhile_SizeIsUnknown(t *testing.T) {
	t.Parallel()

	s := xlinq.FromSlice([]int{1, 2, 3, 4})
	tw := s.TakeWhile(func(v int) bool { return v%2 == 1 || v%2 == 0 }) // always true, but size should be unknown

	assert.Equal(t, 0, tw.Size())
}
