package xerrgroup

import (
	"errors"

	"github.com/SlamJam/go-libs/co/promise"
	"golang.org/x/sync/errgroup"
)

var ErrNotLaunched = errors.New("promise was not launched")

type Group struct {
	g errgroup.Group
}

func (g *Group) Wait() error {
	return g.g.Wait()
}

func (g *Group) Go(fn func() error) {
	g.g.Go(fn)
}

func (g *Group) GoDeferred[T any](fn func() (T, error)) promise.Promise[T] {
	p, r := promise.NewDeferred[T]()

	g.g.Go(func() error {
		r.SafetyFulfill(fn)
		return nil
	})

	return p
}

func (g *Group) TryGo(f func() error) bool {
	return g.g.TryGo(f)
}

func (g *Group) TryGoDeferred[T any](fn func() (T, error)) (promise.Promise[T], bool) {
	p, r := promise.NewDeferred[T]()

	launched := g.g.TryGo(func() error {
		r.SafetyFulfill(fn)
		return nil
	})

	if !launched {
		r.Reject(ErrNotLaunched)
	}

	return p, launched
}

func (g *Group) SetLimit(n int) {
	g.g.SetLimit(n)
}

// func _() {
// 	var g errgroup.Group
// 	_ = g
// }
