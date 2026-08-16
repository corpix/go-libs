package xiter

import (
	"iter"
	"sync"
)

func (s XSeq[T]) PullCh[V any]() (func() <-chan T, func()) {
	ch := make(chan T)
	barrier := make(chan struct{}, 1)
	done := make(chan struct{})

	go func() {
		defer close(ch)

		next, stop := iter.Pull(s.AsSeq())
		defer stop()

		for {
			select {
			case <-barrier:
				val, ok := next()
				if !ok {
					return
				}
				select {
				case ch <- val:
				case <-done:
					return
				}

			case <-done:
				return
			}
		}
	}()

	next := func() <-chan T {
		select {
		case barrier <- struct{}{}:
		default:
		}

		return ch
	}

	var onceStop sync.Once
	stop := func() {
		onceStop.Do(func() {
			close(done)
		})
	}

	return next, stop
}

// Batched разбивает исходную последовательность на части размером не более n элементов
func (s XSeq[T]) Batched[V any](n int, flusher <-chan V) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		ch := make(chan T)
		barrier := make(chan struct{}, 1)
		done := make(chan struct{})
		defer close(done)

		go func() {
			defer close(ch)

			next, stop := iter.Pull(s.AsSeq())
			defer stop()

			for {
				select {
				case <-barrier:
					val, ok := next()
					if !ok {
						return
					}
					select {
					case ch <- val:
					case <-done:
						return
					}

				case <-done:
					return
				}
			}
		}()

		getNext := func() <-chan T {
			select {
			case barrier <- struct{}{}:
			default:
			}

			return ch
		}

		batch := make([]T, 0, n)

		flush := func() bool {
			if len(batch) == 0 {
				return true
			}

			result := yield(batch)
			batch = make([]T, 0, n)

			return result
		}

	loop:
		for {
			select {
			case item, ok := <-getNext():
				if !ok {
					break loop
				}

				batch = append(batch, item)
				if len(batch) > n {
					if !flush() {
						break loop
					}
				}
			case <-ch:
				if !flush() {
					break loop
				}
			}
		}

		flush()
	}
}
