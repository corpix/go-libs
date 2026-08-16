package scope

import (
	"context"
	"sync"
	"time"

	std "github.com/SlamJam/go-libs"

	"github.com/SlamJam/go-libs/co/promise"
	"github.com/SlamJam/go-libs/xerrors"
)

var S = &Scope{}

type Scope struct {
	ctx context.Context

	wg sync.WaitGroup

	panics chan error
}

type ScopeResult struct {
	scope *Scope

	onceWaiter sync.Once

	done chan struct{}
}

func (sr *ScopeResult) Await(ctx context.Context) error {
	sr.onceWaiter.Do(func() {
		sr.scope.wg.Wait()
		sr.done = make(chan struct{})
		close(sr.done)
	})

	select {
	case <-sr.done:
	case <-ctx.Done():
		return xerrors.FromContext(ctx)
	}

	return nil
}

func (s *Scope) Ctx() context.Context {
	return s.ctx
}

func (s *Scope) EnablePanicsChan() {
	if s.panics != nil {
		return
	}

	s.panics = make(chan error)
}

func (s *Scope) Panics() <-chan error {
	return s.panics
}

func (s *Scope) publishPanic(err error) {
	if err != nil || s.panics == nil {
		return
	}

	s.panics <- err
}

func (s *Scope) Go(fn func()) {
	s.wg.Go(func() {
		err := xerrors.RecoverVoid(func() error { fn(); return nil })
		s.publishPanic(err)
	})
}

func (s *Scope) Launch[T any](f func() (T, error)) promise.Promise[T] {
	p, r := promise.NewDeferred[T]()
	s.wg.Go(func() {
		r.SafetyFulfill(f)
	})

	return p
}

func (s *Scope) LaunchAt[T any](ctx context.Context, f func() (T, error), t time.Time) promise.Promise[T] {
	return s.LaunchAfter(ctx, f, time.Until(t))
}

func (s *Scope) LaunchAfter[T any](ctx context.Context, f func() (T, error), d time.Duration) promise.Promise[T] {
	p, r := promise.NewDeferred[T]()

	s.wg.Go(func() {
		select {
		case <-time.After(d):
			r.SafetyFulfill(f)

		// В процессе ожидания времени запуска отменился контекст
		case <-ctx.Done():
			r.Reject(xerrors.FromContext(ctx))

		case <-s.ctx.Done():
			r.Reject(xerrors.FromContext(s.ctx))
		}
	})

	return p
}

func (s *Scope) MapPromise[T, V any](p promise.Promise[T], fn func(T) (V, error)) promise.Promise[V] {
	return s.Launch(func() (V, error) {
		return p.GetOrAwaitResult(s.Ctx()).
			Map(fn).
			Unwrap()
	})
}

func (s *Scope) FlatMapPromise[T, V any](p promise.Promise[T], fn func(T) promise.Promise[V]) promise.Promise[V] {
	return s.Launch(func() (V, error) {
		val, err := p.GetOrAwait(s.Ctx())
		if err != nil {
			return std.Zero[V](), err
		}

		return fn(val).GetOrAwait(s.Ctx())
	})
}

type contextKey string

const scopeKey contextKey = "scope"

func FromContext(ctx context.Context) *Scope {
	s := ctx.Value(scopeKey).(*Scope)
	if s != nil {
		return s
	}

	return S
}

func WithScope[T any](rootCtx context.Context, f func(*Scope) (T, error)) (T, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	s := Scope{}
	s.ctx = context.WithValue(ctx, scopeKey, &s)

	return f(&s)
}

func _() {
	_, _ = WithScope(context.TODO(), func(s *Scope) (int, error) {
		p := s.Launch(func() (int, error) {
			s.Ctx()
			return 5, nil
		})

		return p.GetOrAwait(s.Ctx())
	})
}
