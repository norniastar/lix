package lix_test

import (
	"fmt"

	"github.com/norniastar/lix"
)

func ExampleSeq_Filter() {
	result := lix.Of(1, 2, 3, 4).
		Filter(func(value int) bool { return value%2 == 0 }).
		Collect()

	fmt.Println(result)
	// Output: [2 4]
}

func ExampleSeq_Reject() {
	result := lix.Of(1, 2, 3, 4).
		Reject(func(value int) bool { return value%2 == 0 }).
		Collect()

	fmt.Println(result)
	// Output: [1 3]
}

func ExampleSeq_FlatMap() {
	result := lix.Of(1, 2).
		FlatMap(func(value int) lix.Seq[int] {
			return lix.Of(value, value*10)
		}).
		Collect()

	fmt.Println(result)
	// Output: [1 10 2 20]
}

func ExampleSeq_FlatMap2() {
	result := lix.Pairs(
		lix.Of("a", "bb").FlatMap2(func(value string) lix.Seq2[int, string] {
			return lix.FromPairs([]lix.Pair[int, string]{
				{Key: len(value), Value: value},
			})
		}),
	).Collect()

	fmt.Println(result)
	// Output: [{1 a} {2 bb}]
}

func ExampleSeq_Take() {
	result := lix.Of(1, 2, 3, 4).Take(2).Collect()

	fmt.Println(result)
	// Output: [1 2]
}

func ExampleSeq_Drop() {
	result := lix.Of(1, 2, 3, 4).Drop(2).Collect()

	fmt.Println(result)
	// Output: [3 4]
}

func ExampleSeq_Enumerate() {
	result := lix.Pairs(lix.Of("a", "b").Enumerate()).Collect()

	fmt.Println(result)
	// Output: [{0 a} {1 b}]
}

func ExampleSeq_KeyBy() {
	result := lix.Of("go", "gopher").
		KeyBy(func(value string) int { return len(value) }).
		Collect()

	fmt.Println(result[2], result[6])
	// Output: go gopher
}

func ExampleSeq_Associate() {
	result := lix.Of("go", "gopher").
		Associate(func(value string) (int, string) {
			return len(value), value
		}).
		Collect()

	fmt.Println(result[2], result[6])
	// Output: go gopher
}
