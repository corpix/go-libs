package co

import (
	"context"
	"sync"
	"sync/atomic"

	std "github.com/SlamJam/go-libs"
	"github.com/pkg/errors"
)

// type Deffered[T any] func() (T, error)

type Promise[T any] = *promise[T]

var _ Awaitable = NewPromise(func() (std.Void, error) { return std.Void{}, nil })

type promise[T any] struct {
	f      func() (T, error)
	result T
	err    error

	lazy        bool
	initialized bool
	launched    atomic.Bool
	comleted    atomic.Bool
	done        chan struct{}

	onceLaunch sync.Once
}

func (p *promise[T]) assertInitialized() {
	if !p.initialized {
		panic("misuse: Promise must be create with NewPromise")
	}
}

func (p *promise[T]) IsLaunched() bool {
	return p.launched.Load()
}

func (p *promise[T]) IsCompleted() bool {
	return p.comleted.Load()
}

func (p *promise[T]) onComplete(result T, err error) {
	p.result = result
	p.err = err
	p.comleted.Store(true)
	close(p.done)
}

func (p *promise[T]) launch() {
	p.launched.Store(true)
	p.result, p.err = p.f()
	p.comleted.Store(true)
	close(p.done)
}

func (p *promise[T]) ensureLaunched() (result bool) {
	if p.launched.Load() {
		return
	}

	p.onceLaunch.Do(func() {
		result = true
		p.launched.Store(true)
		go p.launch()
	})

	return
}

func (p *promise[T]) Poll(ctx context.Context) (T, error) {
	p.assertInitialized()

	if p.IsCompleted() {
		return p.result, p.err
	}

	p.ensureLaunched()

	select {
	case <-p.done:
		return p.result, p.err

	case <-ctx.Done():
		return std.Zero[T](), ctx.Err()
	}
}

func (p *promise[T]) Await(ctx context.Context) (err error) {
	_, err = p.Poll(ctx)
	return
}

func (p *promise[T]) Value() (result T, err error) {
	if !p.IsCompleted() {
		panic("trying to read a value that is not set")
	}

	return p.result, p.err
}

func newPromise[T any](f func() (T, error), lazy bool) Promise[T] {
	p := &promise[T]{
		f:    f,
		lazy: lazy,

		initialized: true,
		done:        make(chan struct{}),
	}

	if !p.lazy {
		p.ensureLaunched()
	}

	return p
}

func NewPromise[T any](f func() (T, error)) Promise[T] {
	return newPromise(f, false)
}

func NewLazyPromise[T any](f func() (T, error)) Promise[T] {
	return newPromise(f, true)
}

func NewResolved[T any](result T) Promise[T] {
	p := NewLazyPromise(func() (T, error) {
		return result, nil
	})
	p.launch()

	return p
}

func NewRejected[T any](err error) Promise[T] {
	p := NewLazyPromise(func() (T, error) {
		return std.Zero[T](), err
	})
	p.launch()

	return p
}

var Resolved = NewResolved(std.Void{})

var ErrRejected = errors.New("rejected")

var Rejected = NewRejected[std.Void](ErrRejected)

func NewNeverResolved[T any]() Promise[T] {
	return NewPromise(func() (T, error) {
		<-context.Background().Done()
		return std.Zero[T](), nil
	})
}

// var Never = NewPromise(func() (struct{}, error) {
// 	<-context.Background().Done()
// 	return struct{}{}, nil
// })
