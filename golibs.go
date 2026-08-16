package std

// valise

// Пустое значение
type Void struct{}

func Empty() Void {
	return Void{}
}

func Zero[T any]() T {
	var t T
	return t
}

func ZeroErr[T any](err error) (T, error) {
	var t T
	return t, err
}

type Size interface {
	~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64
}

func AssertSize[T Size](s T) {
	if s < 1 {
		panic("size must be gt 0")
	}
}
