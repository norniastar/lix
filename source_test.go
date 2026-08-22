package lix_test

import (
	"maps"
	"math"
	"slices"
	"testing"

	"github.com/norniastar/lix"
)

func TestFromSlice(t *testing.T) {
	source := []int{1, 2, 3}
	sequence := lix.FromSlice(source)

	if result := sequence.Collect(); !slices.Equal(result, source) {
		t.Fatalf("Collect() = %v, want %v", result, source)
	}
	if result := sequence.Collect(); !slices.Equal(result, source) {
		t.Fatalf("second Collect() = %v, want replayable sequence %v", result, source)
	}
}

func TestFromMap(t *testing.T) {
	source := map[string]int{"one": 1, "two": 2}
	result := lix.FromMap(source).Collect()

	if !maps.Equal(result, source) {
		t.Fatalf("Collect() = %v, want %v", result, source)
	}
}

func TestFromPairsAndPairs(t *testing.T) {
	source := []lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}
	sequence := lix.FromPairs(source)

	if result := sequence.Collect(); !maps.Equal(result, map[string]int{"one": 1, "two": 2}) {
		t.Fatalf("Collect() = %v, want source pairs", result)
	}
	if result := lix.Pairs(sequence).Collect(); !slices.Equal(result, source) {
		t.Fatalf("Pairs() = %v, want %v", result, source)
	}
}

func TestIteratorInterop(t *testing.T) {
	values := lix.FromSeq(slices.Values([]int{1, 2, 3}))
	if result := slices.Collect(values.Iter()); !slices.Equal(result, []int{1, 2, 3}) {
		t.Fatalf("Seq.Iter() = %v", result)
	}

	source := map[string]int{"one": 1, "two": 2}
	pairs := lix.FromSeq2(maps.All(source))
	if result := maps.Collect(pairs.Iter()); !maps.Equal(result, source) {
		t.Fatalf("Seq2.Iter() = %v", result)
	}
}

func TestRangeWithSteps(t *testing.T) {
	tests := []struct {
		name        string
		start, stop int
		step        int
		want        []int
	}{
		{name: "ascending", start: 1, stop: 7, step: 2, want: []int{1, 3, 5}},
		{name: "descending", start: 5, stop: -1, step: -2, want: []int{5, 3, 1}},
		{name: "empty ascending", start: 3, stop: 1, step: 1, want: []int{}},
		{name: "empty descending", start: 1, stop: 3, step: -1, want: []int{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := lix.RangeWithSteps(test.start, test.stop, test.step).Collect()
			if !slices.Equal(result, test.want) {
				t.Fatalf("RangeWithSteps(%d, %d, %d) = %v, want %v", test.start, test.stop, test.step, result, test.want)
			}
		})
	}
}

func TestRangeWithStepsReturnsEmptyForZeroStep(t *testing.T) {
	result := lix.RangeWithSteps(0, 10, 0).Collect()
	if len(result) != 0 {
		t.Fatalf("RangeWithSteps with a zero step = %v, want empty sequence", result)
	}
}

func TestRangeWithStepsStopsOnIntegerOverflow(t *testing.T) {
	ascending := lix.RangeWithSteps(math.MaxInt-1, math.MaxInt, 2).Collect()
	if !slices.Equal(ascending, []int{math.MaxInt - 1}) {
		t.Fatalf("ascending overflow range = %v", ascending)
	}

	descending := lix.RangeWithSteps(math.MinInt+1, math.MinInt, -2).Collect()
	if !slices.Equal(descending, []int{math.MinInt + 1}) {
		t.Fatalf("descending overflow range = %v", descending)
	}
}
