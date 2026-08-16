package xsingleflight

import "golang.org/x/sync/singleflight"

type Group struct{ g singleflight.Group }

type Result[T any] struct {
	Val    T
	Err    error
	Shared bool
}

func (g *Group) Do[T any](key string, fn func() (T, error)) (v T, err error, shared bool) {
	val, err, shared := g.g.Do(key, func() (any, error) { return fn() })
	return val.(T), err, shared
}

func (g *Group) DoChan[T any](key string, fn func() (T, error)) <-chan Result[T] {
	result := make(chan Result[T], 1)

	go func() {
		// defer close(result)

		r := <-g.g.DoChan(key, func() (any, error) { return fn() })
		result <- Result[T]{Val: r.Val.(T), Err: r.Err, Shared: r.Shared}
	}()

	return result
}

func (g *Group) Forget(key string) {
	g.g.Forget(key)
}
