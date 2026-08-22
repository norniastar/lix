package lix

import (
	"iter"
	"maps"
	"slices"
)

// FromSlice 基于 source 创建可重复遍历的序列。
func FromSlice[S ~[]V, V any](source S) Seq[V] {
	return Seq[V](slices.Values(source))
}

// FromMap 基于 source 创建可重复遍历的键值序列。
// map 的遍历顺序不确定。
func FromMap[M ~map[K]V, K comparable, V any](source M) Seq2[K, V] {
	return Seq2[K, V](maps.All(source))
}

// FromSeq 将标准库迭代器转换为 Seq。
func FromSeq[V any](source iter.Seq[V]) Seq[V] {
	return Seq[V](source)
}

// FromSeq2 将标准库键值迭代器转换为 Seq2。
func FromSeq2[K comparable, V any](source iter.Seq2[K, V]) Seq2[K, V] {
	return Seq2[K, V](source)
}

// FromPairs 基于强类型键值对切片创建可重复遍历的键值序列。
func FromPairs[S ~[]Pair[K, V], K comparable, V any](source S) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, pair := range source {
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	}
}

// Of 创建包含 values 的可重复遍历序列。
func Of[V any](values ...V) Seq[V] {
	return FromSlice(values)
}

// RangeWithSteps 创建从 start 开始、到 stop 之前结束的序列。
// step 为零时返回空序列；正数递增，负数递减。
func RangeWithSteps(start, stop, step int) Seq[int] {
	if step == 0 {
		return func(func(int) bool) {}
	}

	return func(yield func(int) bool) {
		for current := start; rangeContinues(current, stop, step); {
			if !yield(current) {
				return
			}

			next := current + step
			if rangeOverflows(current, next, step) {
				return
			}
			current = next
		}
	}
}

func rangeContinues(current, stop, step int) bool {
	if step > 0 {
		return current < stop
	}
	return current > stop
}

func rangeOverflows(current, next, step int) bool {
	if step > 0 {
		return next <= current
	}
	return next >= current
}
