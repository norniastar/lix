package lix

import "iter"

// Pair 表示一个强类型键值对。
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// Seq 表示值的惰性序列。
type Seq[V any] iter.Seq[V]

// Iter 将序列转换为标准库迭代器。
func (s Seq[V]) Iter() iter.Seq[V] {
	return iter.Seq[V](s)
}

// Seq2 表示键值对的惰性序列。
// 键必须可比较，因此 Seq2 始终可以收集为 map。
type Seq2[K comparable, V any] iter.Seq2[K, V]

// Iter 将序列转换为标准库键值迭代器。
func (s Seq2[K, V]) Iter() iter.Seq2[K, V] {
	return iter.Seq2[K, V](s)
}
