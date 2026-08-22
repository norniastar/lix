package lix_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/norniastar/lix"
)

func TestSeqTerminals(t *testing.T) {
	sequence := lix.Of(1, 2, 3)

	first, ok := sequence.First()
	if !ok || first != 1 {
		t.Fatalf("First() = (%d, %t), want (1, true)", first, ok)
	}
	if !sequence.SomeBy(func(value int) bool { return value == 2 }) {
		t.Fatal("SomeBy() = false, want true")
	}
	if !sequence.EveryBy(func(value int) bool { return value > 0 }) {
		t.Fatal("EveryBy() = false, want true")
	}
	if count := sequence.Length(); count != 3 {
		t.Fatalf("Length() = %d, want 3", count)
	}

	if _, ok := lix.Of[int]().First(); ok {
		t.Fatal("First() on an empty sequence reported a value")
	}
}

func TestSeqSearchAndForEach(t *testing.T) {
	sequence := lix.Of(3, 1, 4, 1, 5)

	last, ok := sequence.Last()
	if !ok || last != 5 {
		t.Fatalf("Last() = (%d, %t), want (5, true)", last, ok)
	}
	found, ok := sequence.Find(func(value int) bool { return value%2 == 0 })
	if !ok || found != 4 {
		t.Fatalf("Find() = (%d, %t), want (4, true)", found, ok)
	}
	index, ok := sequence.FindIndex(func(value int) bool { return value == 4 })
	if !ok || index != 2 {
		t.Fatalf("FindIndex() = (%d, %t), want (2, true)", index, ok)
	}
	if index, ok := sequence.FindIndex(func(value int) bool { return value == 9 }); ok || index != -1 {
		t.Fatalf("missing FindIndex() = (%d, %t), want (-1, false)", index, ok)
	}

	visited := make([]int, 0)
	sequence.ForEach(func(value int) { visited = append(visited, value) })
	if !slices.Equal(visited, []int{3, 1, 4, 1, 5}) {
		t.Fatalf("ForEach() visited %v", visited)
	}

	if _, ok := lix.Of[int]().Last(); ok {
		t.Fatal("Last() on an empty sequence reported a value")
	}
}

func TestSeq2Terminals(t *testing.T) {
	source := map[string]int{"one": 1, "two": 2}
	sequence := lix.FromMap(source)

	if result := sequence.Collect(); !maps.Equal(result, source) {
		t.Fatalf("Collect() = %v, want %v", result, source)
	}
	if !sequence.SomeBy(func(_ string, value int) bool { return value == 2 }) {
		t.Fatal("SomeBy() = false, want true")
	}
	if !sequence.EveryBy(func(key string, _ int) bool { return key != "" }) {
		t.Fatal("EveryBy() = false, want true")
	}
	if count := sequence.Length(); count != len(source) {
		t.Fatalf("Length() = %d, want %d", count, len(source))
	}

	keys := sequence.Keys().SortBy(func(value string) string { return value }).Collect()
	if !slices.Equal(keys, []string{"one", "two"}) {
		t.Fatalf("Keys() = %v, want [one two]", keys)
	}
	values := sequence.Values().SortBy(func(value int) int { return value }).Collect()
	if !slices.Equal(values, []int{1, 2}) {
		t.Fatalf("Values() = %v, want [1 2]", values)
	}
}

func TestSeq2SearchAndForEach(t *testing.T) {
	sequence := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
		{Key: "three", Value: 3},
	})
	firstKey, firstValue, ok := sequence.First()
	if !ok || firstKey != "one" || firstValue != 1 {
		t.Fatalf("First() = (%q, %d, %t)", firstKey, firstValue, ok)
	}

	key, value, ok := sequence.Find(func(_ string, value int) bool { return value == 2 })
	if !ok || key != "two" || value != 2 {
		t.Fatalf("Find() = (%q, %d, %t)", key, value, ok)
	}
	index, ok := sequence.FindIndex(func(key string, _ int) bool { return key == "three" })
	if !ok || index != 2 {
		t.Fatalf("FindIndex() = (%d, %t), want (2, true)", index, ok)
	}
	key, value, ok = sequence.Last()
	if !ok || key != "three" || value != 3 {
		t.Fatalf("Last() = (%q, %d, %t)", key, value, ok)
	}

	visited := make([]string, 0)
	sequence.ForEach(func(key string, _ int) { visited = append(visited, key) })
	if !slices.Equal(visited, []string{"one", "two", "three"}) {
		t.Fatalf("ForEach() visited %v", visited)
	}

	if _, _, ok := lix.FromPairs([]lix.Pair[string, int]{}).First(); ok {
		t.Fatal("First() on an empty pair sequence reported a value")
	}
	if _, _, ok := sequence.Find(func(_ string, value int) bool { return value == 9 }); ok {
		t.Fatal("Find() reported a missing pair")
	}
}
