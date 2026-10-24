package lix_test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/norniastar/lix"
)

func trackedValues(produced *int) lix.Seq[int] {
	return lix.FromSeq(func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3, 4} {
			*produced++
			if !yield(value) {
				return
			}
		}
	})
}

func trackedPairs(produced *int) lix.Seq2[string, int] {
	return lix.FromSeq2(func(yield func(string, int) bool) {
		for _, pair := range []lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "two", Value: 2},
			{Key: "three", Value: 3},
			{Key: "four", Value: 4},
		} {
			*produced++
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	})
}

// 每个用例只构造中间序列，并以方法值返回其遍历入口；
// 方法值会求值接收者，但不会调用遍历。
func TestSeqIntermediateOperationsDeferSourceConsumption(t *testing.T) {
	identity := func(value int) int { return value }
	tests := []struct {
		name      string
		construct func(lix.Seq[int]) func() int
	}{
		{name: "Map", construct: func(s lix.Seq[int]) func() int { return s.Map(identity).Length }},
		{name: "Filter", construct: func(s lix.Seq[int]) func() int { return s.Filter(func(int) bool { return true }).Length }},
		{name: "Reject", construct: func(s lix.Seq[int]) func() int { return s.Reject(func(int) bool { return false }).Length }},
		{name: "FlatMap", construct: func(s lix.Seq[int]) func() int {
			return s.FlatMap(func(value int) lix.Seq[int] { return lix.Of(value) }).Length
		}},
		{name: "FlatMap2", construct: func(s lix.Seq[int]) func() int {
			return s.FlatMap2(func(value int) lix.Seq2[int, int] {
				return lix.FromPairs([]lix.Pair[int, int]{{Key: value, Value: value}})
			}).Length
		}},
		{name: "Take", construct: func(s lix.Seq[int]) func() int { return s.Take(2).Length }},
		{name: "Drop", construct: func(s lix.Seq[int]) func() int { return s.Drop(1).Length }},
		{name: "Enumerate", construct: func(s lix.Seq[int]) func() int { return s.Enumerate().Length }},
		{name: "KeyBy", construct: func(s lix.Seq[int]) func() int { return s.KeyBy(identity).Length }},
		{name: "Associate", construct: func(s lix.Seq[int]) func() int {
			return s.Associate(func(value int) (int, int) { return value, value }).Length
		}},
		{name: "UniqBy", construct: func(s lix.Seq[int]) func() int { return s.UniqBy(identity).Length }},
		{name: "SortBy", construct: func(s lix.Seq[int]) func() int { return s.SortBy(identity).Length }},
		{name: "SortByDescending", construct: func(s lix.Seq[int]) func() int { return s.SortByDescending(identity).Length }},
		{name: "Reverse", construct: func(s lix.Seq[int]) func() int { return s.Reverse().Length }},
		{name: "GroupBy", construct: func(s lix.Seq[int]) func() int { return lix.GroupBy(s, identity).Length }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			produced := 0
			traverse := test.construct(trackedValues(&produced))
			if produced != 0 {
				t.Fatalf("construction consumed %d values, want 0", produced)
			}
			if traverse() == 0 || produced == 0 {
				t.Fatalf("traversal produced %d source values, want the source to be consumed", produced)
			}
		})
	}
}

