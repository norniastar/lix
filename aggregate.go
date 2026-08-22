package lix

// Reduce 消费 s，并将其中的值合并为一个结果。
func (s Seq[V]) Reduce[U any](reducer func(U, V) U, initial U) U {
	result := initial
	for value := range s {
		result = reducer(result, value)
	}
	return result
}

// GroupBy 按选出的键对值进行分组。
// 遍历返回序列时会消费源序列，分组按键首次出现的顺序产出。
func GroupBy[V any, K comparable](s Seq[V], selector func(V) K) Seq2[K, []V] {
	return func(yield func(K, []V) bool) {
		groups := make(map[K][]V)
		keys := make([]K, 0)
		for value := range s {
			key := selector(value)
			if _, exists := groups[key]; !exists {
				keys = append(keys, key)
			}
			groups[key] = append(groups[key], value)
		}

		for _, key := range keys {
			if !yield(key, groups[key]) {
				return
			}
		}
	}
}

// UniqBy 针对每个选出的键，惰性保留首次出现的值。
func (s Seq[V]) UniqBy[K comparable](selector func(V) K) Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[K]struct{})
		for value := range s {
			key := selector(value)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			if !yield(value) {
				return
			}
		}
	}
}

// UniqBy 针对每个选出的分组键，惰性保留首次出现的键值对。
func (s Seq2[K, V]) UniqBy[G comparable](selector func(K, V) G) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seen := make(map[G]struct{})
		for key, value := range s {
			group := selector(key, value)
			if _, exists := seen[group]; exists {
				continue
			}
			seen[group] = struct{}{}
			if !yield(key, value) {
				return
			}
		}
	}
}

// Reduce 消费 s，并将其中的键值对合并为一个结果。
func (s Seq2[K, V]) Reduce[U any](reducer func(U, K, V) U, initial U) U {
	result := initial
	for key, value := range s {
		result = reducer(result, key, value)
	}
	return result
}

// GroupPairsBy 按选出的键对键值序列中的值进行分组。
// 遍历返回序列时会消费源序列，分组按键首次出现的顺序产出。
func GroupPairsBy[K comparable, V any, G comparable](s Seq2[K, V], selector func(K, V) G) Seq2[G, []V] {
	return func(yield func(G, []V) bool) {
		groups := make(map[G][]V)
		keys := make([]G, 0)
		for key, value := range s {
			group := selector(key, value)
			if _, exists := groups[group]; !exists {
				keys = append(keys, group)
			}
			groups[group] = append(groups[group], value)
		}

		for _, group := range keys {
			if !yield(group, groups[group]) {
				return
			}
		}
	}
}
