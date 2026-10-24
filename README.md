# lix

[![CI](https://github.com/norniastar/lix/actions/workflows/ci.yml/badge.svg)](https://github.com/norniastar/lix/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/norniastar/lix.svg)](https://pkg.go.dev/github.com/norniastar/lix)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Version](https://img.shields.io/github/v/tag/norniastar/lix)](https://github.com/norniastar/lix/releases/tag/v1.0.0)

`lix` 是基于 Go 1.27 泛型方法和标准库 `iter` 构建的类型安全惰性集合库。
它不使用反射，不进行运行时类型分派，错误的回调类型会在编译期被拒绝。

## 要求

Go 1.27 或更高版本。

## 安装

```bash
go get github.com/norniastar/lix@latest
```

## 快速开始

```go
result := lix.FromSlice([]int{1, 2, 3, 4}).
	Filter(func(value int) bool {
		return value%2 == 0
	}).
	Map(func(value int) string {
		return strconv.Itoa(value)
	}).
	Collect()

// result: []string{"2", "4"}
```

泛型方法可以在保持链式调用的同时改变元素类型，例如
`Seq[int].Map(...)` 可以直接得到 `Seq[string]`。

常用操作采用 `lo/it` 的命名习惯，例如 `Drop`、`UniqBy`、`KeyBy`、
`SomeBy`、`EveryBy` 和 `Length`。`lix` 使用链式泛型方法，因此与 `lo/it`
共享操作词汇，但不是仅替换 import 路径即可迁移的兼容实现。

## 类型模型

```go
type Seq[V any] iter.Seq[V]
type Seq2[K comparable, V any] iter.Seq2[K, V]

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}
```

- `Seq[V]` 表示值序列。
- `Seq2[K,V]` 表示键值序列。
- `Pair[K,V]` 用于在值序列与键值序列之间进行强类型转换。
- `Iter()` 可以将 `Seq`、`Seq2` 转回标准库迭代器。

## 功能

下表按数据流阶段划分，每个操作只归入一个类别。

| 类别 | `Seq` | `Seq2` 或顶层函数 |
| --- | --- | --- |
| 构造与互操作 | `FromSlice`、`FromSeq`、`Of`、`RangeWithSteps` | `FromMap`、`FromSeq2`、`FromPairs`、`Iter`、`Pairs` |
| 形状与键转换 | `Map`、`FlatMap`、`FlatMap2`、`Enumerate`、`KeyBy`、`Associate` | `MapValues`、`MapKeys`、`FlatMap`、`FlatMap2`、`Keys`、`Values` |
| 筛选与窗口 | `Filter`、`Reject`、`UniqBy`、`Take`、`Drop` | 同名方法 |
| 排序 | `SortBy`、`SortByDescending`、`Reverse` | 同名方法 |
| 分组与归约 | `Reduce`、`GroupBy` | `Reduce`、`GroupPairsBy` |
| 查询与遍历 | `First`、`Last`、`Find`、`FindIndex`、`SomeBy`、`EveryBy`、`Length`、`ForEach` | 同名方法 |
| 物化 | `Collect() []V` | `Collect() map[K]V` |

## API 示例

每个公开函数和方法都有一个使用 Go 官方 Example 机制编写的简短示例。
这些示例会被 `go test` 执行，并在发布后显示于
[Go Reference](https://pkg.go.dev/github.com/norniastar/lix) 对应的 API 下。

| 范围 | 可执行示例 |
| --- | --- |
| 构造与顶层函数 | [`FromSlice`](https://pkg.go.dev/github.com/norniastar/lix#example-FromSlice)、[`FromMap`](https://pkg.go.dev/github.com/norniastar/lix#example-FromMap)、[`FromSeq`](https://pkg.go.dev/github.com/norniastar/lix#example-FromSeq)、[`FromSeq2`](https://pkg.go.dev/github.com/norniastar/lix#example-FromSeq2)、[`FromPairs`](https://pkg.go.dev/github.com/norniastar/lix#example-FromPairs)、[`Of`](https://pkg.go.dev/github.com/norniastar/lix#example-Of)、[`RangeWithSteps`](https://pkg.go.dev/github.com/norniastar/lix#example-RangeWithSteps)、[`GroupBy`](https://pkg.go.dev/github.com/norniastar/lix#example-GroupBy)、[`GroupPairsBy`](https://pkg.go.dev/github.com/norniastar/lix#example-GroupPairsBy)、[`Pairs`](https://pkg.go.dev/github.com/norniastar/lix#example-Pairs) |
| `Seq` 转换 | [`Iter`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Iter)、[`Map`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Map)、[`Filter`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Filter)、[`Reject`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Reject)、[`FlatMap`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.FlatMap)、[`FlatMap2`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.FlatMap2)、[`Take`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Take)、[`Drop`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Drop)、[`Enumerate`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Enumerate)、[`KeyBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.KeyBy)、[`Associate`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Associate) |
| `Seq` 聚合与排序 | [`Reduce`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Reduce)、[`UniqBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.UniqBy)、[`SortBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.SortBy)、[`SortByDescending`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.SortByDescending)、[`Reverse`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Reverse) |
| `Seq` 查询与物化 | [`Collect`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Collect)、[`First`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.First)、[`Last`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Last)、[`Find`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Find)、[`FindIndex`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.FindIndex)、[`SomeBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.SomeBy)、[`EveryBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.EveryBy)、[`Length`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.Length)、[`ForEach`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq.ForEach) |
| `Seq2` 转换 | [`Iter`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Iter)、[`MapValues`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.MapValues)、[`MapKeys`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.MapKeys)、[`Filter`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Filter)、[`Reject`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Reject)、[`FlatMap`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.FlatMap)、[`FlatMap2`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.FlatMap2)、[`Take`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Take)、[`Drop`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Drop) |
| `Seq2` 聚合与排序 | [`Reduce`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Reduce)、[`UniqBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.UniqBy)、[`SortBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.SortBy)、[`SortByDescending`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.SortByDescending)、[`Reverse`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Reverse) |
| `Seq2` 查询与物化 | [`Collect`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Collect)、[`First`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.First)、[`Last`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Last)、[`Find`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Find)、[`FindIndex`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.FindIndex)、[`SomeBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.SomeBy)、[`EveryBy`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.EveryBy)、[`Length`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Length)、[`ForEach`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.ForEach)、[`Keys`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Keys)、[`Values`](https://pkg.go.dev/github.com/norniastar/lix#example-Seq2.Values) |

示例源码按 API 类别放在仓库的 `example*_test.go` 文件中。本地可以一次运行全部示例：

```bash
go test -run '^Example' ./...
```

### 键和值

`Enumerate` 为值附加从零开始的遍历序号：

```go
result := lix.Of("a", "b").
	Enumerate().
	MapValues(func(index int, value string) string {
		return fmt.Sprintf("%d:%s", index, value)
	}).
	Collect()
```

`KeyBy` 保留值并生成键；`Associate` 可以同时生成新的键和值：

```go
byID := users.KeyBy(func(user User) int {
	return user.ID
})

nameLengths := users.Associate(func(user User) (int, int) {
	return user.ID, len(user.Name)
})
```

### 分组

分组返回 `Seq2`，因此可以继续链式处理，最后再决定是否收集为 map：

```go
groups := lix.GroupBy(numbers, func(value int) string {
	if value%2 == 0 {
		return "even"
	}
	return "odd"
})

result := groups.Collect() // map[string][]int
```

Go 会拒绝导致方法集无限实例化的泛型方法，因此 `GroupBy`、`GroupPairsBy` 和
`Pairs` 使用顶层泛型函数表达。它们返回的仍是 `Seq` 或 `Seq2`，不会中断后续链式处理。

## 执行语义

- `Map`、`MapValues`、`MapKeys`、`Filter`、`Reject`、`FlatMap`、`Enumerate`、`KeyBy`、
  `UniqBy`、`Take` 和 `Drop` 都是惰性的。
- `SortBy`、`SortByDescending`、`Reverse` 和分组操作在返回序列被遍历时物化源数据。
- `First`、`Find`、`FindIndex`、`SomeBy`、`EveryBy` 和 `Take` 会尽早停止上游迭代。
- `Collect`、`Reduce`、`Last`、`Length` 和 `ForEach` 会消费完整序列。
- 稳定排序会保留排序键相同元素的原始顺序。
- 分组按分组键首次出现的顺序产出。
- `Seq2.Collect` 遇到重复键时保留最后一个值。

从切片、map 和强类型 Pair 切片构造的序列可以重复遍历。通过 `FromSeq` 或
`FromSeq2` 转换的序列遵循原始迭代器自己的可重复遍历规则。map 的遍历顺序不确定。

## 边界

- `RangeWithSteps` 的 `step` 为零时返回空序列；结果溢出 `int` 范围前停止。
- `Take` 的 `count` 非正时返回空序列，`Drop` 的 `count` 非正时返回全部元素。
  这与 `lo/it` 遇到负数参数时 panic 的行为不同。
- `Seq.Collect` 对空序列返回非 nil 的空切片，而 `slices.Collect` 返回 nil；
  `Seq2.Collect` 与 `maps.Collect` 一致，返回非 nil 的空 map。
- 与标准库 `iter.Seq` 相同，零值 `Seq`、`Seq2` 是 nil 函数，遍历时会 panic；
  空序列请使用 `lix.Of[T]()`。`FlatMap`、`FlatMap2` 的 mapper 同样不应返回零值序列。
- `Seq2` 要求键可比较，以保证始终可以 `Collect` 为 map；
  因此键类型不可比较的 `iter.Seq2` 无法通过 `FromSeq2` 转换。
- 不提供基于反射的动态数据源转换。
- 不支持通过字符串访问结构体字段；请使用类型安全的 selector。
- 不向任意运行时目标指针写值；请使用 `Collect` 或 `Reduce`。
- 数组应先切片，例如 `lix.FromSlice(array[:])`。
- 所有操作顺序执行；并行执行不在当前范围内。

## 兼容性

当前稳定版本为 `v1.0.0`。项目遵循语义化版本；`v1` 内保持公开 API
向后兼容，破坏性变更只会在新的主版本中发布。

## 质量验证

项目包含停止传播、短路、边界、性质、Fuzz、变异测试和 CRAP 分析。
全部公开 API 的示例也属于测试套件，示例输出变化会直接导致测试失败。
这些本地检查的运行方法与等价变异说明见 [QUALITY.md](QUALITY.md)。

## 许可证

[MIT](LICENSE)
