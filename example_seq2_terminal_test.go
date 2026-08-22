package lix_test

import (
	"fmt"

	"github.com/norniastar/lix"
)

func ExampleSeq2_Iter() {
	sequence := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	})

	for key, value := range sequence.Iter() {
		fmt.Printf("%s=%d\n", key, value)
	}
	// Output:
	// one=1
	// two=2
}

func ExampleSeq2_Collect() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).Collect()

	fmt.Println(result["one"], result["two"])
	// Output: 1 2
}

func ExampleSeq2_First() {
	key, value, found := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).First()

	fmt.Println(key, value, found)
	// Output: one 1 true
}

func ExampleSeq2_Last() {
	key, value, found := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).Last()

	fmt.Println(key, value, found)
	// Output: two 2 true
}

func ExampleSeq2_Find() {
	key, value, found := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).Find(func(_ string, value int) bool {
		return value%2 == 0
	})

	fmt.Println(key, value, found)
	// Output: two 2 true
}

func ExampleSeq2_FindIndex() {
	index, found := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).FindIndex(func(key string, _ int) bool {
		return key == "two"
	})

	fmt.Println(index, found)
	// Output: 1 true
}

func ExampleSeq2_SomeBy() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).SomeBy(func(_ string, value int) bool {
		return value%2 == 0
	})

	fmt.Println(result)
	// Output: true
}

func ExampleSeq2_EveryBy() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).EveryBy(func(key string, _ int) bool {
		return len(key) > 0
	})

	fmt.Println(result)
	// Output: true
}

func ExampleSeq2_Length() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).Length()

	fmt.Println(result)
	// Output: 2
}

func ExampleSeq2_ForEach() {
	total := 0
	lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).ForEach(func(_ string, value int) {
		total += value
	})

	fmt.Println(total)
	// Output: 3
}

func ExampleSeq2_Keys() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).Keys().Collect()

	fmt.Println(result)
	// Output: [one two]
}

func ExampleSeq2_Values() {
	result := lix.FromPairs([]lix.Pair[string, int]{
		{Key: "one", Value: 1},
		{Key: "two", Value: 2},
	}).Values().Collect()

	fmt.Println(result)
	// Output: [1 2]
}
