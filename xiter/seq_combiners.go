package xiter

import "iter"

func Concat[T any](items ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range items {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Zip[T any](items ...iter.Seq[T]) iter.Seq[[]T] {
	pulls := make([]func() (T, bool), 0, len(items))
	stops := make([]func(), 0, len(items))

	defer func() {
		for _, stop := range stops {
			stop()
		}
	}()

	for _, seq := range items {
		pull, stop := iter.Pull(seq)
		pulls = append(pulls, pull)
		stops = append(stops, stop)
	}

	return func(yield func([]T) bool) {
		for {
			vs := make([]T, 0, len(items))
			for _, pull := range pulls {
				v, ok := pull()
				if !ok {
					return
				}
				vs = append(vs, v)
			}

			if !yield(vs) {
				return
			}
		}
	}
}
