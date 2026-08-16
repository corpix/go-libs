package awaitable

import (
	"context"
	"errors"
	"iter"
	"maps"
	"slices"
	"sync"

	"github.com/SlamJam/go-libs/xchan"
	"github.com/SlamJam/go-libs/xiter"
)

// Последовательность Awaitableэ'ов, ассоциированных с ключами
type AwaitableSeq[K comparable] iter.Seq2[K, Awaitable]

func FromSlice[T Awaitable, S ~[]T](s S) AwaitableSeq[int] {
	return AwaitableSeq[int](xiter.Map2(slices.All(s), func(_ int, item T) Awaitable { return item }))
}

func FromItems(s ...Awaitable) AwaitableSeq[int] {
	return FromSlice(s)
}

func FromMap[K comparable](m map[K]Awaitable) AwaitableSeq[K] {
	return AwaitableSeq[K](maps.All(m))
}

// TODO: rename IterAwaitResults
func (as AwaitableSeq[K]) AwaitResults(ctx context.Context) ResultSeq[K] {
	type errorWithKey struct {
		Key K
		Err error
	}

	return func(yield func(K, error) bool) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		errCh := make(chan errorWithKey)
		defer xchan.Drain(errCh)

		var wg sync.WaitGroup
		for key, item := range as {
			wg.Go(func() {
				errCh <- errorWithKey{Key: key, Err: item.Await(ctx)}
			})
		}

		go func() {
			defer close(errCh)
			wg.Wait()
		}()

		for item := range errCh {
			if !yield(item.Key, item.Err) {
				return
			}
		}
	}
}

func (as AwaitableSeq[K]) AwaitResultsSkipPending(rootCtx context.Context) ResultSeq[K] {
	return func(yield func(K, error) bool) {
		sentinetErr := errors.New("sentinel")

		ctx, cancel := context.WithCancelCause(context.TODO())
		defer cancel(sentinetErr)

		done := make(chan struct{})
		defer close(done)

		go func() {
			select {
			case <-rootCtx.Done():
				cancel(sentinetErr)
			case <-done:
			}
		}()

		for key, err := range as.AwaitResults(ctx) {
			if errors.Is(err, sentinetErr) {
				continue
			}

			if !yield(key, err) {
				return
			}
		}
	}
}

func (as AwaitableSeq[K]) NowReadyResults() ResultSeq[K] {
	ctx, cancel := context.WithCancel(context.TODO())
	cancel()

	return as.AwaitResultsSkipPending(ctx)
}
