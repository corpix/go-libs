package co

import (
	"iter"

	"github.com/SlamJam/go-libs/types/result"
	"go.uber.org/multierr"
)

type valueWithKey[K comparable, T any] struct {
	Key   K
	Value T
}

type errorWithKey[K comparable] struct {
	Key K
	Err error
}

type ResultSet[K comparable, T any] struct {
	resultsWithKey []valueWithKey[K, T]
	errorsWithKey  []errorWithKey[K]
}

func (pr *ResultSet[K, T]) addError(key K, err error) {
	pr.errorsWithKey = append(pr.errorsWithKey, errorWithKey[K]{Key: key, Err: err})
}

func (pr *ResultSet[K, T]) addValue(key K, value T) {
	pr.resultsWithKey = append(pr.resultsWithKey, valueWithKey[K, T]{Key: key, Value: value})
}

func (pr *ResultSet[K, T]) addResult(key K, result result.Result[T]) {
	if val, err := result.Unwrap(); err == nil {
		pr.addValue(key, val)
	} else {
		pr.addError(key, err)
	}
}

func (pr ResultSet[K, T]) Values() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		for _, item := range pr.resultsWithKey {
			if !yield(item.Key, item.Value) {
				return
			}
		}
	}
}

func (pr ResultSet[T, K]) Count() int {
	return len(pr.resultsWithKey) + len(pr.errorsWithKey)
}

func (pr ResultSet[K, T]) MultiErr() error {
	var result error
	for _, err := range pr.errorsWithKey {
		multierr.AppendInto(&result, err.Err)
	}

	return result
}

func (pr ResultSet[K, T]) ValuesMap() map[K]T {
	result := make(map[K]T, len(pr.resultsWithKey))
	for _, r := range pr.resultsWithKey {
		result[r.Key] = r.Value
	}

	return result
}

func (pr ResultSet[K, T]) ErrorsMap() map[K]error {
	result := make(map[K]error, len(pr.errorsWithKey))
	for _, e := range pr.errorsWithKey {
		result[e.Key] = e.Err
	}

	return result
}
