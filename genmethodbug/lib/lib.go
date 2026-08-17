package lib

type Alias[T any] = *impl[T]

type impl[T any] struct{}

func (p *impl[T]) M[V any]() Alias[V] {
	return nil
}
