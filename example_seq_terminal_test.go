package lix_test

import (
	"fmt"
	"slices"

	"github.com/norniastar/lix"
)

func ExampleSeq_Iter() {
	result := slices.Collect(lix.Of(1, 2, 3).Iter())

	fmt.Println(result)
	// Output: [1 2 3]
}

func ExampleSeq_Collect() {
	result := lix.Of("a", "b").Collect()

	fmt.Println(result)
	// Output: [a b]
}

func ExampleSeq_First() {
	value, found := lix.Of("a", "b").First()

	fmt.Println(value, found)
	// Output: a true
}

func ExampleSeq_Last() {
	value, found := lix.Of("a", "b").Last()

	fmt.Println(value, found)
	// Output: b true
}

func ExampleSeq_Find() {
	value, found := lix.Of(1, 2, 3).
		Find(func(value int) bool { return value > 1 })

	fmt.Println(value, found)
	// Output: 2 true
}

func ExampleSeq_FindIndex() {
	index, found := lix.Of("a", "bb", "ccc").
		FindIndex(func(value string) bool { return len(value) == 2 })

	fmt.Println(index, found)
	// Output: 1 true
}

func ExampleSeq_SomeBy() {
	result := lix.Of(1, 2, 3).
		SomeBy(func(value int) bool { return value%2 == 0 })

	fmt.Println(result)
	// Output: true
}

func ExampleSeq_EveryBy() {
	result := lix.Of(1, 2, 3).
		EveryBy(func(value int) bool { return value > 0 })

	fmt.Println(result)
	// Output: true
}

func ExampleSeq_Length() {
	result := lix.Of("a", "b", "c").Length()

	fmt.Println(result)
	// Output: 3
}

func ExampleSeq_ForEach() {
	total := 0
	lix.Of(1, 2, 3).ForEach(func(value int) {
		total += value
	})

	fmt.Println(total)
	// Output: 6
}
