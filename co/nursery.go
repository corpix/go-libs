package co

import (
	"context"
	"sync"
	"time"

	"github.com/SlamJam/go-libs/co/awaitable"
	"github.com/SlamJam/go-libs/co/promise"
	"github.com/pkg/errors"
)

var (
	ErrCancelled = errors.New("nursery context canceled")
)

type Nursery struct {
	ctx      context.Context
	cancel   context.CancelFunc
	promises []awaitable.Awaitable
	mu       *sync.Mutex
	// isCompleted   *atomic.Bool
	isInitialized bool
}

type NurseryResult struct {
	promises []awaitable.Awaitable
}

func (nr NurseryResult) Await(ctx context.Context) error {
	// n.assertIsInitialized()

	return awaitable.FromSlice(nr.promises).AwaitResults(ctx).Collect().Err()
}

func (n *Nursery) assertIsInitialized() {
	if !n.isInitialized {
		panic("Nursery is not initialized")
	}
}

func (n *Nursery) Ctx() context.Context {
	n.assertIsInitialized()

	return n.ctx
}

// func (n *Nursery) IsCompleted() bool {
// 	n.assertIsInitialized()

// 	return n.isCompleted.Load()
// }

func (n *Nursery) getResult() NurseryResult {
	n.assertIsInitialized()

	return NurseryResult{
		promises: n.promises,
	}
}

func (n *Nursery) onComplete() {
	n.assertIsInitialized()

	n.cancel()
	// n.isCompleted.Store(true)
}

func NewNursery(ctx context.Context) Nursery {
	ctx, cancel := context.WithCancelCause(ctx)

	return Nursery{
		ctx: ctx,
		mu:  &sync.Mutex{},
		// isCompleted:   &atomic.Bool{},
		isInitialized: true,
		cancel:        func() { cancel(ErrCancelled) },
	}
}

func WithContext(ctx context.Context, f func(n Nursery)) NurseryResult {
	n := NewNursery(ctx)

	defer n.onComplete()

	f(n)

	return n.getResult()
}

func WithContextResult[RES any](ctx context.Context, f func(n Nursery) (RES, error)) (NurseryResult, RES, error) {
	n := NewNursery(ctx)

	defer n.onComplete()

	res, err := f(n)
	return n.getResult(), res, err
}

// f - лямбда
func (n *Nursery) Fork[T any](f func() (T, error)) promise.Promise[T] {
	p := S.Launch(f)

	n.mu.Lock()
	defer n.mu.Unlock()

	n.promises = append(n.promises, p)

	return p
}

func (n *Nursery) ForkInMultiPromise[T any](mp PromiseSlice[T], f func() (T, error)) {
	p := n.Fork(f)
	mp.Append(p)
}

func xxx1() {
	type Foo struct{}
	type Bar struct{}
	type Baz struct{}

	type Result struct {
		Foo Foo
		Baz Baz
	}

	WithContext(context.TODO(), func(n Nursery) {
		var mpFoo PromiseSlice[Foo]

		for range 5 {
			n.ForkInMultiPromise(mpFoo, func() (Foo, error) {
				n.Ctx()

				return Foo{}, nil
			})
		}

		pBar := n.Fork(func() (Bar, error) {
			n.Ctx()

			return Bar{}, nil
		})

		pBaz := n.Fork(func() (Baz, error) {
			bar, err := pBar.GetOrAwait(n.Ctx())
			if err != nil {
				return Baz{}, err
			}

			_ = bar

			return Baz{}, nil
		})

		foos, err1 := mpFoo.AwaitAllOrFirstError(n.Ctx())
		_ = foos
		_ = err1

		baz, err2 := pBaz.GetOrAwait(n.Ctx())
		_ = baz
		_ = err2
	})
	// тут n.Ctx() уже будет кенсельнут
}

