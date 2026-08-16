package xiter

import (
	"iter"
	"sync"
	"sync/atomic"
)

func FanIn[T any, S iter.Seq[T]](seqs ...S) (iter.Seq[T], func()) {
	ch := make(chan T)

	// Fixme: cycle import
	// defer xchan.Drain(ch)

	var done atomic.Bool
	defer done.Store(true)

	var wg sync.WaitGroup
	for _, seq := range seqs {
		wg.Go(func() {
			if done.Load() {
				return
			}
			for val := range seq {
				ch <- val
				if done.Load() {
					return
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	iter := func(yeild func(T) bool) {
		for val := range ch {
			if !yeild(val) {
				return
			}
		}
	}

	stop := func() { done.Store(true) }

	return iter, stop
}
