<div align="center">

# gust 🌬️

**将 Rust 的优雅带入 Go**

*一个生产就绪的库，让错误处理、可选值和迭代在 Go 中变得像在 Rust 中一样优雅和安全。*

[![GitHub release](https://img.shields.io/github/release/andeya/gust.svg)](https://github.com/andeya/gust/releases)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-00ADD8?style=flat&logo=go)](https://golang.org)
[![GoDoc](https://pkg.go.dev/badge/github.com/andeya/gust.svg)](https://pkg.go.dev/github.com/andeya/gust)
[![CI Status](https://github.com/andeya/gust/actions/workflows/go-ci.yml/badge.svg)](https://github.com/andeya/gust/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/andeya/gust)](https://goreportcard.com/report/github.com/andeya/gust)
[![Code Coverage](https://codecov.io/gh/andeya/gust/branch/main/graph/badge.svg)](https://codecov.io/gh/andeya/gust)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

**[English](./README.md)** | **[中文](./README_ZH.md)**

</div>

---

## 🎯 什么是 gust？

**gust** 是一个全面的 Go 库，将 Rust 最强大的模式带入 Go，使您能够编写更安全、更简洁、更具表现力的代码。**零依赖**且**生产就绪**，gust 改变了您在 Go 中处理错误、可选值和数据迭代的方式。

### ✨ 核心特性

- 🛡️ **类型安全的错误处理** - 用可链式调用的 `Result[T]` 替换 `(T, error)`
- 🎯 **不再有 Nil Panic** - 使用 `Option[T]` 替代 `*T` 或 `(T, bool)`
- 🚀 **Rust 风格迭代器** - 完整的 Iterator trait 实现，包含 60+ 方法
- ⚡ **零依赖** - 纯 Go，无外部依赖
- 📚 **文档完善** - 包含真实世界示例的全面文档
- 🔒 **生产就绪** - 高测试覆盖率和经过实战检验

---

## 🚀 快速开始

```bash
go get github.com/andeya/gust
```

### 您的第一个 gust 程序

```go
package main

import (
    "fmt"
    "github.com/andeya/gust/result"
)

func main() {
    // 优雅地链式操作
    res := result.Ok(10).
        Map(func(x int) int { return x * 2 }).
        AndThen(func(x int) result.Result[int] {
            if x > 20 {
                return result.TryErr[int]("too large")
            }
            return result.Ok(x + 5)
        })

    if res.IsOk() {
        fmt.Println("Success:", res.Unwrap())
    } else {
        fmt.Println("Error:", res.UnwrapErr())
    }
}
```

---

## 💡 为什么选择 gust？

### 传统 Go 的问题

传统 Go 代码冗长且容易出错：

```go
func fetchUserData(userID int) (string, error) {
    // 步骤 1: 从数据库获取
    user, err := db.GetUser(userID)
    if err != nil {
        return "", fmt.Errorf("db error: %w", err)
    }
    
    // 步骤 2: 验证用户
    if user == nil {
        return "", fmt.Errorf("user not found")
    }
    if user.Email == "" {
        return "", fmt.Errorf("invalid user: no email")
    }
    
    // 步骤 3: 获取配置文件
    profile, err := api.GetProfile(user.Email)
    if err != nil {
        return "", fmt.Errorf("api error: %w", err)
    }
    
    // 步骤 4: 格式化结果
    return fmt.Sprintf("%s: %s", user.Name, profile.Bio), nil
}
```

**问题：**
- ❌ 重复的错误处理样板代码
- ❌ 嵌套的 if-else 语句
- ❌ 难以组合和测试
- ❌ 容易忘记错误检查

### gust 解决方案

使用 gust，编写声明式、可组合的代码：

```go
import "github.com/andeya/gust/result"

func fetchUserData(userID int) result.Result[string] {
    return result.AndThen(result.Ret(getUser(userID)), func(user *User) result.Result[string] {
        if user == nil || user.Email == "" {
            return result.TryErr[string]("invalid user")
        }
        return result.Map(result.Ret(getProfile(user.Email)), func(profile *Profile) string {
            return fmt.Sprintf("%s: %s", user.Name, profile.Bio)
        })
    })
}
```

**优势：**
- ✅ **没有错误样板代码** - 错误自然地在链中流动
- ✅ **线性流程** - 易于阅读和理解
- ✅ **自动传播** - 错误自动停止链的执行
- ✅ **可组合** - 每个步骤都是独立且可测试的
- ✅ **类型安全** - 编译器强制正确的错误处理

### 从命令式到声明式

gust 帮助您从**命令式**（关注*如何*）转向**声明式**（关注*什么*）编程：

![声明式 vs 命令式](./doc/declarative_vs_imperative.jpg)

使用 gust，您描述的是**想要实现什么**，而不是**如何一步步实现**。这使得您的代码更易读、更易维护，且更不容易出错。

---

## 📚 核心功能

### 1. Result<T> - 优雅的错误处理

用可链式调用的 `Result[T]` 替换 `(T, error)`，实现类型安全的错误处理：

```go
import "github.com/andeya/gust/result"

// 链式操作可能失败的操作
res := result.Ok(10).
    Map(func(x int) int { return x * 2 }).
    AndThen(func(x int) result.Result[int] {
        if x > 15 {
            return result.TryErr[int]("too large")
        }
        return result.Ok(x + 5)
    }).
    OrElse(func(err error) result.Result[int] {
        fmt.Println("Error handled:", err)
        return result.Ok(0) // 回退值
    })

fmt.Println("Final value:", res.Unwrap())
// Output: Error handled: too large
// Final value: 0
```

**核心方法：**
- `Map` - 如果 Ok 则转换值
- `AndThen` - 链式返回 Result 的操作
- `OrElse` - 使用回退值处理错误
- `Unwrap` / `UnwrapOr` - 安全地提取值
- `IsOk` / `IsErr` - 检查结果状态

**优势：**
- ✅ 不再需要 `if err != nil` 样板代码
- ✅ 自动错误传播
- ✅ 优雅地链式多个操作
- ✅ 类型安全的错误处理

### 2. Option<T> - 不再有 Nil Panic

用安全的 `Option[T]` 替换 `*T` 和 `(T, bool)`：

```go
import "github.com/andeya/gust/option"

// 安全的除法，无需 nil 检查
divide := func(a, b float64) option.Option[float64] {
    if b == 0 {
        return option.None[float64]()
    }
    return option.Some(a / b)
}

res := divide(10, 2).
    Map(func(x float64) float64 { return x * 2 }).
    UnwrapOr(0)

fmt.Println(res) // 10
```

**核心方法：**
- `Map` - 如果 Some 则转换值
- `AndThen` - 链式返回 Option 的操作
- `Filter` - 条件过滤值
- `Unwrap` / `UnwrapOr` - 安全地提取值
- `IsSome` / `IsNone` - 检查选项状态

**优势：**
- ✅ 消除 nil 指针 panic
- ✅ 明确的可选值
- ✅ 安全地链式操作
- ✅ 编译器强制安全

### 3. Iterator - Go 中的 Rust 风格迭代

完整的 Rust Iterator trait 实现，支持方法链和惰性求值：

```go
import "github.com/andeya/gust/iterator"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

sum := iterator.FromSlice(numbers).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Take(3).
    Fold(0, func(acc int, x int) int {
        return acc + x
    })

fmt.Println(sum) // 56 (4 + 16 + 36)
```

**可用方法：**

| 类别 | 方法 |
|------|------|
| **构造函数** | `FromSlice`, `FromElements`, `FromRange`, `FromFunc`, `FromIterable`, `Empty`, `Once`, `Repeat` |
| **位集合迭代器** | `FromBitSet`, `FromBitSetOnes`, `FromBitSetZeros`, `FromBitSetBytes`, 等 |
| **Go 集成** | `FromSeq`, `FromSeq2`, `FromPull`, `FromPull2`, `Seq`, `Seq2`, `Pull`, `Pull2` |
| **基础适配器** | `Map`, `FilterMap`, `RetMap`, `OptMap`, `Chain`, `Zip`, `Enumerate` |
| **过滤适配器** | `Filter`, `Skip`, `Take`, `StepBy`, `SkipWhile`, `TakeWhile` |
| **转换适配器** | `MapWhile`, `Scan`, `FlatMap`, `Flatten` |
| **分块适配器** | `MapWindows`, `ArrayChunks`, `ChunkBy` |
| **工具适配器** | `Fuse`, `Inspect`, `Intersperse`, `IntersperseWith`, `Cycle`, `Peekable` |
| **消费者** | `Fold`, `Reduce`, `Collect`, `Count`, `Last`, `All`, `Any`, `Find`, `Sum`, `Product`, `Partition`, `AdvanceBy`, `Nth`, `NextChunk` |
| **查找与搜索** | `Find`, `FindMap`, `Position`, `All`, `Any` |
| **最值** | `Max`, `Min`, `MaxBy`, `MinBy`, `MaxByKey`, `MinByKey` |
| **Try 方法** | `TryFold`, `TryForEach`, `TryReduce`, `TryFind` |
| **双端** | `NextBack`, `Rfold`, `TryRfold`, `Rfind`, `AdvanceBackBy`, `NthBack` |

**60+ 个方法**来自 Rust Iterator trait！

**代码组织：**

iterator 包按功能模块化组织，便于维护：

- **核心** (`core.go`): 核心接口 (`Iterable`, `Iterator`, `DoubleEndedIterator`) 和基础类型，包括双端迭代器方法 (`NextBack`, `AdvanceBackBy`, `NthBack`, `Remaining`)
- **构造函数** (`constructors.go`): 从各种数据源创建迭代器的函数
- **基础适配器** (`basic.go`): Map, FilterMap, Chain, Zip, Enumerate, FlatMap
- **过滤适配器** (`filtering.go`): Skip, Take, StepBy, SkipWhile, TakeWhile
- **转换适配器** (`transforming.go`): MapWhile, Scan, Flatten
- **分块适配器** (`chunking.go`): MapWindows, ArrayChunks, ChunkBy
- **工具适配器** (`utility.go`): Fuse, Inspect, Intersperse, IntersperseWith, Cycle, Peekable, Cloned
- **消费者** (`consumers.go`): Collect, Count, Last, Partition, AdvanceBy, Nth, NextChunk, Sum, Product, Unzip, TryReduce, TryForEach
- **折叠与归约** (`fold_reduce.go`): Fold, Reduce, ForEach, TryFold, Rfold, TryRfold
- **查找与搜索** (`find_search.go`): Find, FindMap, Position, All, Any, TryFind, Rfind
- **最值** (`min_max.go`): Max, Min, MaxBy, MinBy, MaxByKey, MinByKey
- **比较** (`comparison.go`): 比较工具

每个模块都是自包含的，包含自己的实现函数 (`_Impl`) 和可迭代结构体 (`_Iterable`)，确保独立性和可维护性。双端迭代器方法已集成到相应的功能模块中（例如，`Rfold` 在 `fold_reduce.go` 中，`Rfind` 在 `find_search.go` 中）。

**注意：** 对于类型转换操作（例如，从 `string` 到 `int` 的 `Map`），请使用函数式 API：

```go
iterator.Map(iterator.FromSlice(strings), func(s string) int { return len(s) })
```

对于相同类型的操作，您可以使用方法链：

```go
iterator.FromSlice(numbers).
    Filter(func(x int) bool { return x > 0 }).
    Map(func(x int) int { return x * 2 })
```

**优势：**
- ✅ Rust 风格的方法链
- ✅ 惰性求值
- ✅ 类型安全的转换
- ✅ 尽可能零拷贝

#### 迭代器构造函数

从各种数据源创建迭代器：

```go
import (
    "github.com/andeya/gust/iterator"
    "github.com/andeya/gust/option"
)

// 从切片创建
iter1 := iterator.FromSlice([]int{1, 2, 3})

// 从单个元素创建
iter2 := iterator.FromElements(1, 2, 3)

// 从范围创建 [start, end)
iter3 := iterator.FromRange(0, 5) // 0, 1, 2, 3, 4

// 从函数创建
count := 0
iter4 := iterator.FromFunc(func() option.Option[int] {
    if count < 3 {
        count++
        return option.Some(count)
    }
    return option.None[int]()
})

// 空迭代器
iter5 := iterator.Empty[int]()

// 单值迭代器
iter6 := iterator.Once(42)

// 无限重复
iter7 := iterator.Repeat("hello") // "hello", "hello", "hello", ...
```

#### Go 标准迭代器集成

gust 迭代器与 Go 1.24+ 标准迭代器无缝集成：

**将 gust Iterator 转换为 Go 的 `iter.Seq[T]`：**

```go
import "github.com/andeya/gust/iterator"

numbers := []int{1, 2, 3, 4, 5}
gustIter := iterator.FromSlice(numbers).Filter(func(x int) bool { return x%2 == 0 })

// 在 Go 标准的 for-range 循环中使用
for v := range gustIter.Seq() {
    fmt.Println(v) // 输出 2, 4
}
```

**将 Go 的 `iter.Seq[T]` 转换为 gust Iterator：**

```go
import "github.com/andeya/gust/iterator"

// 创建 Go 标准迭代器序列
goSeq := func(yield func(int) bool) {
    for i := 0; i < 5; i++ {
        if !yield(i) {
            return
        }
    }
}

// 转换为 gust Iterator 并使用 gust 方法
gustIter, deferStop := iterator.FromSeq(goSeq)
defer deferStop()
result := gustIter.Map(func(x int) int { return x * 2 }).Collect()
fmt.Println(result) // [0 2 4 6 8]
```

### 4. 双端迭代器

从两端高效迭代：

```go
import "github.com/andeya/gust/iterator"

numbers := []int{1, 2, 3, 4, 5}
deIter := iterator.FromSlice(numbers).MustToDoubleEnded()

// 从前端迭代
if val := deIter.Next(); val.IsSome() {
    fmt.Println("Front:", val.Unwrap()) // Front: 1
}

// 从后端迭代
if val := deIter.NextBack(); val.IsSome() {
    fmt.Println("Back:", val.Unwrap()) // Back: 5
}
```

---

## 📖 真实世界示例

### 解析和过滤错误处理

```go
import (
    "github.com/andeya/gust/iterator"
    "github.com/andeya/gust/result"
    "strconv"
)

// 将字符串解析为整数，自动过滤错误
numbers := []string{"1", "2", "three", "4", "five"}

results := iterator.FilterMap(
    iterator.RetMap(iterator.FromSlice(numbers), strconv.Atoi),
    result.Result[int].Ok,
).Collect()

fmt.Println("Parsed numbers:", results)
// Output: Parsed numbers: [1 2 4]
```

### 数据处理管道

```go
import (
    "github.com/andeya/gust/iterator"
    "github.com/andeya/gust/result"
    "strconv"
)

// 处理用户输入：解析、验证、转换、限制
input := []string{"10", "20", "invalid", "30", "0", "40"}

results := iterator.FilterMap(
    iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
    result.Result[int].Ok,
).
    Filter(func(x int) bool { return x > 0 }).
    Map(func(x int) int { return x * 2 }).
    Take(3).
    Collect()

fmt.Println(results) // [20 40 60]
```

### Option 链式操作

```go
import (
    "fmt"
    "github.com/andeya/gust/option"
)

// 在可选值上链式操作并过滤
res := option.Some(5).
    Map(func(x int) int { return x * 2 }).
    Filter(func(x int) bool { return x > 8 }).
    XMap(func(x int) any {
        return fmt.Sprintf("Value: %d", x)
    }).
    UnwrapOr("No value")

fmt.Println(res) // "Value: 10"
```

### 数据分区

```go
import (
    "fmt"
    "github.com/andeya/gust/iterator"
)

// 将数字分为偶数和奇数
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

evens, odds := iterator.FromSlice(numbers).
    Partition(func(x int) bool { return x%2 == 0 })

fmt.Println("Evens:", evens) // [2 4 6 8 10]
fmt.Println("Odds:", odds)   // [1 3 5 7 9]
```

### 位集合迭代

使用完整的迭代器支持迭代位集合或字节切片中的位：

```go
import (
    "fmt"
    "github.com/andeya/gust/iterator"
)

// 迭代字节切片中的位
bytes := []byte{0b10101010, 0b11001100}

// 获取所有设置为 1 的位的偏移量
setBits := iterator.FromBitSetBytesOnes(bytes).
    Filter(func(offset int) bool { return offset > 5 }).
    Collect()
fmt.Println(setBits) // [6 8 9 12 13]

// 统计设置为 1 的位的数量
count := iterator.FromBitSetBytesOnes(bytes).Count()
fmt.Println(count) // 8

// 设置为 1 的位的偏移量之和
sum := iterator.FromBitSetBytesOnes(bytes).
    Fold(0, func(acc, offset int) int { return acc + offset })
fmt.Println(sum) // 54 (0+2+4+6+8+9+12+13)
```

---

## 📦 附加包

gust 提供了多个工具包来扩展其功能：

| 包 | 描述 |
|---------|-------------|
| **`gust/dict`** | 通用 map 工具（Filter, Map, Keys, Values, Get 等） |
| **`gust/vec`** | 通用 slice 工具（MapAlone, Get, Copy, Dict 等） |
| **`gust/conv`** | 类型安全的值转换和反射工具 |
| **`gust/digit`** | 数字转换工具（进制转换, FormatByDict, ParseByDict） |
| **`gust/opt`** | `Option[T]` 辅助函数（Map, AndThen, Zip, Unzip, Assert） |
| **`gust/result`** | `Result[T]` 辅助函数（Map, AndThen, Assert, Flatten） |
| **`gust/iterator`** | Rust 风格迭代器实现（参见上面的[迭代器部分](#3-iterator---go-中的-rust-风格迭代)） |
| **`gust/syncutil`** | 并发工具（SyncMap, Mutex 包装器, 懒加载初始化） |
| **`gust/errutil`** | 错误工具（堆栈跟踪, Panic 恢复, ErrBox） |
| **`gust/constraints`** | 类型约束（Ordering, Numeric 等） |

### 快速示例

**Dict 工具：**
```go
import "github.com/andeya/gust/dict"

m := map[string]int{"a": 1, "b": 2, "c": 3}
value := dict.Get(m, "b").UnwrapOr(0) // 2
filtered := dict.Filter(m, func(k string, v int) bool { return v > 1 })
```

**Vec 工具：**
```go
import "github.com/andeya/gust/vec"

numbers := []int{1, 2, 3, 4, 5}
doubled := vec.MapAlone(numbers, func(x int) int { return x * 2 })
```

**SyncUtil 工具：**
```go
import "github.com/andeya/gust/syncutil"

// 线程安全的 map
var m syncutil.SyncMap[string, int]
m.Store("key", 42)
value := m.Load("key") // 返回 Option[int]

// 懒加载初始化
lazy := syncutil.NewLazyValueWithFunc(func() result.Result[int] {
    return result.Ok(expensiveComputation())
})
value := lazy.TryGetValue() // 只计算一次
```

更多详细信息，请参阅[完整文档](https://pkg.go.dev/github.com/andeya/gust)和[示例](./examples/)。

---

## 🔗 资源

- 📖 **[完整文档](https://pkg.go.dev/github.com/andeya/gust)** - 包含示例的完整 API 参考
- 💡 **[示例](./examples/)** - 按功能组织的综合示例
- 🌐 **[English Documentation](./README.md)** - 英文文档
- 🐛 **[问题追踪](https://github.com/andeya/gust/issues)** - 报告 bug 或请求功能
- 💬 **[讨论](https://github.com/andeya/gust/discussions)** - 提问和分享想法

---

## 📋 要求

- **Go 1.24+**（需要支持泛型和标准迭代器）

---

## 🤝 贡献

我们欢迎贡献！无论您是：

- 🐛 **报告 bug** - 通过报告问题帮助我们改进
- 💡 **建议功能** - 分享您对新功能的想法
- 📝 **改进文档** - 帮助改进我们的文档
- 🔧 **提交 PR** - 贡献代码改进

每一个贡献都让 gust 变得更好！请查看我们的[贡献指南](CONTRIBUTING.md)（如果有）或随时提交 Pull Request 或打开 issue。

### 开发设置

```bash
# 克隆仓库
git clone https://github.com/andeya/gust.git
cd gust

# 运行测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📄 许可证

本项目采用 **MIT 许可证** - 详情请参阅 [LICENSE](./LICENSE) 文件。

---

<div align="center">

**为 Go 社区用心制作 ❤️**

*灵感来自 Rust 的 `Result`、`Option` 和 `Iterator` traits*

[⭐ 在 GitHub 上给我们 Star](https://github.com/andeya/gust) • [📖 文档](https://pkg.go.dev/github.com/andeya/gust) • [🐛 报告 Bug](https://github.com/andeya/gust/issues) • [💡 请求功能](https://github.com/andeya/gust/issues/new)

</div>
