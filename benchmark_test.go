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
	benchmarkMap    map[int]int
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

// benchmarkCardinalities 覆盖键大量重复、中等重复（约按天聚合一年）和几乎全部唯一三种情况，
// 分别对应 map 查找以命中为主和以插入为主的路径。
var benchmarkCardinalities = []int{16, 365, benchmarkSize}

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

// Seq2 基准使用 Enumerate 生成的唯一键；除 Collect 外均以 Length 消费结果，
// 避免把构建 map 的开销计入被测操作。
func benchmarkPairs() lix.Seq2[int, int] {
	return lix.FromSlice(benchmarkValues()).Enumerate()
}

func BenchmarkCollect(b *testing.B) {
	source := benchmarkValues()
	pairs := benchmarkPairs()

	b.Run("Seq", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkResult = lix.FromSlice(source).Collect()
		}
	})
	b.Run("Seq2", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMap = pairs.Collect()
		}
	})
}

func BenchmarkSortBy(b *testing.B) {
	source := benchmarkValues()
	pairs := benchmarkPairs()
	identity := func(value int) int { return value }
	byValue := func(_ int, value int) int { return value }

	b.Run("Seq/ascending", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkResult = lix.FromSlice(source).SortBy(identity).Collect()
		}
	})
	b.Run("Seq/descending", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkResult = lix.FromSlice(source).SortByDescending(identity).Collect()
		}
	})
	b.Run("Seq2/ascending", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkCount = pairs.SortBy(byValue).Length()
		}
	})
	b.Run("Seq2/descending", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkCount = pairs.SortByDescending(byValue).Length()
		}
	})
}

func BenchmarkReverse(b *testing.B) {
	source := benchmarkValues()
	pairs := benchmarkPairs()

	b.Run("Seq", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkResult = lix.FromSlice(source).Reverse().Collect()
		}
	})
	b.Run("Seq2", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkCount = pairs.Reverse().Length()
		}
	})
}

func BenchmarkUniqBy(b *testing.B) {
	source := benchmarkValues()
	pairs := benchmarkPairs()
	for _, keys := range benchmarkCardinalities {
		b.Run("Seq/keys="+strconv.Itoa(keys), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkResult = lix.FromSlice(source).UniqBy(func(value int) int {
					return value % keys
				}).Collect()
			}
		})
		b.Run("Seq2/keys="+strconv.Itoa(keys), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkCount = pairs.UniqBy(func(_ int, value int) int {
					return value % keys
				}).Length()
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
	source := benchmarkPairs()
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
