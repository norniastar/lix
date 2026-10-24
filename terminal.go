package lix

import (
	"maps"
	"slices"
)

// Collect 消费 s，并以切片形式返回其中的值。
func (s Seq[V]) Collect() []V {
	// 与 slices.Collect 不同，空序列返回非 nil 的空切片。
	return slices.AppendSeq(make([]V, 0), s.Iter())
}

// First 返回 s 中的第一个值。
func (s Seq[V]) First() (V, bool) {
	for value := range s {
		return value, true
	}
	var zero V
	return zero, false
}

// Last 返回 s 中的最后一个值。
func (s Seq[V]) Last() (V, bool) {
	var last V
	found := false
	for value := range s {
		last = value
		found = true
	}
	return last, found
}

// Find 返回 s 中第一个被 predicate 接受的值。
func (s Seq[V]) Find(predicate func(V) bool) (V, bool) {
	for value := range s {
		if predicate(value) {
			return value, true
		}
	}
	var zero V
	return zero, false
}

// FindIndex 返回 s 中第一个被 predicate 接受的值的遍历序号。
func (s Seq[V]) FindIndex(predicate func(V) bool) (int, bool) {
	index := 0
	for value := range s {
		if predicate(value) {
			return index, true
		}
		index++
	}
	return -1, false
}

// SomeBy 报告是否至少有一个值被 predicate 接受。
func (s Seq[V]) SomeBy(predicate func(V) bool) bool {
	for value := range s {
		if predicate(value) {
			return true
		}
	}
	return false
}

// EveryBy 报告是否所有值都被 predicate 接受。
func (s Seq[V]) EveryBy(predicate func(V) bool) bool {
	for value := range s {
		if !predicate(value) {
			return false
		}
	}
	return true
}

// Length 返回 s 中值的数量。
func (s Seq[V]) Length() int {
	count := 0
	for range s {
		count++
	}
	return count
}

// ForEach 消费 s，并按遍历顺序对每个值执行 action。
func (s Seq[V]) ForEach(action func(V)) {
	for value := range s {
		action(value)
	}
}

// Collect 消费 s，并以 map 形式返回其中的键值对。
// 如果 s 中存在重复键，则保留最后一个值。
func (s Seq2[K, V]) Collect() map[K]V {
	return maps.Collect(s.Iter())
}

// First 返回 s 中的第一个键值对。
func (s Seq2[K, V]) First() (K, V, bool) {
	for key, value := range s {
		return key, value, true
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}

// Last 返回 s 中的最后一个键值对。
func (s Seq2[K, V]) Last() (K, V, bool) {
	var lastKey K
	var lastValue V
	found := false
	for key, value := range s {
		lastKey = key
		lastValue = value
		found = true
	}
	return lastKey, lastValue, found
}

// Find 返回 s 中第一个被 predicate 接受的键值对。
func (s Seq2[K, V]) Find(predicate func(K, V) bool) (K, V, bool) {
	for key, value := range s {
		if predicate(key, value) {
			return key, value, true
		}
	}
	var zeroKey K
	var zeroValue V
	return zeroKey, zeroValue, false
}

// FindIndex 返回 s 中第一个被 predicate 接受的键值对的遍历序号。
func (s Seq2[K, V]) FindIndex(predicate func(K, V) bool) (int, bool) {
	index := 0
	for key, value := range s {
		if predicate(key, value) {
			return index, true
		}
		index++
	}
	return -1, false
}

// SomeBy 报告是否至少有一个键值对被 predicate 接受。
func (s Seq2[K, V]) SomeBy(predicate func(K, V) bool) bool {
	for key, value := range s {
		if predicate(key, value) {
			return true
		}
	}
	return false
}

// EveryBy 报告是否所有键值对都被 predicate 接受。
func (s Seq2[K, V]) EveryBy(predicate func(K, V) bool) bool {
	for key, value := range s {
		if !predicate(key, value) {
			return false
		}
	}
	return true
}

// Length 返回 s 中键值对的数量。
func (s Seq2[K, V]) Length() int {
	count := 0
	for range s {
		count++
	}
	return count
}

// ForEach 消费 s，并按遍历顺序对每个键值对执行 action。
func (s Seq2[K, V]) ForEach(action func(K, V)) {
	for key, value := range s {
		action(key, value)
	}
}

// Keys 惰性返回 s 中的键。
func (s Seq2[K, V]) Keys() Seq[K] {
	return func(yield func(K) bool) {
		for key := range s {
			if !yield(key) {
				return
			}
		}
	}
}

// Values 惰性返回 s 中的值。
func (s Seq2[K, V]) Values() Seq[V] {
	return func(yield func(V) bool) {
		for _, value := range s {
			if !yield(value) {
				return
			}
		}
	}
}

// Pairs 惰性地将键值序列转换为 Pair 序列。
func Pairs[K comparable, V any](s Seq2[K, V]) Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for key, value := range s {
			if !yield(Pair[K, V]{Key: key, Value: value}) {
				return
			}
		}
	}
}