func TestSeq2IntermediateOperationsDeferSourceConsumption(t *testing.T) {
	byValue := func(_ string, value int) int { return value }
	tests := []struct {
		name      string
		construct func(lix.Seq2[string, int]) func() int
	}{
		{name: "MapValues", construct: func(s lix.Seq2[string, int]) func() int { return s.MapValues(byValue).Length }},
		{name: "MapKeys", construct: func(s lix.Seq2[string, int]) func() int { return s.MapKeys(byValue).Length }},
		{name: "Filter", construct: func(s lix.Seq2[string, int]) func() int {
			return s.Filter(func(string, int) bool { return true }).Length
		}},
		{name: "Reject", construct: func(s lix.Seq2[string, int]) func() int {
			return s.Reject(func(string, int) bool { return false }).Length
		}},
		{name: "FlatMap", construct: func(s lix.Seq2[string, int]) func() int {
			return s.FlatMap(func(_ string, value int) lix.Seq[int] { return lix.Of(value) }).Length
		}},
		{name: "FlatMap2", construct: func(s lix.Seq2[string, int]) func() int {
			return s.FlatMap2(func(key string, value int) lix.Seq2[string, int] {
				return lix.FromPairs([]lix.Pair[string, int]{{Key: key, Value: value}})
			}).Length
		}},
		{name: "Take", construct: func(s lix.Seq2[string, int]) func() int { return s.Take(2).Length }},
		{name: "Drop", construct: func(s lix.Seq2[string, int]) func() int { return s.Drop(1).Length }},
		{name: "UniqBy", construct: func(s lix.Seq2[string, int]) func() int { return s.UniqBy(byValue).Length }},
		{name: "SortBy", construct: func(s lix.Seq2[string, int]) func() int { return s.SortBy(byValue).Length }},
		{name: "SortByDescending", construct: func(s lix.Seq2[string, int]) func() int {
			return s.SortByDescending(byValue).Length
		}},
		{name: "Reverse", construct: func(s lix.Seq2[string, int]) func() int { return s.Reverse().Length }},
		{name: "GroupPairsBy", construct: func(s lix.Seq2[string, int]) func() int {
			return lix.GroupPairsBy(s, byValue).Length
		}},
		{name: "Keys", construct: func(s lix.Seq2[string, int]) func() int { return s.Keys().Length }},
		{name: "Values", construct: func(s lix.Seq2[string, int]) func() int { return s.Values().Length }},
		{name: "Pairs", construct: func(s lix.Seq2[string, int]) func() int { return lix.Pairs(s).Length }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			produced := 0
			traverse := test.construct(trackedPairs(&produced))
			if produced != 0 {
				t.Fatalf("construction consumed %d pairs, want 0", produced)
			}
			if traverse() == 0 || produced == 0 {
				t.Fatalf("traversal produced %d source pairs, want the source to be consumed", produced)
			}
		})
	}
}

func replaySeq[V any](s lix.Seq[V]) func() string {
	return func() string { return fmt.Sprint(s.Collect()) }
}

func replaySeq2[K comparable, V any](s lix.Seq2[K, V]) func() string {
	return func() string { return fmt.Sprint(lix.Pairs(s).Collect()) }
}

// 带内部状态的派生序列必须在每次遍历时重新初始化状态。
func TestStatefulOperationsAreReplayable(t *testing.T) {
	values := lix.Of(3, 1, 2, 3, 1)
	pairs := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "c", Value: 3},
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c2", Value: 3},
		{Key: "a2", Value: 1},
	})
	identity := func(value int) int { return value }
	byValue := func(_ string, value int) int { return value }

	tests := []struct {
		name     string
		traverse func() string
	}{
		{name: "Seq/Take", traverse: replaySeq(values.Take(2))},
		{name: "Seq/Drop", traverse: replaySeq(values.Drop(2))},
		{name: "Seq/Enumerate", traverse: replaySeq2(values.Enumerate())},
		{name: "Seq/UniqBy", traverse: replaySeq(values.UniqBy(identity))},
		{name: "Seq/SortBy", traverse: replaySeq(values.SortBy(identity))},
		{name: "Seq/SortByDescending", traverse: replaySeq(values.SortByDescending(identity))},
		{name: "Seq/Reverse", traverse: replaySeq(values.Reverse())},
		{name: "Seq/GroupBy", traverse: replaySeq2(lix.GroupBy(values, identity))},
		{name: "Seq2/Take", traverse: replaySeq2(pairs.Take(2))},
		{name: "Seq2/Drop", traverse: replaySeq2(pairs.Drop(2))},
		{name: "Seq2/UniqBy", traverse: replaySeq2(pairs.UniqBy(byValue))},
		{name: "Seq2/SortBy", traverse: replaySeq2(pairs.SortBy(byValue))},
		{name: "Seq2/SortByDescending", traverse: replaySeq2(pairs.SortByDescending(byValue))},
		{name: "Seq2/Reverse", traverse: replaySeq2(pairs.Reverse())},
		{name: "Seq2/GroupPairsBy", traverse: replaySeq2(lix.GroupPairsBy(pairs, byValue))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			first := test.traverse()
			if first == "[]" {
				t.Fatal("first traversal is empty; the case does not exercise state")
			}
			if second := test.traverse(); second != first {
				t.Fatalf("second traversal = %s, want %s", second, first)
			}
		})
	}
}

