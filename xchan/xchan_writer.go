package xchan

import (
	"context"
	"time"

	std "github.com/SlamJam/go-libs"
)

type Writer[T any] chan<- T

func NewWriter[T any](ch chan<- T) Writer[T] {
	return ch
}

func (ch Writer[T]) Put(ctx context.Context, value T) error {
	select {
	case ch <- value:
	case <-ctx.Done():
		return context.Cause(ctx)
	}

	return nil
}

func (ch Writer[T]) PutWithTimeout(value T, timeout time.Duration) error {
	select {
	case ch <- value:
	case <-time.After(timeout):
		return std.ErrTimeout
	}

	return nil
}

// TryPut пытается отправить значение в канал без блокировки.
// Возвращает true, если значение было отправлено, и false, если канал переполнен.
func (ch Writer[T]) TryPut(value T) bool {
	select {
	case ch <- value:
		return true
	default:
		return false
	}
}
