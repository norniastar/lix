package lix_test

import (
	"slices"
	"testing"

	"github.com/norniastar/lix"
)

func TestSortByIsStableAndLazy(t *testing.T) {
	type item struct {
		Group int
		Name  string
	}

	calls := 0
	sequence := lix.Of(
		item{Group: 2, Name: "second-a"},
		item{Group: 1, Name: "first"},
		item{Group: 2, Name: "second-b"},
	).SortBy(func(value item) int {
		calls++
		return value.Group
	})

	if calls != 0 {
		t.Fatalf("SortBy evaluated eagerly: calls = %d", calls)
	}

	result := sequence.Collect()
	want := []item{
		{Group: 1, Name: "first"},
		{Group: 2, Name: "second-a"},
		{Group: 2, Name: "second-b"},
	}
	if !slices.Equal(result, want) {
		t.Fatalf("SortBy() = %v, want %v", result, want)
	}
	if calls != len(want) {
		t.Fatalf("selector calls = %d, want %d", calls, len(want))
	}
}

func TestSortByDescendingAndReverse(t *testing.T) {
	type item struct {
		Rank int
		Name string
	}
	source := lix.Of(
		item{Rank: 1, Name: "first-a"},
		item{Rank: 2, Name: "second"},
		item{Rank: 1, Name: "first-b"},
	)

	descending := source.SortByDescending(func(value item) int { return value.Rank }).Collect()
	wantDescending := []item{
		{Rank: 2, Name: "second"},
		{Rank: 1, Name: "first-a"},
		{Rank: 1, Name: "first-b"},
	}
	if !slices.Equal(descending, wantDescending) {
		t.Fatalf("SortByDescending() = %v, want %v", descending, wantDescending)
	}

	reversed := source.Reverse().Collect()
	wantReversed := []item{
		{Rank: 1, Name: "first-b"},
		{Rank: 2, Name: "second"},
		{Rank: 1, Name: "first-a"},
	}
	if !slices.Equal(reversed, wantReversed) {
		t.Fatalf("Reverse() = %v, want %v", reversed, wantReversed)
	}
}

func TestSeq2Ordering(t *testing.T) {
	source := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "three", Value: 3},
		{Key: "two", Value: 2},
	})

	descending := source.SortByDescending(func(_ string, value int) int { return value })
	wantDescending := []lix.Pair[string, int]{
		{Key: "three", Value: 3},
		{Key: "two", Value: 2},
		{Key: "one", Value: 1},
	}
	if result := lix.Pairs(descending).Collect(); !slices.Equal(result, wantDescending) {
		t.Fatalf("SortByDescending() = %v, want %v", result, wantDescending)
	}

	ascending := source.SortBy(func(_ string, value int) int { return value })
	wantAscending := []lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
		{Key: "three", Value: 3},
	}
	if result := lix.Pairs(ascending).Collect(); !slices.Equal(result, wantAscending) {
		t.Fatalf("SortBy() = %v, want %v", result, wantAscending)
	}

	reversed := source.Reverse()
	wantReversed := []lix.Pair[string, int]{
		{Key: "two", Value: 2},
		{Key: "three", Value: 3},
		{Key: "one", Value: 1},
	}
	if result := lix.Pairs(reversed).Collect(); !slices.Equal(result, wantReversed) {
		t.Fatalf("Reverse() = %v, want %v", result, wantReversed)
	}
}
