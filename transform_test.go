package lix_test

import (
	"maps"
	"slices"
	"strconv"
	"testing"

	"github.com/norniastar/lix"
)

func TestMapChangesValueTypeLazily(t *testing.T) {
	calls := 0
	sequence := lix.FromSlice([]int{1, 2, 3}).Map(func(value int) string {
		calls++
		return strconv.Itoa(value)
	})

	if calls != 0 {
		t.Fatalf("Map evaluated eagerly: calls = %d", calls)
	}

	want := []string{"1", "2", "3"}
	if result := sequence.Collect(); !slices.Equal(result, want) {
		t.Fatalf("Collect() = %v, want %v", result, want)
	}
	if calls != len(want) {
		t.Fatalf("mapper calls = %d, want %d", calls, len(want))
	}
}

func TestFilterFlatMapDropAndTake(t *testing.T) {
	result := lix.RangeWithSteps(1, 8, 1).
		Filter(func(value int) bool {
			return value%2 != 0
		}).
		FlatMap(func(value int) lix.Seq[int] {
			return lix.Of(value, -value)
		}).
		Drop(1).
		Take(4).
		Collect()

	want := []int{-1, 3, -3, 5}
	if !slices.Equal(result, want) {
		t.Fatalf("pipeline result = %v, want %v", result, want)
	}
}

func TestRejectEnumerateKeyByAndAssociate(t *testing.T) {
	type user struct {
		ID   int
		Name string
	}

	users := lix.Of(
		user{ID: 1, Name: "one"},
		user{ID: 2, Name: "two"},
		user{ID: 3, Name: "three"},
	)

	indexed := users.
		Reject(func(value user) bool { return value.ID == 2 }).
		KeyBy(func(value user) int { return value.ID }).
		Map(func(_ int, value user) string { return value.Name }).
		Collect()
	if !maps.Equal(indexed, map[int]string{1: "one", 3: "three"}) {
		t.Fatalf("indexed result = %v", indexed)
	}

	enumerated := users.Enumerate().Collect()
	if !maps.Equal(enumerated, map[int]user{0: {ID: 1, Name: "one"}, 1: {ID: 2, Name: "two"}, 2: {ID: 3, Name: "three"}}) {
		t.Fatalf("Enumerate() = %v", enumerated)
	}

	associated := users.Associate(func(value user) (string, int) {
		return value.Name, value.ID * 10
	}).Collect()
	if !maps.Equal(associated, map[string]int{"one": 10, "two": 20, "three": 30}) {
		t.Fatalf("Associate() = %v", associated)
	}
}

func TestFlatMap2PreservesMappedKeys(t *testing.T) {
	sequence := lix.Of("a", "b").
		FlatMap2(func(value string) lix.Seq2[int, string] {
			return lix.FromPairs([]lix.Pair[int, string]{
				{Key: 0, Value: value + "0"},
				{Key: 1, Value: value + "1"},
			})
		})
	result := lix.Pairs(sequence).Collect()

	want := []lix.Pair[int, string]{
		{Key: 0, Value: "a0"},
		{Key: 1, Value: "a1"},
		{Key: 0, Value: "b0"},
		{Key: 1, Value: "b1"},
	}
	if !slices.Equal(result, want) {
		t.Fatalf("FlatMap2() = %v, want %v", result, want)
	}
}

func TestTakeStopsSourceEarly(t *testing.T) {
	produced := 0
	source := lix.FromSeq(func(yield func(int) bool) {
		for value := range 100 {
			produced++
			if !yield(value) {
				return
			}
		}
	})

	result := source.Take(3).Collect()
	if !slices.Equal(result, []int{0, 1, 2}) {
		t.Fatalf("Take(3) = %v, want [0 1 2]", result)
	}
	if produced != 3 {
		t.Fatalf("source produced %d values, want early stop after 3", produced)
	}
}

func TestSeq2Transformations(t *testing.T) {
	result := lix.FromMap(map[string]int{"a": 1, "b": 2, "c": 3}).
		Filter(func(_ string, value int) bool {
			return value%2 != 0
		}).
		Map(func(key string, value int) string {
			return key + strconv.Itoa(value)
		}).
		MapKeys(func(key string, _ string) string {
			return "key-" + key
		}).
		Collect()

	want := map[string]string{"key-a": "a1", "key-c": "c3"}
	if !maps.Equal(result, want) {
		t.Fatalf("Seq2 pipeline = %v, want %v", result, want)
	}
}

func TestSeq2RejectUniqByAndFlatMap2(t *testing.T) {
	sequence := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
		{Key: "d", Value: 4},
	})

	result := sequence.
		Reject(func(_ string, value int) bool { return value == 2 }).
		UniqBy(func(_ string, value int) int { return value % 2 }).
		FlatMap2(func(key string, value int) lix.Seq2[string, int] {
			return lix.FromPairs([]lix.Pair[string, int]{{Key: "mapped-" + key, Value: value * 10}})
		})

	want := []lix.Pair[string, int]{{Key: "mapped-a", Value: 10}, {Key: "mapped-d", Value: 40}}
	if collected := lix.Pairs(result).Collect(); !slices.Equal(collected, want) {
		t.Fatalf("Seq2 transformations = %v, want %v", collected, want)
	}
}

func TestSeq2FlatMapDropAndTake(t *testing.T) {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
		{Key: "three", Value: 3},
	}).
		Drop(1).
		Take(1).
		FlatMap(func(key string, value int) lix.Seq[string] {
			return lix.Of(key, strconv.Itoa(value))
		}).
		Collect()

	if !slices.Equal(result, []string{"two", "2"}) {
		t.Fatalf("Seq2 FlatMap/Drop/Take = %v", result)
	}
}
