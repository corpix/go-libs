package promise

import (
	"github.com/SlamJam/go-libs/xerrors"
	"github.com/pkg/errors"
)

type Resolver[T any] struct {
	promise Promise[T]
}

func (r Resolver[T]) Resolve(val T) bool {
	return r.promise.fulfill(val, nil)
}

func (r Resolver[T]) Reject(err error) bool {
	var t T

	if err == nil {
		err = errors.New("promise rejected with nil error")
	}

	return r.promise.fulfill(t, err)
}

func (r Resolver[T]) Wrap(val T, err error) bool {
	return r.promise.fulfill(val, err)
}

func (r Resolver[T]) SafetyFulfill(f func() (T, error)) (result bool) {
	return r.Wrap(xerrors.Recover(f))
}
