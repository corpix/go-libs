package promise

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"

	std "github.com/SlamJam/go-libs"
	"github.com/SlamJam/go-libs/co/awaitable"
	"github.com/SlamJam/go-libs/types/option"
	"github.com/SlamJam/go-libs/types/result"
	"github.com/SlamJam/go-libs/xerrors"
	"github.com/pkg/errors"
)

var ErrAwaitCanceled = errors.New("promise await canceled")

type Promise[T any] = *promise[T]

var _ awaitable.Awaitable = &promise[any]{}

type promise[T any] struct {
	once sync.Once

	fulfilled atomic.Bool
	done      chan struct{}

	result T
	err    error

	createdStack []uintptr // raw PCs, форматируется лениво
}

func (p *promise[T]) assertInitialized() {
	if p == nil || p.done == nil {
		panic("misuse: promise must be created via constructors")
	}
}

func (p *promise[T]) AsAwaitable() awaitable.Awaitable {
	return p
}

func (p *promise[T]) IsFulfilled() bool {
	return p.fulfilled.Load()
}

// Fulfilled - мост для оператора select. Аналогично context.Context.Done()
func (p *promise[T]) Fulfilled() <-chan struct{} {
	return p.done
}

func (p *promise[T]) fulfill(result T, err error) (fulfilled bool) {
	p.once.Do(func() {
		if err != nil {
			p.err = errors.WithStack(err)
		} else {
			p.result = result
		}

		p.fulfilled.Store(true)
		close(p.done)

		fulfilled = true
	})

	return
}

func (p *promise[T]) Await(ctx context.Context) error {
	_, err := p.GetOrAwait(ctx)
	return err
}

func (p *promise[T]) GetOrAwait(ctx context.Context) (T, error) {
	p.assertInitialized()

	// Perf optimization: Fast path if already completed
	if p.IsFulfilled() {
		return p.result, p.err
	}

	select {
	case <-p.done:
		return p.result, p.err

	case <-ctx.Done():
		return std.Zero[T](), xerrors.WithCauseFromContext(ErrAwaitCanceled, ctx)
	}
}

func (p *promise[T]) GetOrAwaitResult(ctx context.Context) result.Result[T] {
	return result.New(p.GetOrAwait(ctx))
}

func (p *promise[T]) OptionalResult() option.Option[result.Result[T]] {
	if !p.IsFulfilled() {
		return option.Empty[result.Result[T]]()
	}

	return option.Just(result.New(p.result, p.err))
}

func (p *promise[T]) TryGet(ctx context.Context) (T, error, bool) {
	p.assertInitialized()

	if p.IsFulfilled() {
		return p.result, p.err, true
	}

	return std.Zero[T](), nil, false
}

func newPromise[T any]() Promise[T] {
	var pcs [16]uintptr
	n := runtime.Callers(3, pcs[:])

	return &promise[T]{
		done:         make(chan struct{}),
		createdStack: pcs[:n],
	}
}

// NewResolved ...
func NewResolved[T any](result T) Promise[T] {
	p := newPromise[T]()
	p.fulfill(result, nil)

	return p
}

func NewRejected[T any](err error) Promise[T] {
	p := newPromise[T]()

	var t T
	p.fulfill(t, err)

	return p
}

// Любой из этих методов роняет компилятор

func (p *promise[T]) Map[V any](ctx context.Context, fn func(T) (V, error)) Promise[V] {
	return nil
}

func (p *promise[T]) FlatMap[V any](ctx context.Context, fn func(T) Promise[V]) Promise[V] {
	return nil
}