func xxx2() {
	type Foo struct{}
	type Bar struct{}
	type Baz struct{}

	type Result struct {
		Foo Foo
		Baz Baz
	}

	getFooPromises := func(n Nursery, count int) PromiseSlice[Foo] {
		var mpFoo PromiseSlice[Foo]

		for range count {
			n.ForkInMultiPromise(mpFoo, func() (Foo, error) {
				n.Ctx()

				return Foo{}, nil
			})
		}

		return mpFoo
	}

	getBarPromise := func(n Nursery) promise.Promise[Bar] {
		return n.Fork(func() (Bar, error) {
			n.Ctx()

			return Bar{}, nil
		})
	}

	getBaz := func(ctx context.Context, bar Bar) (Baz, error) {
		_, _ = ctx, bar

		return Baz{}, nil
	}

	getBazPromise := func(n Nursery, pBar promise.Promise[Bar]) promise.Promise[Baz] {
		return n.Fork(func() (Baz, error) {
			bar, err := pBar.GetOrAwait(n.Ctx())
			if err != nil {
				return Baz{}, err
			}

			// получаем Baz, используя Bar
			return getBaz(n.Ctx(), bar)
		})
	}

	WithContext(context.TODO(), func(n Nursery) {
		mpFoo := getFooPromises(n, 5)
		pBar := getBarPromise(n)
		pBaz := getBazPromise(n, pBar)

		foos, err1 := mpFoo.AwaitAllOrFirstError(n.Ctx())
		_, _ = foos, err1

		baz, err2 := pBaz.GetOrAwait(n.Ctx())
		_, _ = baz, err2
	})
	// тут n.Ctx() уже будет кенсельнут

	_, res, err := WithContextResult(context.TODO(), func(n Nursery) (int, error) {
		mpFoo := getFooPromises(n, 5)
		pBar := getBarPromise(n)
		pBaz := getBazPromise(n, pBar)

		foos, err1 := mpFoo.AwaitAllOrFirstError(n.Ctx())
		_, _ = foos, err1

		baz, err2 := pBaz.GetOrAwait(n.Ctx())
		_, _ = baz, err2

		return 42, nil
	})

	_, _ = res, err
}

type Response struct {
}

func RequesReplica(context.Context, string) (Response, error) {
	return Response{}, nil
}

func requesShard(n Nursery, addrs []string) promise.Promise[Response] {
	var replicaReqs PromiseSlice[Response]

	for _, addr := range addrs {
		n.ForkInMultiPromise(replicaReqs, func() (Response, error) {
			return RequesReplica(n.Ctx(), addr)
		})
	}

	return n.Fork(func() (Response, error) {
		_, resp, _ := replicaReqs.Results(n.Ctx()).CollectFirstResult()
		return resp, nil
	})
}

var ErrReplicaResultTimeout = errors.New("replica time budget exeeded")

func requesShardWithDelay(n Nursery, addrs []string) promise.Promise[Response] {
	return n.Fork(func() (Response, error) {
		var replicaReqs PromiseSlice[Response]

		for _, addr := range addrs {
			n.ForkInMultiPromise(replicaReqs, func() (Response, error) {
				return RequesReplica(n.Ctx(), addr)
			})

			waitCtx, cancel := context.WithTimeoutCause(n.Ctx(), 50*time.Millisecond, ErrReplicaResultTimeout)
			defer cancel()

			_, resp, _ := replicaReqs.Results(waitCtx).CollectFirstResult()
			// if err == nil {
			return resp, nil
			// }
		}

		_, resp, _ := replicaReqs.Results(n.Ctx()).CollectFirstResult()
		return resp, nil
	})
}

// P1        | P2        | FirstResult(waitCtx)

// In-Flight | In-Flight | error(ErrReplicaResultTimeout)
// In-Flight | Result    | Result
// In-Flight | Error     | error(ErrReplicaResultTimeout)
// Result    | In-Flight | Result
// Result    | Result    | Result (one of)
// Result    | Error     | Result
// Error     | In-Flight | error(ErrReplicaResultTimeout)
// Error     | Result    | Result
// Error     | Error     | multierr

func xxx3() {
	// shard -> []replica
	cluster := [][]string{
		{"shard1.replica1", "shard1.replica2", "shard1.replica3"},
		{"shard2.replica1", "shard2.replica2", "shard2.replica3"},
		{"shard3.replica1", "shard3.replica2", "shard3.replica3"},
		{"shard4.replica1", "shard4.replica2", "shard4.replica3"},
	}

	nr, resp, err := WithContextResult(context.TODO(), func(n Nursery) ([]Response, error) {
		var shardReqs PromiseSlice[Response]
		for _, shard := range cluster {
			p := requesShard(n, shard)
			shardReqs.Append(p)
		}

		// Хотим все результаты
		// return shardReqs.AllResults(n.Ctx())

		// Зачем ждать все, если кто-то не ответил?
		// return shardReqs.AllResultsOrFirstError(n.Ctx())

		// Соберём частичный результат
		// partialResult := shardReqs.PartialResult(n.Ctx())
		// if err := partialResult.MultiErr(); err != nil {
		// 	log.Printf("WARN: partial result with errors: %v", err)
		// }

		// return partialResult.AvailableResults(), nil
		return nil, nil
	})

	_, _ = resp, err

	// Wait until all jobs to be done
	// Wait until shutdown complete
	nr.Await(context.TODO())

	// Focus
	manyNurseryResults := awaitable.FromItems(nr, nr, nr)
	_ = manyNurseryResults.AwaitResults(context.TODO()).Collect()
}
