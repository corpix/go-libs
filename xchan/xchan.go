package xchan

import (
	"context"
	"sync"

	"github.com/SlamJam/go-libs/xerrors"
)

func FanIn[T, V any](fn func(merged <-chan T, cancel func()) (V, error), chans ...<-chan T) (V, error) {
	result := make(chan T)

	readCtx, readCancel := context.WithCancel(context.Background())
	defer readCancel()

	writeCtx, writeCancel := context.WithCancel(context.Background())
	defer writeCancel()

	wg := sync.WaitGroup{}
	for _, ch := range chans {
		wg.Go(func() {
			for {
				val, ok := NewReader(ch).ReadValueWithContext(readCtx)
				if !ok {
					return
				}

				err := NewWriter(result).Put(writeCtx, val)
				if err != nil {
					return
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return fn(result, readCancel)
}

func Drain[T any](ch <-chan T) {
	for range ch {
	}
}

func DrainAsync[T any](ctx context.Context, ch <-chan T) <-chan error {
	res := make(chan error, 1)

	go func() {
		defer close(res)

		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
			case <-ctx.Done():
				res <- xerrors.FromContext(ctx)
			}
		}
	}()

	return res
}
