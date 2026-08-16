package co

import (
	"errors"
	"iter"

	"github.com/SlamJam/go-libs/co/promise"
	"github.com/SlamJam/go-libs/types/result"
)

type ResultSeq[K comparable, T any] iter.Seq2[K, result.Result[T]]

func (it ResultSeq[K, T]) SkipPending() ResultSeq[K, T] {
	return func(yield func(K, result.Result[T]) bool) {
		for key, res := range it {
			// Нас интересует ТОЛЬКО самая высокоуровневая ошибка
			if errors.Is(res.Err(), promise.ErrAwaitCanceled) {
				continue
			}

			if !yield(key, res) {
				return
			}
		}
	}
}

func (it ResultSeq[K, T]) CollectAll() (map[K]T, map[K]error) {
	values := make(map[K]T)
	errors := make(map[K]error)

	for key, item := range it {
		if val, err := item.Unwrap(); err != nil {
			errors[key] = err
		} else {
			values[key] = val
		}
	}

	return values, errors
}

func (it ResultSeq[K, T]) CollectUntilFirstError() (map[K]T, K, error) {
	values := make(map[K]T)

	for key, item := range it {
		if val, err := item.Unwrap(); err != nil {
			return values, key, err
		} else {
			values[key] = val
		}
	}

	var k K

	return values, k, nil
}

func (it ResultSeq[K, T]) CollectFirstResult() (K, T, bool) {
	for key, item := range it {
		if val, err := item.Unwrap(); err != nil {
			return key, val, true
		}
	}

	var k K
	var t T

	return k, t, false
}
