package xiter

import (
	"iter"
)

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

type XSeqKV[K comparable, V any] XSeq[KV[K, V]]

func ToKV[K comparable, V any](seq iter.Seq[KV[K, V]]) XSeqKV[K, V] {
	return XSeqKV[K, V](seq)
}

func FromMap[M ~map[K]V, K comparable, V any](m M) XSeqKV[K, V] {
	return func(yield func(KV[K, V]) bool) {
		for k, v := range m {
			kv := KV[K, V]{Key: k, Value: v}
			if !yield(kv) {
				return
			}
		}
	}
}

func (seq XSeqKV[K, V]) Collect(capacity int) map[K]V {
	res := make(map[K]V, capacity)
	for kv := range seq {
		res[kv.Key] = kv.Value
	}

	return res
}