func TestSeqStreamingOperationsPropagateStop(t *testing.T) {
	tests := []struct {
		name         string
		run          func(lix.Seq[int]) int
		wantProduced int
	}{
		{name: "Map", run: func(s lix.Seq[int]) int { return s.Map(func(value int) int { return value * 2 }).Take(1).Length() }, wantProduced: 1},
		{name: "Filter", run: func(s lix.Seq[int]) int { return s.Filter(func(int) bool { return true }).Take(1).Length() }, wantProduced: 1},
		{name: "Reject", run: func(s lix.Seq[int]) int { return s.Reject(func(int) bool { return false }).Take(1).Length() }, wantProduced: 1},
		{name: "FlatMap", run: func(s lix.Seq[int]) int {
			return s.FlatMap(func(value int) lix.Seq[int] { return lix.Of(value, -value) }).Take(1).Length()
		}, wantProduced: 1},
		{name: "FlatMap2", run: func(s lix.Seq[int]) int {
			return s.FlatMap2(func(value int) lix.Seq2[int, int] {
				return lix.FromPairs([]lix.Pair[int, int]{{Key: value, Value: value}})
			}).Take(1).Length()
		}, wantProduced: 1},
		{name: "Take", run: func(s lix.Seq[int]) int { return s.Take(4).Take(1).Length() }, wantProduced: 1},
		{name: "Drop", run: func(s lix.Seq[int]) int { return s.Drop(1).Take(1).Length() }, wantProduced: 2},
		{name: "Enumerate", run: func(s lix.Seq[int]) int { return s.Enumerate().Take(1).Length() }, wantProduced: 1},
		{name: "KeyBy", run: func(s lix.Seq[int]) int { return s.KeyBy(func(value int) int { return value }).Take(1).Length() }, wantProduced: 1},
		{name: "Associate", run: func(s lix.Seq[int]) int {
			return s.Associate(func(value int) (int, int) { return value, -value }).Take(1).Length()
		}, wantProduced: 1},
		{name: "UniqBy", run: func(s lix.Seq[int]) int { return s.UniqBy(func(value int) int { return value }).Take(1).Length() }, wantProduced: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			produced := 0
			if count := test.run(trackedValues(&produced)); count != 1 {
				t.Fatalf("downstream count = %d, want 1", count)
			}
			if produced != test.wantProduced {
				t.Fatalf("source produced %d values, want %d", produced, test.wantProduced)
			}
		})
	}
}

func TestSeq2StreamingOperationsPropagateStop(t *testing.T) {
	tests := []struct {
		name         string
		run          func(lix.Seq2[string, int]) int
		wantProduced int
	}{
		{name: "MapValues", run: func(s lix.Seq2[string, int]) int {
			return s.MapValues(func(_ string, value int) int { return value * 2 }).Take(1).Length()
		}, wantProduced: 1},
		{name: "MapKeys", run: func(s lix.Seq2[string, int]) int {
			return s.MapKeys(func(_ string, value int) int { return value }).Take(1).Length()
		}, wantProduced: 1},
		{name: "Filter", run: func(s lix.Seq2[string, int]) int {
			return s.Filter(func(string, int) bool { return true }).Take(1).Length()
		}, wantProduced: 1},
		{name: "Reject", run: func(s lix.Seq2[string, int]) int {
			return s.Reject(func(string, int) bool { return false }).Take(1).Length()
		}, wantProduced: 1},
		{name: "FlatMap", run: func(s lix.Seq2[string, int]) int {
			return s.FlatMap(func(_ string, value int) lix.Seq[int] { return lix.Of(value, -value) }).Take(1).Length()
		}, wantProduced: 1},
		{name: "FlatMap2", run: func(s lix.Seq2[string, int]) int {
			return s.FlatMap2(func(key string, value int) lix.Seq2[string, int] {
				return lix.FromPairs([]lix.Pair[string, int]{{Key: key, Value: value}})
			}).Take(1).Length()
		}, wantProduced: 1},
		{name: "Take", run: func(s lix.Seq2[string, int]) int { return s.Take(4).Take(1).Length() }, wantProduced: 1},
		{name: "Drop", run: func(s lix.Seq2[string, int]) int { return s.Drop(1).Take(1).Length() }, wantProduced: 2},
		{name: "UniqBy", run: func(s lix.Seq2[string, int]) int {
			return s.UniqBy(func(_ string, value int) int { return value }).Take(1).Length()
		}, wantProduced: 1},
		{name: "Keys", run: func(s lix.Seq2[string, int]) int { return s.Keys().Take(1).Length() }, wantProduced: 1},
		{name: "Values", run: func(s lix.Seq2[string, int]) int { return s.Values().Take(1).Length() }, wantProduced: 1},
		{name: "Pairs", run: func(s lix.Seq2[string, int]) int { return lix.Pairs(s).Take(1).Length() }, wantProduced: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			produced := 0
			if count := test.run(trackedPairs(&produced)); count != 1 {
				t.Fatalf("downstream count = %d, want 1", count)
			}
			if produced != test.wantProduced {
				t.Fatalf("source produced %d pairs, want %d", produced, test.wantProduced)
			}
		})
	}
}

