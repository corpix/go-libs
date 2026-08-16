package xlinq

// func (s streamComparable[V]) Unique(sizeHint int) StreamComparable[V] {

// 	return streamComparable[V]{
// 		stream: stream[V]{
// 			iterator: xiter.SeqComparable[V](s.iterator).
// 				Unique(s.sizeInfo.GetPreAllocSize()).
// 				ToSeq(),
// 			sizeInfo: s.sizeInfo.CopyWithSize(0),
// 		},
// 	}
// }
