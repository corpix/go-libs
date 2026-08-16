package promise

func NewDeferred[T any]() (Promise[T], Resolver[T]) {
	p := newPromise[T]()

	return p, Resolver[T]{promise: p}
}
