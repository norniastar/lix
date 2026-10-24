package lix_test

import (
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/norniastar/lix"
)

const benchmarkSize = 10_000

var (
	benchmarkResult []int
	benchmarkCount  int
)

// benchmarkValues 返回固定种子的伪随机数据，使不同运行之间可比较。
func benchmarkValues() []int {
	random := rand.New(rand.NewPCG(1, 2))
	values := make([]int, benchmarkSize)
	for index := range values {
		values[index] = random.IntN(benchmarkSize)
	}
	return values
}

// benchmarkCardinalities 覆盖键大量重复和几乎全部唯一两种情况，
// 分别对应 map 查找以命中为主和以插入为主的路径。
var benchmarkCardinalities = []int{16, benchmarkSize}

func BenchmarkPipeline(b *testing.B) {
	source := make([]int, 1_000)
	for index := range source {
		source[index] = index
	}

	b.ReportAllocs()
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

// BenchmarkPipelineHandwritten 是 BenchmarkPipeline 的手写循环基线，
// 用于量化迭代器组合本身的开销。
func BenchmarkPipelineHandwritten(b *testing.B) {
	source := make([]int, 1_000)
	for index := range source {
		source[index] = index
	}

	b.ReportAllocs()
	for b.Loop() {
		result := make([]int, 0)
		for _, value := range source {
			if value%2 != 0 {
				continue
			}
			result = append(result, value*value)
			if len(result) == 100 {
				break
			}
		}
		benchmarkResult = result
	}
}

func BenchmarkCollect(b *testing.B) {
	source := benchmarkValues()

	b.ReportAllocs()
	for b.Loop() {
		benchmarkResult = lix.FromSlice(source).Collect()
	}
}

func BenchmarkSortBy(b *testing.B) {
	source := benchmarkValues()

	b.ReportAllocs()
	for b.Loop() {
		benchmarkResult = lix.FromSlice(source).SortBy(func(value int) int {
			return value
		}).Collect()
	}
}

func BenchmarkReverse(b *testing.B) {
	source := benchmarkValues()

	b.ReportAllocs()
	for b.Loop() {
		benchmarkResult = lix.FromSlice(source).Reverse().Collect()
	}
}

func BenchmarkUniqBy(b *testing.B) {
	source := benchmarkValues()
	for _, keys := range benchmarkCardinalities {
		b.Run("keys="+strconv.Itoa(keys), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkResult = lix.FromSlice(source).UniqBy(func(value int) int {
					return value % keys
				}).Collect()
			}
		})
	}
}

// 分组基准只遍历分组结果，不再收集为 map，以隔离分组本身的开销。
func BenchmarkGroupBy(b *testing.B) {
	source := benchmarkValues()
	for _, keys := range benchmarkCardinalities {
		b.Run("keys="+strconv.Itoa(keys), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkCount = lix.GroupBy(lix.FromSlice(source), func(value int) int {
					return value % keys
				}).Length()
			}
		})
	}
}

func BenchmarkGroupPairsBy(b *testing.B) {
	source := lix.FromSlice(benchmarkValues()).Enumerate()
	for _, keys := range benchmarkCardinalities {
		b.Run("keys="+strconv.Itoa(keys), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkCount = lix.GroupPairsBy(source, func(_ int, value int) int {
					return value % keys
				}).Length()
			}
		})
	}
}
