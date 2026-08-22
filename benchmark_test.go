package lix_test

import (
	"testing"

	"github.com/norniastar/lix"
)

var benchmarkResult []int

func BenchmarkPipeline(b *testing.B) {
	source := make([]int, 1_000)
	for index := range source {
		source[index] = index
	}

	for b.Loop() {
		benchmarkResult = lix.FromSlice(source).
			Filter(func(value int) bool {
				return value%2 == 0
			}).
			Map(func(value int) int {
				return value * value
			}).
			Take(100).
			Collect()
	}
}
