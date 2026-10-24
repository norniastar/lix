# 变更日志

本项目的所有重要变更都会记录在此文件中。

格式参考 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，
版本遵循[语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 变更

- `GroupBy`、`GroupPairsBy` 改用“键到位置”的索引加有序分组切片，每个元素只需一次 map 查找。
  键大量重复时更快（基准：16 个分组快 18%–26%，365 个分组快 5%–13%）；
  键几乎全部唯一时耗时不变，内存增加约 29%。签名与行为不变。
- `SortBy`、`SortByDescending` 的 `Seq` 与 `Seq2` 实现合并为共享的稳定排序帮助函数；
  `Reverse` 改用 `slices.Backward`，`Collect` 改用 `slices.AppendSeq` 与 `maps.Collect`。
  行为与性能不变。

### 文档

- README 边界说明补充负数 `Take`/`Drop`、空序列 `Collect`、零值序列和 `Seq2` 键约束的语义。

### 测试

- 新增全部中间操作的构造期惰性契约测试，以及带状态派生序列的可重复遍历契约测试。
- 新增 `Collect`、`SortBy`、`Reverse`、`UniqBy`、`GroupBy`、`GroupPairsBy` 基准与手写循环基线，
  覆盖 `Seq`、`Seq2` 和升降序变体。
- 新增 `GroupPairsBy` 按首次出现顺序分组的确定性测试。
- 新增覆盖完整 `int` 范围的 `RangeWithSteps` Fuzz，以及 `Seq`/`Seq2` 升降序排序稳定性 Fuzz。

## [1.0.0] - 2026-08-22

### 新增

- 基于标准库 `iter` 的 `Seq`、`Seq2` 和 `Pair` 类型。
- 覆盖映射、展开、筛选、去重、窗口、排序、分组、归约、检索和物化的完整操作集。
- `Enumerate`、`KeyBy` 和 `Associate` 等类型安全键转换能力。
- `Drop`、`UniqBy`、`SomeBy`、`EveryBy` 和 `Length` 等与 `lo/it` 一致的操作名称。
- 切片、map、标准迭代器、Pair 切片、可变参数和整数区间数据源。
- `Seq` 与 `Seq2` 对称的查询、短路和遍历操作。
- 可执行示例、行为测试、基准测试和 CI。
- 为全部公开函数和方法提供由 `go test` 验证的 Go Example。
- 停止传播、短路、空序列、计数边界和重复键契约测试。
- `RangeWithSteps`、窗口、反转、去重、Pair 往返、稳定排序和分组性质的 Fuzz 测试。
- 固定版本的本地变异测试与 CRAP 报告脚本。

[Unreleased]: https://github.com/norniastar/lix/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/norniastar/lix/releases/tag/v1.0.0
