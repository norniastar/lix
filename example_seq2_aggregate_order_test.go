package lix_test

import (
	"fmt"

	"github.com/norniastar/lix"
)

func ExampleSeq2_UniqBy() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "another-odd", Value: 3},
			{Key: "two", Value: 2},
		}).UniqBy(func(_ string, value int) int {
			return value % 2
		}),
	).Collect()

	fmt.Println(result)
	// Output: [{one 1} {two 2}]
}

func ExampleSeq2_Reduce() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).Reduce(func(total int, _ string, value int) int {
		return total + value
	}, 0)

	fmt.Println(result)
	// Output: 3
}

func ExampleSeq2_SortBy() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "three", Value: 3},
			{Key: "one", Value: 1},
			{Key: "two", Value: 2},
		}).SortBy(func(_ string, value int) int {
			return value
		}),
	).Collect()

	fmt.Println(result)
	// Output: [{one 1} {two 2} {three 3}]
}

func ExampleSeq2_SortByDescending() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "three", Value: 3},
			{Key: "two", Value: 2},
		}).SortByDescending(func(_ string, value int) int {
			return value
		}),
	).Collect()

	fmt.Println(result)
	// Output: [{three 3} {two 2} {one 1}]
}

func ExampleSeq2_Reverse() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "two", Value: 2},
			{Key: "three", Value: 3},
		}).Reverse(),
	).Collect()

	fmt.Println(result)
	// Output: [{three 3} {two 2} {one 1}]
}
