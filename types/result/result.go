package result

type Result[T any] struct {
	err   error
	value T
}

func (r Result[T]) IsError() bool {
	return r.err != nil
}

func (r Result[T]) MustValue() T {
	if r.IsError() {
		panic("Result contains no value")
	}
	return r.value
}

func (r Result[T]) ValueOrDefault(t T) T {
	if r.IsError() {
		return t
	}
	return r.value
}

func (r Result[T]) ValueOrZero() T {
	var t T
	return r.ValueOrDefault(t)
}

func (r Result[T]) Err() error {
	return r.err
}

func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

func (r Result[T]) Map[V any](f func(T) (V, error)) Result[V] {
	if r.IsError() {
		return Err[V](r.err)
	}

	return New(f(r.value))
}

func (r Result[T]) FlatMap[V any](f func(T) Result[V]) Result[V] {
	if r.IsError() {
		return Err[V](r.err)
	}

	return f(r.value)
}

func New[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}

	return Ok(value)
}

func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: value,
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		err: err,
	}
}
