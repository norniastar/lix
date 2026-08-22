# 质量验证

`lix` 使用三层测试证明公共行为：

1. 单操作契约测试验证边界、惰性执行和停止传播。
2. 组合管道测试验证不同操作能够正确连接。
3. 性质测试与 Fuzz 将实现和简单参考模型进行比较。

这些检查作为本地质量工具使用，不属于 CI 门禁。

## 普通测试

```bash
go test -race -coverprofile=coverage.out ./...
```

## Fuzz

运行全部 Fuzz 目标，每个目标默认执行三秒：

```bash
./scripts/fuzz.sh
```

可以通过环境变量调整单个目标的运行时间：

```bash
FUZZ_TIME=30s ./scripts/fuzz.sh
```

## 变异测试

项目固定使用 Gremlins `v0.6.0`，通过 `go run` 执行，不进行全局安装：

```bash
./scripts/mutation.sh
```

机器可读报告写入 `target/quality/mutation.json`。可以调整并发和超时预算：

```bash
MUTATION_WORKERS=2 \
MUTATION_TEST_CPU=0 \
MUTATION_TIMEOUT_COEFFICIENT=20 \
./scripts/mutation.sh
```

`MUTATION_TEST_CPU=0` 会让脚本完全省略 `--test-cpu` 参数。过低的 CPU 限制会
制造大量假超时，因此默认不限制单次测试可用的 CPU。

以下四种边界变异在当前前置条件下与原程序等价，应人工审核而不是通过弱化测试隐藏：

- `rangeContinues` 中 `step > 0` 变为 `step >= 0`：`step == 0` 会在进入该函数前返回。
- `rangeOverflows` 中 `step > 0` 变为 `step >= 0`：`step == 0` 会在进入该函数前返回。
- 正步长分支中 `next <= current` 变为 `next < current`：非零 `int` 步长不可能令 `next == current`。
- 负步长分支中 `next >= current` 变为 `next > current`：同上。

当前基线：

- 可运行变异体：47
- 被测试杀死：43
- 等价变异体：4
- 未覆盖、超时、不可运行：0
- Gremlins 原始 efficacy：91.49%
- 排除已审核等价变异后的有效杀伤率：100%

## CRAP

项目固定使用 `crap4go` 的伪版本
`v0.0.0-20260521190544-bee16dbdadb4`：

```bash
./scripts/crap.sh
```

报告写入 `target/quality/crap.txt`。CRAP 计算公式为：

```text
CRAP = CC² × (1 - coverage)³ + CC
```

当前目标是每个生产函数的 CRAP 小于 6。发现超标时，先判断缺少的行为测试，
再判断函数是否需要拆分，不能只通过无意义断言提高覆盖率。

当前报告覆盖率为 100%，最高 CRAP 为 5.0；最高复杂度函数是 `RangeWithSteps`、
`GroupBy`、`GroupPairsBy`、`Seq.Take` 和 `Seq2.Take`。
