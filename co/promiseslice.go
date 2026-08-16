package co

import (
	"context"
	"slices"

	"github.com/SlamJam/go-libs/co/awaitable"
	"github.com/SlamJam/go-libs/co/promise"
)

// var ErrEmptyPromiseSlice = errors.New("empty PromiseSlice")

type PromiseSlice[T any] []promise.Promise[T]

func (mp *PromiseSlice[T]) Add(f func() (T, error)) {
	mp.Append(S.Launch(f))
}

func (mp *PromiseSlice[T]) Append(p promise.Promise[T]) {
	*mp = append(*mp, p)
}

func (mp PromiseSlice[T]) AsAwaitables() awaitable.AwaitableSeq[int] {
	// return xslices.Map(mp, func(p Promise[T]) Awaitable { return p })
	return awaitable.FromSlice(mp)
}

var _ awaitable.Awaitable = PromiseSlice[int]{}

func (mp PromiseSlice[T]) Await(ctx context.Context) error {
	return mp.AsAwaitables().AwaitResults(ctx).Collect().Err()
}

func (mp PromiseSlice[T]) Results(ctx context.Context) ResultSeq[int, T] {
	return AwaitSeq[int, T](ctx, slices.All(mp))
}

// Альтернатива:
// ctx, cancel := context.WithCancel(context.TODO())
// cancel()
// mp.Results(ctx).SkipPending().CollectAll2()
func (mp PromiseSlice[T]) ReadyResults() (map[int]T, map[int]error) {
	return ReadySeq[int, T](slices.All(mp)).CollectAll()
}

// AwaitAll дожидается выполнения всех задач
// Возвращает или все результаты или все ошибки, собранные в multierr
// Например, пишем в несколько шардов и хотим получить ВСЕ возникшие ошибки
func (mp PromiseSlice[T]) AwaitAll(ctx context.Context) ([]T, map[int]error) {
	results := make([]T, len(mp))
	resultErr := make(map[int]error, len(mp))

	for idx, res := range mp.Results(ctx) {
		if val, err := res.Unwrap(); err != nil {
			resultErr[idx] = err
		} else {
			results[idx] = val
		}
	}

	if len(resultErr) != 0 {
		return nil, resultErr
	}

	return results, nil
}

// AwaitAllOrFirstError дожидается выполнения всех задач или первой возникшей ошибки
// Возвращает или все результаты или первую возникшую ошибку
// Реализует концепт ИЛИ: "прочитать  реплики"
func (mp PromiseSlice[T]) AwaitAllOrFirstError(ctx context.Context) ([]T, error) {
	results := make([]T, len(mp))

	for idx, res := range mp.Results(ctx) {
		if val, err := res.Unwrap(); err != nil {
			return nil, err
		} else {
			results[idx] = val
		}
	}

	return results, nil
}
