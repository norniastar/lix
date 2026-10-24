package lix_test

import (
	"fmt"

	"github.com/norniastar/lix"
)

func ExampleSeq2_MapValues() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).MapValues(func(_ string, value int) string {
		return fmt.Sprintf("value-%d", value)
	}).Collect()

	fmt.Println(result["one"], result["two"])
	// Output: value-1 value-2
}

func ExampleSeq2_MapKeys() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "a", Value: 1},
			{Key: "three", Value: 3},
		}).MapKeys(func(key string, _ int) int {
			return len(key)
		}),
	).Collect()

	fmt.Println(result)
	// Output: [{1 1} {5 3}]
}

func ExampleSeq2_Filter() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "two", Value: 2},
			{Key: "three", Value: 3},
		}).Filter(func(_ string, value int) bool {
			return value%2 == 0
		}),
	).Collect()

	fmt.Println(result)
	// Output: [{two 2}]
}

func ExampleSeq2_Reject() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "two", Value: 2},
			{Key: "three", Value: 3},
		}).Reject(func(_ string, value int) bool {
			return value%2 == 0
		}),
	).Collect()

	fmt.Println(result)
	// Output: [{one 1} {three 3}]
}

func ExampleSeq2_FlatMap() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).FlatMap(func(_ string, value int) lix.Seq[int] {
		return lix.Of(value, value*10)
	}).Collect()

	fmt.Println(result)
	// Output: [1 10 2 20]
}

func ExampleSeq2_FlatMap2() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
		}).FlatMap2(func(key string, value int) lix.Seq2[int, string] {
			return lix.FromPairs([]lix.Pair[int, string]{
				{Key: value, Value: key},
			})
		}),
	).Collect()

	fmt.Println(result)
	// Output: [{1 a} {2 b}]
}

func ExampleSeq2_Take() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "two", Value: 2},
			{Key: "three", Value: 3},
		}).Take(2),
	).Collect()

	fmt.Println(result)
	// Output: [{one 1} {two 2}]
}

func ExampleSeq2_Drop() {
	result := lix.Pairs(
		lix.FromPairs([]lix.Pair[string, int]{
			{Key: "one", Value: 1},
			{Key: "two", Value: 2},
			{Key: "three", Value: 3},
		}).Drop(2),
	).Collect()

	fmt.Println(result)
	// Output: [{three 3}]
}
