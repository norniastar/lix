package lix_test

import (
	"slices"
	"testing"

	"github.com/norniastar/lix"
)

func TestReduceChangesResultType(t *testing.T) {
	result := lix.Of(1, 2, 3).Reduce(func(result string, value int) string {
		return result + " " + string(rune('0'+value))
	}, "values:")

	if result != "values: 1 2 3" {
		t.Fatalf("Reduce() = %q, want %q", result, "values: 1 2 3")
	}
}

func TestGroupBy(t *testing.T) {
	groups := lix.GroupBy(lix.RangeWithSteps(1, 6, 1), func(value int) string {
		if value%2 == 0 {
			return "even"
		}
		return "odd"
	})
	result := groups.Collect()

	if !slices.Equal(result["odd"], []int{1, 3, 5}) {
		t.Fatalf("odd group = %v, want [1 3 5]", result["odd"])
	}
	if !slices.Equal(result["even"], []int{2, 4}) {
		t.Fatalf("even group = %v, want [2 4]", result["even"])
	}
}

func TestGroupByIsLazyAndKeepsFirstKeyOrder(t *testing.T) {
	calls := 0
	groups := lix.GroupBy(lix.Of(1, 2, 3, 4), func(value int) string {
		calls++
		if value%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if calls != 0 {
		t.Fatalf("GroupBy evaluated eagerly: calls = %d", calls)
	}

	pairs := lix.Pairs(groups).Collect()
	want := []lix.Pair[string, []int]{
		{Key: "odd", Value: []int{1, 3}},
		{Key: "even", Value: []int{2, 4}},
	}
	if len(pairs) != len(want) {
		t.Fatalf("group count = %d, want %d", len(pairs), len(want))
	}
	for index := range pairs {
		if pairs[index].Key != want[index].Key || !slices.Equal(pairs[index].Value, want[index].Value) {
			t.Fatalf("group %d = %v, want %v", index, pairs[index], want[index])
		}
	}
	if calls != 4 {
		t.Fatalf("selector calls = %d, want 4", calls)
	}
}

func TestUniqByKeepsFirstValue(t *testing.T) {
	type user struct {
		ID   int
		Name string
	}

	result := lix.Of(
		user{ID: 1, Name: "first"},
		user{ID: 1, Name: "duplicate"},
		user{ID: 2, Name: "second"},
	).UniqBy(func(value user) int {
		return value.ID
	}).Collect()

	want := []user{{ID: 1, Name: "first"}, {ID: 2, Name: "second"}}
	if !slices.Equal(result, want) {
		t.Fatalf("UniqBy() = %v, want %v", result, want)
	}
}

func TestSeq2ReduceAndGroupBy(t *testing.T) {
	sequence := lix.FromMap(map[string]int{"one": 1, "two": 2, "three": 3})
	total := sequence.Reduce(func(result int, _ string, value int) int {
		return result + value
	}, 0)
	if total != 6 {
		t.Fatalf("Reduce() = %d, want 6", total)
	}

	groups := lix.GroupPairsBy(sequence, func(_ string, value int) string {
		if value%2 == 0 {
			return "even"
		}
		return "odd"
	}).Collect()
	if !slices.Equal(groups["even"], []int{2}) {
		t.Fatalf("even group = %v, want [2]", groups["even"])
	}
	if len(groups["odd"]) != 2 {
		t.Fatalf("odd group length = %d, want 2", len(groups["odd"]))
	}
}
