package lix

// Map 惰性转换 s 中的每个值。
func (s Seq[V]) Map[U any](mapper func(V) U) Seq[U] {
	return func(yield func(U) bool) {
		for value := range s {
			if !yield(mapper(value)) {
				return
			}
		}
	}
}

// Filter 惰性保留 predicate 接受的值。
func (s Seq[V]) Filter(predicate func(V) bool) Seq[V] {
	return func(yield func(V) bool) {
		for value := range s {
			if predicate(value) && !yield(value) {
				return
			}
		}
	}
}

// Reject 惰性保留 predicate 拒绝的值。
func (s Seq[V]) Reject(predicate func(V) bool) Seq[V] {
	return s.Filter(func(value V) bool {
		return !predicate(value)
	})
}

// FlatMap 惰性地把每个值映射为序列，并将结果展开一层。
func (s Seq[V]) FlatMap[U any](mapper func(V) Seq[U]) Seq[U] {
	return func(yield func(U) bool) {
		for value := range s {
			for mapped := range mapper(value) {
				if !yield(mapped) {
					return
				}
			}
		}
	}
}

// FlatMap2 惰性地把每个值映射为键值序列，并将结果展开一层。
func (s Seq[V]) FlatMap2[K comparable, U any](mapper func(V) Seq2[K, U]) Seq2[K, U] {
	return func(yield func(K, U) bool) {
		for value := range s {
			for key, mapped := range mapper(value) {
				if !yield(key, mapped) {
					return
				}
			}
		}
	}
}

// Take 惰性返回至多 count 个值；count 非正时返回空序列。
func (s Seq[V]) Take(count int) Seq[V] {
	return func(yield func(V) bool) {
		if count <= 0 {
			return
		}

		taken := 0
		for value := range s {
			if !yield(value) {
				return
			}
			taken++
			if taken == count {
				return
			}
		}
	}
}

// Drop 惰性跳过前 count 个值；count 非正时返回全部值。
func (s Seq[V]) Drop(count int) Seq[V] {
	return func(yield func(V) bool) {
		skipped := 0
		for value := range s {
			if skipped < count {
				skipped++
				continue
			}
			if !yield(value) {
				return
			}
		}
	}
}

// Enumerate 惰性地将从零开始的遍历序号附加到每个值。
func (s Seq[V]) Enumerate() Seq2[int, V] {
	return func(yield func(int, V) bool) {
		index := 0
		for value := range s {
			if !yield(index, value) {
				return
			}
			index++
		}
	}
}

// KeyBy 惰性地使用 selector 为每个值生成键。
func (s Seq[V]) KeyBy[K comparable](selector func(V) K) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for value := range s {
			if !yield(selector(value), value) {
				return
			}
		}
	}
}

// Associate 惰性地将每个值转换为一个键值对。
func (s Seq[V]) Associate[K comparable, U any](mapper func(V) (K, U)) Seq2[K, U] {
	return func(yield func(K, U) bool) {
		for value := range s {
			key, mapped := mapper(value)
			if !yield(key, mapped) {
				return
			}
		}
	}
}

// Map 惰性转换每个值，并保留对应的键。
func (s Seq2[K, V]) Map[U any](mapper func(K, V) U) Seq2[K, U] {
	return func(yield func(K, U) bool) {
		for key, value := range s {
			if !yield(key, mapper(key, value)) {
				return
			}
		}
	}
}

// MapKeys 惰性转换每个键，并保留对应的值。
func (s Seq2[K, V]) MapKeys[L comparable](mapper func(K, V) L) Seq2[L, V] {
	return func(yield func(L, V) bool) {
		for key, value := range s {
			if !yield(mapper(key, value), value) {
				return
			}
		}
	}
}

// Filter 惰性保留 predicate 接受的键值对。
func (s Seq2[K, V]) Filter(predicate func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, value := range s {
			if predicate(key, value) && !yield(key, value) {
				return
			}
		}
	}
}

// Reject 惰性保留 predicate 拒绝的键值对。
func (s Seq2[K, V]) Reject(predicate func(K, V) bool) Seq2[K, V] {
	return s.Filter(func(key K, value V) bool {
		return !predicate(key, value)
	})
}

// FlatMap 惰性地把每个键值对映射为值序列，并将结果展开一层。
func (s Seq2[K, V]) FlatMap[U any](mapper func(K, V) Seq[U]) Seq[U] {
	return func(yield func(U) bool) {
		for key, value := range s {
			for mapped := range mapper(key, value) {
				if !yield(mapped) {
					return
				}
			}
		}
	}
}

// FlatMap2 惰性地把每个键值对映射为新的键值序列，并将结果展开一层。
func (s Seq2[K, V]) FlatMap2[L comparable, U any](mapper func(K, V) Seq2[L, U]) Seq2[L, U] {
	return func(yield func(L, U) bool) {
		for key, value := range s {
			for mappedKey, mappedValue := range mapper(key, value) {
				if !yield(mappedKey, mappedValue) {
					return
				}
			}
		}
	}
}

// Take 惰性返回至多 count 个键值对。
func (s Seq2[K, V]) Take(count int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if count <= 0 {
			return
		}

		taken := 0
		for key, value := range s {
			if !yield(key, value) {
				return
			}
			taken++
			if taken == count {
				return
			}
		}
	}
}

// Drop 惰性跳过前 count 个键值对。
func (s Seq2[K, V]) Drop(count int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skipped := 0
		for key, value := range s {
			if skipped < count {
				skipped++
				continue
			}
			if !yield(key, value) {
				return
			}
		}
	}
}
