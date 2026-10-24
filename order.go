package lix

import (
	"cmp"
	"slices"
)

// SortBy 返回按选定值排序的序列。遍历返回序列时才会消费源序列；
// 排序键相等的元素保持原有顺序。
func (s Seq[V]) SortBy[O cmp.Ordered](selector func(V) O) Seq[V] {
	return sortSeq(s, selector, false)
}

// SortByDescending 返回按选定值降序排列的序列。
// 遍历返回序列时才会消费源序列；排序键相等的元素保持原有顺序。
func (s Seq[V]) SortByDescending[O cmp.Ordered](selector func(V) O) Seq[V] {
	return sortSeq(s, selector, true)
}

// Reverse 返回遍历顺序反转后的序列。
// 遍历返回序列时才会消费源序列。
func (s Seq[V]) Reverse() Seq[V] {
	return func(yield func(V) bool) {
		for _, value := range slices.Backward(s.Collect()) {
			if !yield(value) {
				return
			}
		}
	}
}

// SortBy 返回按选定值排序的键值序列。遍历返回序列时才会消费源序列；
// 排序键相等的元素保持原有顺序。
func (s Seq2[K, V]) SortBy[O cmp.Ordered](selector func(K, V) O) Seq2[K, V] {
	return sortSeq2(s, selector, false)
}

// SortByDescending 返回按选定值降序排列的键值序列。
// 遍历返回序列时才会消费源序列；排序键相等的元素保持原有顺序。
func (s Seq2[K, V]) SortByDescending[O cmp.Ordered](selector func(K, V) O) Seq2[K, V] {
	return sortSeq2(s, selector, true)
}

// Reverse 返回遍历顺序反转后的键值序列。
// 遍历返回序列时才会消费源序列。
func (s Seq2[K, V]) Reverse() Seq2[K, V] {
	return func(yield func(K, V) bool) {
		pairs := make([]Pair[K, V], 0)
		for key, value := range s {
			pairs = append(pairs, Pair[K, V]{Key: key, Value: value})
		}
		for _, pair := range slices.Backward(pairs) {
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	}
}

// keyed 保存元素及其排序键，使 selector 对每个元素只调用一次。
type keyed[T any, O cmp.Ordered] struct {
	item  T
	order O
}

// sortKeyed 按排序键稳定排序 items；descending 为 true 时按降序排列。
func sortKeyed[T any, O cmp.Ordered](items []keyed[T, O], descending bool) {
	slices.SortStableFunc(items, func(left, right keyed[T, O]) int {
		if descending {
			return cmp.Compare(right.order, left.order)
		}
		return cmp.Compare(left.order, right.order)
	})
}

func sortSeq[V any, O cmp.Ordered](s Seq[V], selector func(V) O, descending bool) Seq[V] {
	return func(yield func(V) bool) {
		items := make([]keyed[V, O], 0)
		for value := range s {
			items = append(items, keyed[V, O]{item: value, order: selector(value)})
		}

		sortKeyed(items, descending)

		for _, item := range items {
			if !yield(item.item) {
				return
			}
		}
	}
}

func sortSeq2[K comparable, V any, O cmp.Ordered](s Seq2[K, V], selector func(K, V) O, descending bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		items := make([]keyed[Pair[K, V], O], 0)
		for key, value := range s {
			items = append(items, keyed[Pair[K, V], O]{
				item:  Pair[K, V]{Key: key, Value: value},
				order: selector(key, value),
			})
		}

		sortKeyed(items, descending)

		for _, item := range items {
			if !yield(item.item.Key, item.item.Value) {
				return
			}
		}
	}
}
