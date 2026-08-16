package xlinq

// import (
// 	"maps"

// 	"github.com/SlamJam/go-libs/xiter"
// )

// func FromMap[K comparable, V any](m map[K]V) streamKV[K, V] {
// 	return streamKV[K, V]{
// 		stream: stream[xiter.KV[K, V]]{
// 			iterator: xiter.FromMap(m),
// 			size:     len(m),
// 		},
// 	}
// }

// func FromMapCopy[K comparable, V any](m map[K]V) streamKV[K, V] {
// 	data := make(map[K]V, len(m))
// 	maps.Copy(data, m)

// 	return FromMap(data)
// }