func TestMaterializingOperationsStopOutput(t *testing.T) {
	seqTests := []struct {
		name string
		run  func(lix.Seq[int]) int
	}{
		{name: "SortBy", run: func(s lix.Seq[int]) int { return s.SortBy(func(value int) int { return value }).Take(1).Length() }},
		{name: "SortByDescending", run: func(s lix.Seq[int]) int {
			return s.SortByDescending(func(value int) int { return value }).Take(1).Length()
		}},
		{name: "Reverse", run: func(s lix.Seq[int]) int { return s.Reverse().Take(1).Length() }},
		{name: "GroupBy", run: func(s lix.Seq[int]) int {
			return lix.GroupBy(s, func(value int) int { return value % 2 }).Take(1).Length()
		}},
	}
	for _, test := range seqTests {
		t.Run("Seq/"+test.name, func(t *testing.T) {
			produced := 0
			if count := test.run(trackedValues(&produced)); count != 1 {
				t.Fatalf("downstream count = %d, want 1", count)
			}
			if produced != 4 {
				t.Fatalf("materializing operation consumed %d values, want 4", produced)
			}
		})
	}

	seq2Tests := []struct {
		name string
		run  func(lix.Seq2[string, int]) int
	}{
		{name: "SortBy", run: func(s lix.Seq2[string, int]) int {
			return s.SortBy(func(_ string, value int) int { return value }).Take(1).Length()
		}},
		{name: "SortByDescending", run: func(s lix.Seq2[string, int]) int {
			return s.SortByDescending(func(_ string, value int) int { return value }).Take(1).Length()
		}},
		{name: "Reverse", run: func(s lix.Seq2[string, int]) int { return s.Reverse().Take(1).Length() }},
		{name: "GroupPairsBy", run: func(s lix.Seq2[string, int]) int {
			return lix.GroupPairsBy(s, func(_ string, value int) int { return value % 2 }).Take(1).Length()
		}},
	}
	for _, test := range seq2Tests {
		t.Run("Seq2/"+test.name, func(t *testing.T) {
			produced := 0
			if count := test.run(trackedPairs(&produced)); count != 1 {
				t.Fatalf("downstream count = %d, want 1", count)
			}
			if produced != 4 {
				t.Fatalf("materializing operation consumed %d pairs, want 4", produced)
			}
		})
	}
}

func TestSeqTerminalsShortCircuit(t *testing.T) {
	tests := []struct {
		name         string
		run          func(lix.Seq[int]) bool
		wantProduced int
	}{
		{name: "First", run: func(s lix.Seq[int]) bool { value, ok := s.First(); return ok && value == 1 }, wantProduced: 1},
		{name: "Find", run: func(s lix.Seq[int]) bool {
			value, ok := s.Find(func(value int) bool { return value == 2 })
			return ok && value == 2
		}, wantProduced: 2},
		{name: "FindIndex", run: func(s lix.Seq[int]) bool {
			index, ok := s.FindIndex(func(value int) bool { return value == 2 })
			return ok && index == 1
		}, wantProduced: 2},
		{name: "SomeBy", run: func(s lix.Seq[int]) bool { return s.SomeBy(func(value int) bool { return value == 2 }) }, wantProduced: 2},
		{name: "EveryBy", run: func(s lix.Seq[int]) bool { return !s.EveryBy(func(value int) bool { return value < 3 }) }, wantProduced: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			produced := 0
			if !test.run(trackedValues(&produced)) {
				t.Fatal("terminal result did not match contract")
			}
			if produced != test.wantProduced {
				t.Fatalf("source produced %d values, want %d", produced, test.wantProduced)
			}
		})
	}
}

