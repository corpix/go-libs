package xchan

import (
	"context"
	"time"

	"github.com/SlamJam/go-libs/xiter"
)

type Reader[T any] <-chan T

func NewReader[T any](ch <-chan T) Reader[T] {
	return ch
}

func (ch Reader[T]) ReadValue(stop <-chan struct{}) (T, bool) {
	select {
	case value, ok := <-ch:
		return value, ok
	case <-stop:
		var t T
		return t, false
	}
}

func (ch Reader[T]) ReadValueWithContext(ctx context.Context) (T, bool) {
	return ch.ReadValue(ctx.Done())
}

// func (ch Reader[T]) ReadValueWithTimeout(d time.Duration) (T, bool) {
// 	return ch.ReadValue()
// }

func (ch Reader[T]) TakeUntil(stop <-chan struct{}) xiter.XSeq[T] {
	return func(yield func(T) bool) {
		for {
			select {
			case item, ok := <-ch:
				if !ok || !yield(item) {
					return
				}
			case <-stop:
				return
			}
		}
	}
}

func (ch Reader[T]) TakeUntilContext[V any](ctx context.Context) xiter.XSeq[T] {
	return ch.TakeUntil(ctx.Done())
}

func (ch Reader[T]) batch(n int, until <-chan struct{}, flusher <-chan time.Time) xiter.XSeq[[]T] {
	return func(yield func([]T) bool) {
		batch := make([]T, 0, n)

		flush := func() bool {
			if len(batch) == 0 {
				return true
			}

			duplicated := make([]T, len(batch))
			copy(duplicated, batch)
			batch = batch[:]

			return yield(duplicated)
		}
		defer flush()

		for {
			select {
			case item, ok := <-ch:
				if !ok {
					return
				}

				batch = append(batch, item)
				if len(batch) >= n {
					if !flush() {
						return
					}
				}
			case <-flusher:
				if !flush() {
					return
				}
			case <-until:
				return
			}
		}
	}
}

type batchOptions struct {
	until   <-chan struct{}
	flusher <-chan time.Time
}

func WithPeriodicFlush(d time.Duration) func(opts *batchOptions) {
	return func(opts *batchOptions) {
		t := time.NewTicker(d)
		opts.flusher = t.C
	}
}
func WithFlushOnTimer(t *time.Timer) func(opts *batchOptions) {
	return func(opts *batchOptions) {
		opts.flusher = t.C
	}
}
func WithStop(ch <-chan struct{}) func(opts *batchOptions) {
	return func(opts *batchOptions) {
		opts.until = ch
	}
}
func WithStopContext(ctx context.Context) func(opts *batchOptions) {
	return func(opts *batchOptions) {
		opts.until = ctx.Done()
	}
}

func (ch Reader[T]) Batch(n int, opts ...func(opts *batchOptions)) xiter.XSeq[[]T] {
	var options batchOptions
	for _, opt := range opts {
		opt(&options)
	}

	return ch.batch(n, options.until, options.flusher)
}

func _() {
	var ch chan int

	t := time.NewTimer(5 * time.Second)
	for item := range NewReader(ch).Batch(10, WithFlushOnTimer(t)) {
		// process
		_ = item
		t.Reset(5 * time.Second)
	}
}
