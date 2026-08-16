package xiter

import (
	"iter"
)

type XSeq2[K, V any] iter.Seq2[K, V]

func (xseq2 XSeq2[K, V]) AsSeq2() iter.Seq2[K, V] {
	return iter.Seq2[K, V](xseq2)
}

func (xseq2 XSeq2[K, V]) Keys() XSeq[K] {
	return func(yield func(K) bool) {
		for k := range xseq2 {
			if !yield(k) {
				return
			}
		}
	}
}

func (xseq2 XSeq2[K, V]) Values() XSeq[V] {
	return func(yield func(V) bool) {
		for _, v := range xseq2 {
			if !yield(v) {
				return
			}
		}
	}
}
