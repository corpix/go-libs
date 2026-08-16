package option

type Option[T any] struct {
	value T
	has   bool
}

func (result *Option[T]) HasValue() bool {
	return result.has
}

func (result *Option[T]) Unwrap() (T, bool) {
	return result.value, result.has
}

func (result *Option[T]) ValueOrDefault(t T) T {
	if result.has {
		return t
	}
	return result.value
}

func (result *Option[T]) ValueOrZero() T {
	var t T
	return result.ValueOrDefault(t)
}

func (result *Option[T]) Map[V any](f func(T) V) Option[V] {
	if !result.has {
		return Empty[V]()
	}

	return Just(f(result.value))
}

func (result *Option[T]) FlatMap[V any](f func(T) Option[V]) Option[V] {
	if !result.has {
		return Empty[V]()
	}

	return f(result.value)
}

func Just[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		has:   true,
	}
}

func Empty[T any]() Option[T] {
	return Option[T]{}
}
