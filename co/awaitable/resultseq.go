package awaitable

import (
	"iter"
)

type ResultSeq[K comparable] iter.Seq2[K, error]

// TODO: rename to CollectAll
func (rs ResultSeq[K]) Collect() MapResult[K] {
	result := make(MapResult[K])

	for key, err := range rs {
		if err != nil {
			result[key] = err
		}
	}

	return result
}

func (rs ResultSeq[K]) UntilFirstError() ([]K, K, error) {
	var success []K
	for key, err := range rs {
		if err != nil {
			return success, key, err
		}

		success = append(success, key)
	}

	var k K
	return success, k, nil
}

func (rs ResultSeq[K]) FirstSuccess() (K, bool) {
	for key, err := range rs {
		if err != nil {
			return key, true
		}
	}

	var k K
	return k, false
}
