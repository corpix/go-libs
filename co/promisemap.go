package co

import (
	"context"
	"slices"

	"github.com/SlamJam/go-libs/co/awaitable"
	"github.com/SlamJam/go-libs/co/promise"
	"github.com/SlamJam/go-libs/xiter"
	"github.com/SlamJam/go-libs/xsync"
	"github.com/pkg/errors"
)

var ErrEmptyPromiseMap = errors.New("empty PromiseMap")

type PromiseMap[K comparable, T any] struct {
	m xsync.Map[K, promise.Promise[T]]
}

func (pm *PromiseMap[K, T]) Append(key K, p promise.Promise[T]) {
	pm.m.Store(key, p)
}

func (pm *PromiseMap[K, T]) AsAwaitables() []awaitable.Awaitable {
	return slices.Collect(xiter.Map(pm.m.Values, promise.Promise[T].AsAwaitable))
}

func (pm *PromiseMap[K, T]) AsAwaitable(ctx context.Context) awaitable.Awaitable {
	return pm
}

func (pm *PromiseMap[K, T]) Await(ctx context.Context) error {
	return awaitable.FromSlice(pm.AsAwaitables()).AwaitResults(ctx).Collect().Err()
}

func (pm *PromiseMap[K, T]) Results(ctx context.Context) ResultSeq[K, T] {
	return AwaitSeq(ctx, pm.m.Range)
}

func (pm *PromiseMap[K, T]) ReadyResults() (map[K]T, map[K]error) {
	return ReadySeq(pm.m.Range).CollectAll()
}
