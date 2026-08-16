package awaitable

import "go.uber.org/multierr"

type MapResult[K comparable] map[K]error

func (m MapResult[K]) HasError() bool {
	for _, err := range m {
		if err != nil {
			return true
		}
	}

	return false
}

func (m MapResult[K]) Err() error {
	var result error

	for _, err := range m {
		multierr.AppendInto(&result, err)
	}

	return result
}