func TestSeq2TerminalsShortCircuit(t *testing.T) {
	tests := []struct {
		name         string
		run          func(lix.Seq2[string, int]) bool
		wantProduced int
	}{
		{name: "First", run: func(s lix.Seq2[string, int]) bool {
			key, value, ok := s.First()
			return ok && key == "one" && value == 1
		}, wantProduced: 1},
		{name: "Find", run: func(s lix.Seq2[string, int]) bool {
			key, value, ok := s.Find(func(_ string, value int) bool { return value == 2 })
			return ok && key == "two" && value == 2
		}, wantProduced: 2},
		{name: "FindIndex", run: func(s lix.Seq2[string, int]) bool {
			index, ok := s.FindIndex(func(_ string, value int) bool { return value == 2 })
			return ok && index == 1
		}, wantProduced: 2},
		{name: "SomeBy", run: func(s lix.Seq2[string, int]) bool {
			return s.SomeBy(func(_ string, value int) bool { return value == 2 })
		}, wantProduced: 2},
		{name: "EveryBy", run: func(s lix.Seq2[string, int]) bool {
			return !s.EveryBy(func(_ string, value int) bool { return value < 3 })
		}, wantProduced: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			produced := 0
			if !test.run(trackedPairs(&produced)) {
				t.Fatal("terminal result did not match contract")
			}
			if produced != test.wantProduced {
				t.Fatalf("source produced %d pairs, want %d", produced, test.wantProduced)
			}
		})
	}
}

func TestTakeAndDropBoundaries(t *testing.T) {
	takeTests := []struct {
		count int
		want  []int
	}{
		{count: -1, want: []int{}},
		{count: 0, want: []int{}},
		{count: 1, want: []int{1}},
		{count: 4, want: []int{1, 2, 3, 4}},
		{count: 8, want: []int{1, 2, 3, 4}},
	}
	for _, test := range takeTests {
		result := lix.Of(1, 2, 3, 4).Take(test.count).Collect()
		if !slices.Equal(result, test.want) {
			t.Fatalf("Take(%d) = %v, want %v", test.count, result, test.want)
		}
	}

	dropTests := []struct {
		count int
		want  []int
	}{
		{count: -1, want: []int{1, 2, 3, 4}},
		{count: 0, want: []int{1, 2, 3, 4}},
		{count: 1, want: []int{2, 3, 4}},
		{count: 4, want: []int{}},
		{count: 8, want: []int{}},
	}
	for _, test := range dropTests {
		result := lix.Of(1, 2, 3, 4).Drop(test.count).Collect()
		if !slices.Equal(result, test.want) {
			t.Fatalf("Drop(%d) = %v, want %v", test.count, result, test.want)
		}
	}
}

func TestSeq2TakeAndDropBoundaries(t *testing.T) {
	source := []lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
		{Key: "three", Value: 3},
	}

	takeTests := []struct {
		count int
		want  []lix.Pair[string, int]
	}{
		{count: -1, want: []lix.Pair[string, int]{}},
		{count: 0, want: []lix.Pair[string, int]{}},
		{count: 1, want: source[:1]},
		{count: 3, want: source},
		{count: 6, want: source},
	}
	for _, test := range takeTests {
		result := lix.Pairs(lix.FromPairs(source).Take(test.count)).Collect()
		if !slices.Equal(result, test.want) {
			t.Fatalf("Seq2.Take(%d) = %v, want %v", test.count, result, test.want)
		}
	}

	dropTests := []struct {
		count int
		want  []lix.Pair[string, int]
	}{
		{count: -1, want: source},
		{count: 0, want: source},
		{count: 1, want: source[1:]},
		{count: 3, want: []lix.Pair[string, int]{}},
		{count: 6, want: []lix.Pair[string, int]{}},
	}
	for _, test := range dropTests {
		result := lix.Pairs(lix.FromPairs(source).Drop(test.count)).Collect()
		if !slices.Equal(result, test.want) {
			t.Fatalf("Seq2.Drop(%d) = %v, want %v", test.count, result, test.want)
		}
	}
}

