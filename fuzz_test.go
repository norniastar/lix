package lix_test

import (
	"math"
	"math/big"
	"slices"
	"testing"

	"github.com/norniastar/lix"
)

func FuzzRangeWithStepsMatchesBoundedReference(f *testing.F) {
	f.Add(int16(1), int16(10), int8(2))
	f.Add(int16(10), int16(-5), int8(-3))
	f.Add(int16(5), int16(5), int8(1))
	f.Add(int16(0), int16(10), int8(0))

	f.Fuzz(func(t *testing.T, rawStart, rawStop int16, rawStep int8) {
		start := int(rawStart % 50)
		stop := int(rawStop % 50)
		step := int(rawStep % 10)

		want := make([]int, 0)
		if step > 0 {
			for value := start; value < stop; value += step {
				want = append(want, value)
			}
		} else if step < 0 {
			for value := start; value > stop; value += step {
				want = append(want, value)
			}
		}

		if result := lix.RangeWithSteps(start, stop, step).Collect(); !slices.Equal(result, want) {
			t.Fatalf("RangeWithSteps(%d, %d, %d) = %v, want %v", start, stop, step, result, want)
		}
	})
}

func FuzzTakeAndDropMatchSliceModel(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4}, uint8(1), uint8(2))
	f.Add([]byte{}, uint8(5), uint8(5))

	f.Fuzz(func(t *testing.T, data []byte, rawDrop, rawTake uint8) {
		if len(data) > 64 {
			data = data[:64]
		}
		values := make([]int, len(data))
		for index, value := range data {
			values[index] = int(value)
		}

		drop := min(int(rawDrop), len(values))
		take := min(int(rawTake), len(values)-drop)
		want := values[drop : drop+take]
		result := lix.FromSlice(values).Drop(int(rawDrop)).Take(int(rawTake)).Collect()
		if !slices.Equal(result, want) {
			t.Fatalf("Drop(%d).Take(%d) = %v, want %v", rawDrop, rawTake, result, want)
		}
	})
}

func FuzzReverseIsInvolution(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 128 {
			data = data[:128]
		}
		result := lix.FromSlice(data).Reverse().Reverse().Collect()
		if !slices.Equal(result, data) {
			t.Fatalf("double Reverse() = %v, want %v", result, data)
		}
	})
}

func FuzzUniqByKeepsFirstOccurrence(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4, 5}, uint8(2))
	f.Add([]byte{}, uint8(1))

	f.Fuzz(func(t *testing.T, data []byte, rawModulo uint8) {
		if len(data) > 128 {
			data = data[:128]
		}
		modulo := int(rawModulo%16) + 1
		seen := make(map[int]struct{})
		want := make([]byte, 0)
		for _, value := range data {
			key := int(value) % modulo
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			want = append(want, value)
		}

		result := lix.FromSlice(data).UniqBy(func(value byte) int {
			return int(value) % modulo
		}).Collect()
		if !slices.Equal(result, want) {
			t.Fatalf("UniqBy() = %v, want %v", result, want)
		}
	})
}

func FuzzPairsRoundTrip(f *testing.F) {
	f.Add([]byte{1, 2, 1}, []byte{10, 20, 30})
	f.Add([]byte{}, []byte{})

	f.Fuzz(func(t *testing.T, keys, values []byte) {
		size := min(len(keys), len(values), 128)
		pairs := make([]lix.Pair[byte, byte], size)
		for index := range size {
			pairs[index] = lix.Pair[byte, byte]{Key: keys[index], Value: values[index]}
		}

		result := lix.Pairs(lix.FromPairs(pairs)).Collect()
		if !slices.Equal(result, pairs) {
			t.Fatalf("pair round trip = %v, want %v", result, pairs)
		}
	})
}

func FuzzSortByIsStable(f *testing.F) {
	type item struct {
		Key      byte
		Position int
	}

	f.Add([]byte{2, 1, 2, 1})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 128 {
			data = data[:128]
		}
		items := make([]item, len(data))
		for index, key := range data {
			items[index] = item{Key: key, Position: index}
		}

		result := lix.FromSlice(items).SortBy(func(value item) byte { return value.Key }).Collect()
		if len(result) != len(items) {
			t.Fatalf("sorted length = %d, want %d", len(result), len(items))
		}
		for index := 1; index < len(result); index++ {
			previous := result[index-1]
			current := result[index]
			if previous.Key > current.Key {
				t.Fatalf("result is not sorted at %d: %v", index, result)
			}
			if previous.Key == current.Key && previous.Position > current.Position {
				t.Fatalf("result is not stable at %d: %v", index, result)
			}
		}
	})
}

