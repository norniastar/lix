package lix_test

import (
	"fmt"
	"strconv"

	"github.com/norniastar/lix"
)

func ExampleSeq_Map() {
	result := lix.FromSlice([]int{1, 2, 3, 4}).
		Filter(func(value int) bool {
			return value%2 == 0
		}).
		Map(func(value int) string {
			return "value-" + strconv.Itoa(value)
		}).
		Collect()

	fmt.Println(result)
	// Output: [value-2 value-4]
}

func ExampleSeq_Reduce() {
	total := lix.RangeWithSteps(1, 5, 1).
		Reduce(func(result, value int) int {
			return result + value
		}, 0)

	fmt.Println(total)
	// Output: 10
}

func ExampleGroupBy() {
	groups := lix.GroupBy(
		lix.RangeWithSteps(1, 6, 1),
		func(value int) string {
			if value%2 == 0 {
				return "even"
			}
			return "odd"
		},
	).Collect()

	fmt.Println(groups["odd"])
	fmt.Println(groups["even"])
	// Output:
	// [1 3 5]
	// [2 4]
}
