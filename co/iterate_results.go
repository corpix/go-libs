package co

import (
	"context"
	"iter"
	"sync"

	"github.com/SlamJam/go-libs/co/promise"
	"github.com/SlamJam/go-libs/types/result"
	"github.com/SlamJam/go-libs/xchan"
)

// Нужно рассмотреть Seq2[K, Promise[T]] -> Seq2[K, Result[T]]
// При ожидании результата промиса с ограниченным временем у нас есть три исхода: T, error и "промис не разрезолвился"
//

type resultWithKey[K comparable, T any] struct {
	Key    K
	Result result.Result[T]
}

func AwaitSeq[K comparable, T any](
	ctx context.Context,
	items iter.Seq2[K, promise.Promise[T]],
) ResultSeq[K, T] {
	return func(yield func(K, result.Result[T]) bool) {
		ch := make(chan resultWithKey[K, T])
		defer xchan.Drain(ch)

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var wg sync.WaitGroup
		for key, p := range items {
			wg.Go(func() {
				ch <- resultWithKey[K, T]{
					Key:    key,
					Result: p.GetOrAwaitResult(ctx),
				}
			})
		}

		go func() {
			defer close(ch)
			wg.Wait()
		}()

		for item := range ch {
			if !yield(item.Key, item.Result) {
				// На всякий случай
				cancel()
				return
			}
		}
	}
}

func ReadySeq[K comparable, T any](
	items iter.Seq2[K, promise.Promise[T]],
) ResultSeq[K, T] {
	return func(yield func(K, result.Result[T]) bool) {
		for key, item := range items {
			optRes := item.OptionalResult()
			if val, ok := optRes.Unwrap(); ok {
				if !yield(key, val) {
					return
				}
			}
		}
	}
}

// func AwaitPartial[K comparable, T any](
// 	rootCtx context.Context,
// 	items iter.Seq2[K, Promise[T]],
// ) ResultSeq[K, T] {
// 	return func(yield func(K, result.Result[T]) bool) {
// 		var errSentinelForPending = errors.New("promise is pending")

// 		var done chan struct{}
// 		defer close(done)

// 		ctx, cancel := context.WithCancelCause(rootCtx)

// 		go func() {
// 			select {
// 			case <-rootCtx.Done():
// 				cancel(errSentinelForPending)
// 			case <-done:
// 			}
// 		}()

// 		for key, result := range AwaitSeq(ctx, items) {
// 			if errors.Is(result.Err(), errSentinelForPending) {
// 				continue
// 			}

// 			if !yield(key, result) {
// 				return
// 			}
// 		}
// 	}
// }