func FuzzGroupByPartitionsValues(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4, 5}, uint8(2))
	f.Add([]byte{}, uint8(1))

	f.Fuzz(func(t *testing.T, data []byte, rawModulo uint8) {
		if len(data) > 128 {
			data = data[:128]
		}
		modulo := int(rawModulo%16) + 1
		wantGroups := make(map[int][]byte)
		wantKeys := make([]int, 0)
		for _, value := range data {
			key := int(value) % modulo
			if _, exists := wantGroups[key]; !exists {
				wantKeys = append(wantKeys, key)
			}
			wantGroups[key] = append(wantGroups[key], value)
		}

		result := lix.Pairs(lix.GroupBy(lix.FromSlice(data), func(value byte) int {
			return int(value) % modulo
		})).Collect()
		if len(result) != len(wantKeys) {
			t.Fatalf("group count = %d, want %d", len(result), len(wantKeys))
		}
		for index, pair := range result {
			wantKey := wantKeys[index]
			if pair.Key != wantKey || !slices.Equal(pair.Value, wantGroups[wantKey]) {
				t.Fatalf("group %d = %v, want key %d values %v", index, pair, wantKey, wantGroups[wantKey])
			}
		}
	})
}

// FuzzRangeWithStepsMatchesExactReference 覆盖完整 int 范围，
// 参考模型使用 big.Int 精确运算，因此可以验证溢出前停止的契约。
func FuzzRangeWithStepsMatchesExactReference(f *testing.F) {
	f.Add(int64(math.MaxInt64-1), int64(math.MaxInt64), int64(2))
	f.Add(int64(math.MinInt64+1), int64(math.MinInt64), int64(-2))
	f.Add(int64(math.MinInt64), int64(math.MaxInt64), int64(math.MaxInt64))
	f.Add(int64(math.MaxInt64), int64(math.MinInt64), int64(math.MinInt64))
	f.Add(int64(-3), int64(3), int64(1))
	f.Add(int64(0), int64(10), int64(0))

	const limit = 64
	f.Fuzz(func(t *testing.T, rawStart, rawStop, rawStep int64) {
		start, stop, step := int(rawStart), int(rawStop), int(rawStep)

		want := make([]int, 0, limit)
		if step != 0 {
			current := big.NewInt(int64(start))
			bigStop := big.NewInt(int64(stop))
			bigStep := big.NewInt(int64(step))
			for len(want) < limit {
				comparison := current.Cmp(bigStop)
				if (step > 0 && comparison >= 0) || (step < 0 && comparison <= 0) {
					break
				}
				want = append(want, int(current.Int64()))
				current.Add(current, bigStep)
			}
		}

		result := lix.RangeWithSteps(start, stop, step).Take(limit).Collect()
		if !slices.Equal(result, want) {
			t.Fatalf("RangeWithSteps(%d, %d, %d) = %v, want %v", start, stop, step, result, want)
		}
	})
}

type stableItem struct {
	Key      byte
	Position int
}

// checkStableOrder 验证 result 按 Key 排序，且相等 Key 保持原始位置顺序。
func checkStableOrder(t *testing.T, name string, result []stableItem, size int, descending bool) {
	t.Helper()
	if len(result) != size {
		t.Fatalf("%s length = %d, want %d", name, len(result), size)
	}
	for index := 1; index < len(result); index++ {
		previous, current := result[index-1], result[index]
		outOfOrder := previous.Key > current.Key
		if descending {
			outOfOrder = previous.Key < current.Key
		}
		if outOfOrder {
			t.Fatalf("%s is not sorted at %d: %v", name, index, result)
		}
		if previous.Key == current.Key && previous.Position > current.Position {
			t.Fatalf("%s is not stable at %d: %v", name, index, result)
		}
	}
}

func FuzzSortVariantsAreStable(f *testing.F) {
	f.Add([]byte{2, 1, 2, 1})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 128 {
			data = data[:128]
		}
		items := make([]stableItem, len(data))
		pairs := make([]lix.Pair[byte, int], len(data))
		for index, key := range data {
			items[index] = stableItem{Key: key, Position: index}
			pairs[index] = lix.Pair[byte, int]{Key: key, Value: index}
		}
		byKey := func(item stableItem) byte { return item.Key }
		pairKey := func(key byte, _ int) byte { return key }
		toItems := func(s lix.Seq2[byte, int]) []stableItem {
			return lix.Pairs(s).Map(func(pair lix.Pair[byte, int]) stableItem {
				return stableItem{Key: pair.Key, Position: pair.Value}
			}).Collect()
		}

		values := lix.FromSlice(items)
		checkStableOrder(t, "Seq.SortByDescending", values.SortByDescending(byKey).Collect(), len(items), true)

		source := lix.FromPairs(pairs)
		checkStableOrder(t, "Seq2.SortBy", toItems(source.SortBy(pairKey)), len(items), false)
		checkStableOrder(t, "Seq2.SortByDescending", toItems(source.SortByDescending(pairKey)), len(items), true)
	})
}