func TestEmptySequenceContracts(t *testing.T) {
	values := lix.Of[int]()
	if result := values.Collect(); result == nil || len(result) != 0 {
		t.Fatalf("empty Collect() = %#v, want non-nil empty slice", result)
	}
	if _, ok := values.First(); ok {
		t.Fatal("empty First() reported a value")
	}
	if _, ok := values.Last(); ok {
		t.Fatal("empty Last() reported a value")
	}
	if _, ok := values.Find(func(int) bool { return true }); ok {
		t.Fatal("empty Find() reported a value")
	}
	if index, ok := values.FindIndex(func(int) bool { return true }); ok || index != -1 {
		t.Fatalf("empty FindIndex() = (%d, %t), want (-1, false)", index, ok)
	}
	if values.SomeBy(func(int) bool { return true }) {
		t.Fatal("empty SomeBy() = true")
	}
	if !values.EveryBy(func(int) bool { return false }) {
		t.Fatal("empty EveryBy() = false")
	}
	if values.Length() != 0 {
		t.Fatal("empty Length() != 0")
	}
	if result := values.SortBy(func(value int) int { return value }).Collect(); result == nil || len(result) != 0 {
		t.Fatalf("empty SortBy() = %#v", result)
	}
	if result := values.Reverse().Collect(); result == nil || len(result) != 0 {
		t.Fatalf("empty Reverse() = %#v", result)
	}
	if result := lix.GroupBy(values, func(value int) int { return value }).Collect(); result == nil || len(result) != 0 {
		t.Fatalf("empty GroupBy() = %#v", result)
	}
}

func TestEmptyPairSequenceContracts(t *testing.T) {
	pairs := lix.FromPairs([]lix.Pair[string, int]{})
	if result := pairs.Collect(); result == nil || len(result) != 0 {
		t.Fatalf("empty Collect() = %#v, want non-nil empty map", result)
	}
	if _, _, ok := pairs.First(); ok {
		t.Fatal("empty First() reported a pair")
	}
	if _, _, ok := pairs.Last(); ok {
		t.Fatal("empty Last() reported a pair")
	}
	if _, _, ok := pairs.Find(func(string, int) bool { return true }); ok {
		t.Fatal("empty Find() reported a pair")
	}
	if index, ok := pairs.FindIndex(func(string, int) bool { return true }); ok || index != -1 {
		t.Fatalf("empty FindIndex() = (%d, %t), want (-1, false)", index, ok)
	}
	if pairs.SomeBy(func(string, int) bool { return true }) {
		t.Fatal("empty SomeBy() = true")
	}
	if !pairs.EveryBy(func(string, int) bool { return false }) {
		t.Fatal("empty EveryBy() = false")
	}
	if pairs.Length() != 0 {
		t.Fatal("empty Length() != 0")
	}
	if result := lix.GroupPairsBy(pairs, func(string, int) int { return 0 }).Collect(); result == nil || len(result) != 0 {
		t.Fatalf("empty GroupPairsBy() = %#v", result)
	}
}

func TestDuplicateKeysKeepLastValue(t *testing.T) {
	pairs := []lix.Pair[string, int]{
		{Key: "same", Value: 1},
		{Key: "same", Value: 2},
	}
	if result := lix.FromPairs(pairs).Collect(); !maps.Equal(result, map[string]int{"same": 2}) {
		t.Fatalf("FromPairs duplicate result = %v", result)
	}

	if result := lix.Of(1, 2).KeyBy(func(int) string { return "same" }).Collect(); !maps.Equal(result, map[string]int{"same": 2}) {
		t.Fatalf("KeyBy duplicate result = %v", result)
	}
	if result := lix.Of(1, 2).Associate(func(value int) (string, int) { return "same", value * 10 }).Collect(); !maps.Equal(result, map[string]int{"same": 20}) {
		t.Fatalf("Associate duplicate result = %v", result)
	}
	if result := lix.FromPairs(pairs).MapKeys(func(string, int) string { return "mapped" }).Collect(); !maps.Equal(result, map[string]int{"mapped": 2}) {
		t.Fatalf("MapKeys duplicate result = %v", result)
	}
}
