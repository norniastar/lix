package lix

import (
	"cmp"
	"slices"
)

// SortBy 返回按选定值排序的序列。遍历返回序列时才会消费源序列；
// 排序键相等的元素保持原有顺序。
func (s Seq[V]) SortBy[O cmp.Ordered](selector func(V) O) Seq[V] {
	type item struct {
		value V
		order O
	}

	return func(yield func(V) bool) {
		items := make([]item, 0)
		for value := range s {
			items = append(items, item{value: value, order: selector(value)})
		}

		slices.SortStableFunc(items, func(left, right item) int {
			return cmp.Compare(left.order, right.order)
		})

		for _, item := range items {
			if !yield(item.value) {
				return
			}
		}
	}
}

// SortByDescending 返回按选定值降序排列的序列。
// 遍历返回序列时才会消费源序列；排序键相等的元素保持原有顺序。
func (s Seq[V]) SortByDescending[O cmp.Ordered](selector func(V) O) Seq[V] {
	type item struct {
		value V
		order O
	}

	return func(yield func(V) bool) {
		items := make([]item, 0)
		for value := range s {
			items = append(items, item{value: value, order: selector(value)})
		}

		slices.SortStableFunc(items, func(left, right item) int {
			return cmp.Compare(right.order, left.order)
		})

		for _, item := range items {
			if !yield(item.value) {
				return
			}
		}
	}
}

// Reverse 返回遍历顺序反转后的序列。
// 遍历返回序列时才会消费源序列。
func (s Seq[V]) Reverse() Seq[V] {
	return func(yield func(V) bool) {
		values := s.Collect()
		for index := len(values) - 1; index >= 0; index-- {
			if !yield(values[index]) {
				return
			}
		}
	}
}

// SortBy 返回按选定值排序的键值序列。遍历返回序列时才会消费源序列；
// 排序键相等的元素保持原有顺序。
func (s Seq2[K, V]) SortBy[O cmp.Ordered](selector func(K, V) O) Seq2[K, V] {
	type item struct {
		key   K
		value V
		order O
	}

	return func(yield func(K, V) bool) {
		items := make([]item, 0)
		for key, value := range s {
			items = append(items, item{
				key:   key,
				value: value,
				order: selector(key, value),
			})
		}

		slices.SortStableFunc(items, func(left, right item) int {
			return cmp.Compare(left.order, right.order)
		})

		for _, item := range items {
			if !yield(item.key, item.value) {
				return
			}
		}
	}
}

// SortByDescending 返回按选定值降序排列的键值序列。
// 遍历返回序列时才会消费源序列；排序键相等的元素保持原有顺序。
func (s Seq2[K, V]) SortByDescending[O cmp.Ordered](selector func(K, V) O) Seq2[K, V] {
	type item struct {
		key   K
		value V
		order O
	}

	return func(yield func(K, V) bool) {
		items := make([]item, 0)
		for key, value := range s {
			items = append(items, item{
				key:   key,
				value: value,
				order: selector(key, value),
			})
		}

		slices.SortStableFunc(items, func(left, right item) int {
			return cmp.Compare(right.order, left.order)
		})

		for _, item := range items {
			if !yield(item.key, item.value) {
				return
			}
		}
	}
}

// Reverse 返回遍历顺序反转后的键值序列。
// 遍历返回序列时才会消费源序列。
func (s Seq2[K, V]) Reverse() Seq2[K, V] {
	type item struct {
		key   K
		value V
	}

	return func(yield func(K, V) bool) {
		items := make([]item, 0)
		for key, value := range s {
			items = append(items, item{key: key, value: value})
		}
		for index := len(items) - 1; index >= 0; index-- {
			item := items[index]
			if !yield(item.key, item.value) {
				return
			}
		}
	}
}
