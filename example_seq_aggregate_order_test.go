package lix_test

import (
	"fmt"

	"github.com/norniastar/lix"
)

func ExampleSeq_UniqBy() {
	result := lix.Of(1, 3, 2, 4, 5).
		UniqBy(func(value int) int { return value % 2 }).
		Collect()

	fmt.Println(result)
	// Output: [1 2]
}

func ExampleSeq_SortBy() {
	result := lix.Of(3, 1, 2).
		SortBy(func(value int) int { return value }).
		Collect()

	fmt.Println(result)
	// Output: [1 2 3]
}

func ExampleSeq_SortByDescending() {
	result := lix.Of(3, 1, 2).
		SortByDescending(func(value int) int { return value }).
		Collect()

	fmt.Println(result)
	// Output: [3 2 1]
}

func ExampleSeq_Reverse() {
	result := lix.Of(1, 2, 3).Reverse().Collect()

	fmt.Println(result)
	// Output: [3 2 1]
}
