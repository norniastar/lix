// Package lix 基于标准库 iter 和 Go 1.27 泛型方法，提供类型安全的惰性迭代管道。
//
// 转换操作默认惰性执行；排序、反转和分组在遍历时物化源数据；
// Collect、Reduce 等终结操作会消费序列。
package lix
