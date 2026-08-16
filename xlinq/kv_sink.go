package xlinq

// import (
// 	"maps"

// 	"github.com/SlamJam/go-libs/xiter"
// )

// func (s streamKV[K, V]) IterKV() xiter.Seq2[K, V] {
// 	return func(yield func(K, V) bool) {
// 		for kv := range s.Iter() {
// 			if !yield(kv.Key, kv.Value) {
// 				return
// 			}
// 		}
// 	}
// }

// func (s streamKV[K, V]) IterKeys() xiter.Seq[K] {
// 	return s.IterKV().Keys()
// }

// func (s streamKV[K, V]) IterValues() xiter.Seq[V] {
// 	return s.IterKV().Values()
// }

// func (s streamKV[K, V]) EnumKeys() stream[K] {
// 	return stream[K]{
// 		iterator: s.IterKeys(),
// 		size:     s.size,
// 	}
// }

// func (s streamKV[K, V]) EnumValues() stream[V] {
// 	return stream[V]{
// 		iterator: s.IterValues(),
// 		size:     s.size,
// 	}
// }

// func (s streamKV[K, V]) ForEach(f func(K, V)) {
// 	for k, v := range s.IterKV() {
// 		f(k, v)
// 	}
// }

// func (s streamKV[K, V]) ToMap() map[K]V {
// 	return maps.Collect(s.IterKV().ToIterSeq2())
// }
