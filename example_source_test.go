package lix_test

import (
	"fmt"
	"iter"
	"slices"

	"github.com/norniastar/lix"
)

func ExampleFromSlice() {
	result := lix.FromSlice([]int{1, 2, 3}).Collect()

	fmt.Println(result)
	// Output: [1 2 3]
}

func ExampleFromMap() {
	result := lix.FromMap(map[string]int{"one": 1, "two": 2}).Collect()

	fmt.Println(result["one"], result["two"])
	// Output: 1 2
}

func ExampleFromSeq() {
	source := slices.Values([]string{"a", "b"})
	result := lix.FromSeq(source).Collect()

	fmt.Println(result)
	// Output: [a b]
}

func ExampleFromSeq2() {
	source := iter.Seq2[string, int](func(yield func(string, int) bool) {
		if !yield("one", 1) {
			return
		}
		yield("two", 2)
	})
	result := lix.FromSeq2(source).Collect()

	fmt.Println(result["one"], result["two"])
	// Output: 1 2
}

func ExampleFromPairs() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).Collect()

	fmt.Println(result["one"], result["two"])
	// Output: 1 2
}

func ExampleOf() {
	result := lix.Of("a", "b", "c").Collect()

	fmt.Println(result)
	// Output: [a b c]
}

func ExampleRangeWithSteps() {
	result := lix.RangeWithSteps(1, 7, 2).Collect()

	fmt.Println(result)
	// Output: [1 3 5]
}

func ExampleGroupPairsBy() {
	groups := lix.GroupPairsBy(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "two", Value: 2},
			{Key: "three", Value: 3},
		}),
		func(_ string, value int) string {
			if value%2 == 0 {
				return "even"
			}
			return "odd"
		},
	).Collect()

	fmt.Println(groups["odd"])
	fmt.Println(groups["even"])
	// Output:
	// [1 3]
	// [2]
}

func ExamplePairs() {
	result := lix.Pairs(lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	})).Collect()

	fmt.Println(result)
	// Output: [{one 1} {two 2}]
}
