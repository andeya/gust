<div align="center">

# gust 🌬️

**编写像 Rust 一样安全、像函数式编程一样优雅、像原生 Go 一样快速的代码。**

*一个零依赖的库，将 Rust 最强大的模式带入 Go，消除错误处理样板代码、nil 指针 panic 和命令式循环。*

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

**gust** 是一个生产就绪的 Go 库，将 Rust 最强大的模式带入 Go。它通过提供以下功能来改变您编写 Go 代码的方式：

- **类型安全的错误处理** - 使用 `Result[T]` 消除 `if err != nil` 样板代码
- **安全的可选值** - 使用 `Option[T]` 告别 nil 指针 panic
- **声明式迭代** - 60+ 迭代器方法，像 Rust 一样编写数据处理管道

**零依赖**且**完全类型安全**，gust 让您编写更安全、更简洁、更具表现力的 Go 代码——同时不牺牲性能。

### ✨ 为什么选择 gust？

| 传统 Go | 使用 gust |
|---------|-----------|
| ❌ 15+ 行错误检查代码 | ✅ 3 行可链式调用的代码 |
| ❌ 到处都是 `if err != nil` | ✅ 错误自动流动 |
| ❌ Nil 指针 panic | ✅ 编译时安全 |
| ❌ 命令式循环 | ✅ 声明式管道 |
| ❌ 难以组合 | ✅ 优雅的方法链式调用 |

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
    // 优雅地链式操作 - 错误自动流动！
    res := result.Ok(10).
        Map(func(x int) int { return x * 2 }).
        AndThen(func(x int) result.Result[int] {
            if x > 20 {
                return result.TryErr[int]("too large")
            }
            return result.Ok(x + 5)
        })

    fmt.Println(res.UnwrapOr(0)) // 25 (安全：如果错误则返回 0)
}
```

**输出：** `25`

---

## 💡 gust 解决的问题

### 之前：传统 Go 代码

```go
func fetchUserData(userID int) (string, error) {
    user, err := db.GetUser(userID)
    if err != nil {
        return "", fmt.Errorf("db error: %w", err)
    }
    if user == nil {
        return "", fmt.Errorf("user not found")
    }
    if user.Email == "" {
        return "", fmt.Errorf("invalid user: no email")
    }
    profile, err := api.GetProfile(user.Email)
    if err != nil {
        return "", fmt.Errorf("api error: %w", err)
    }
    return fmt.Sprintf("%s: %s", user.Name, profile.Bio), nil
}
```

**问题：**
- ❌ 4 个重复的 `if err != nil` 检查
- ❌ 3 个嵌套的条件语句
- ❌ 难以测试单个步骤
- ❌ 容易忘记错误处理
- ❌ 15 行样板代码

### 之后：使用 gust

```go
import "github.com/andeya/gust/result"

func fetchUserData(userID int) result.Result[string] {
    return result.Ret(db.GetUser(userID)).
        AndThen(func(user *User) result.Result[string] {
            if user == nil || user.Email == "" {
                return result.TryErr[string]("invalid user")
            }
            return result.Ret(api.GetProfile(user.Email)).
                Map(func(profile *Profile) string {
                    return fmt.Sprintf("%s: %s", user.Name, profile.Bio)
                })
        })
}
```

**优势：**
- ✅ **代码减少 70%** - 错误自然流动
- ✅ **线性流程** - 易于从上到下阅读
- ✅ **自动传播** - 错误自动停止链式调用
- ✅ **可组合** - 每个步骤独立且可测试
- ✅ **类型安全** - 编译器强制正确的错误处理

---

## 📚 核心功能

### 1. Result<T> - 类型安全的错误处理

用可链式调用的 `Result[T]` 替换 `(T, error)`，消除错误处理样板代码：

```go
import "github.com/andeya/gust/result"

res := result.Ok(10).
    Map(func(x int) int { return x * 2 }).
    AndThen(func(x int) result.Result[int] {
        if x > 15 {
            return result.TryErr[int]("too large")
        }
        return result.Ok(x + 5)
    }).
    OrElse(func(err error) result.Result[int] {
        return result.Ok(0) // 回退值
    })

fmt.Println(res.UnwrapOr(0)) // 25 (安全，永不 panic)
```

**关键方法：**
- `Map` - 如果 Ok 则转换值
- `AndThen` - 链式调用返回 Result 的操作
- `OrElse` - 使用回退值处理错误
- `UnwrapOr` - 安全提取值（带默认值，**永不 panic**）
- `Unwrap` - 提取值（⚠️ **如果错误会 panic** - 仅在 `IsOk()` 检查后使用）

**实际应用场景：**
- API 调用链
- 数据库操作
- 文件 I/O 操作
- 数据验证管道

### 2. Option<T> - 告别 Nil Panic

用安全的 `Option[T]` 替换 `*T` 和 `(T, bool)`，防止 nil 指针 panic：

```go
import "github.com/andeya/gust/option"

divide := func(a, b float64) option.Option[float64] {
    if b == 0 {
        return option.None[float64]()
    }
    return option.Some(a / b)
}

quotient := divide(10, 2).
    Map(func(x float64) float64 { return x * 2 }).
    Filter(func(x float64) bool { return x > 5 }).
    UnwrapOr(0)

fmt.Println(quotient) // 10
```

**关键方法：**
- `Map` - 如果 Some 则转换值
- `AndThen` - 链式调用返回 Option 的操作
- `Filter` - 条件过滤值
- `UnwrapOr` - 安全提取值（带默认值，**永不 panic**）
- `Unwrap` - 提取值（⚠️ **如果为 None 会 panic** - 仅在 `IsSome()` 检查后使用）

**实际应用场景：**
- 配置读取
- 可选函数参数
- Map 查找
- JSON 反序列化

### 3. Iterator - Rust 风格迭代

完整的 Rust Iterator trait 实现，包含 **60+ 方法**，用于声明式数据处理：

```go
import "github.com/andeya/gust/iterator"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

sum := iterator.FromSlice(numbers).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Take(3).
    Fold(0, func(acc, x int) int { return acc + x })

fmt.Println(sum) // 56 (4 + 16 + 36)
```

**亮点：**
- 🚀 **60+ 方法**来自 Rust Iterator trait
- 🔄 **惰性求值** - 按需计算
- 🔗 **方法链式调用** - 优雅组合复杂操作
- 🔌 **Go 1.24+ 集成** - 与标准 `iter.Seq[T]` 协同工作
- 🎯 **类型安全** - 编译时保证
- ⚡ **零开销抽象** - 无运行时性能损失

**方法分类：**
- **构造函数**: `FromSlice`, `FromRange`, `FromFunc`, `Empty`, `Once`, `Repeat`
- **BitSet 迭代器**: `FromBitSet`, `FromBitSetOnes`, `FromBitSetZeros`
- **Go 集成**: `FromSeq`, `Seq`, `Pull` (Go 1.24+ 标准迭代器)
- **基础适配器**: `Map`, `Filter`, `Chain`, `Zip`, `Enumerate`
- **过滤**: `Skip`, `Take`, `StepBy`, `SkipWhile`, `TakeWhile`
- **转换**: `MapWhile`, `Scan`, `FlatMap`, `Flatten`
- **分块**: `MapWindows`, `ArrayChunks`, `ChunkBy`
- **消费者**: `Collect`, `Fold`, `Reduce`, `Count`, `Sum`, `Product`, `Partition`
- **搜索**: `Find`, `FindMap`, `Position`, `All`, `Any`
- **最值**: `Max`, `Min`, `MaxBy`, `MinBy`, `MaxByKey`, `MinByKey`
- **双端**: `NextBack`, `Rfold`, `Rfind`, `NthBack`

---

## 🌟 实际案例

### 案例 1：数据处理管道

在单个链中解析、验证、转换并限制用户输入：

```go
import (
    "github.com/andeya/gust/iterator"
    "github.com/andeya/gust/result"
    "strconv"
)

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

### 案例 2：带错误处理的 API 调用链

处理多个 API 调用，自动传播错误：

```go
import "github.com/andeya/gust/result"

func fetchUserProfile(userID int) result.Result[string] {
    return result.Ret(db.GetUser(userID)).
        AndThen(func(user *User) result.Result[string] {
            if user == nil || user.Email == "" {
                return result.TryErr[string]("invalid user")
            }
            return result.Ret(api.GetProfile(user.Email)).
                Map(func(profile *Profile) string {
                    return fmt.Sprintf("%s: %s", user.Name, profile.Bio)
                })
        })
}

// 使用
profileRes := fetchUserProfile(123)
if profileRes.IsOk() {
    fmt.Println(profileRes.Unwrap())
} else {
    fmt.Println("Error:", profileRes.UnwrapErr())
}
```

### 案例 3：配置管理

使用 Option 安全地读取和验证配置：

```go
import (
    "github.com/andeya/gust/option"
    "os"
    "strconv"
)

type Config struct {
    APIKey option.Option[string]
    Port   option.Option[int]
}

func loadConfig() Config {
    apiKeyEnv := os.Getenv("API_KEY")
    var apiKeyPtr *string
    if apiKeyEnv != "" {
        apiKeyPtr = &apiKeyEnv
    }
    return Config{
        APIKey: option.ElemOpt(apiKeyPtr),
        Port:   option.RetOpt(strconv.Atoi(os.Getenv("PORT"))),
    }
}

config := loadConfig()
port := config.Port.UnwrapOr(8080) // 如果未设置，默认为 8080
apiKey := config.APIKey.UnwrapOr("") // 如果未设置，默认为空字符串
```

### 案例 4：BitSet 与迭代器

使用迭代器方法处理位集合：

```go
import (
    "github.com/andeya/gust/bitset"
    "github.com/andeya/gust/iterator"
)

bs := bitset.New()
bs.Set(0, true).Unwrap()
bs.Set(5, true).Unwrap()

// 使用迭代器获取所有设置的位
setBits := iterator.FromBitSetOnes(bs).Collect() // [0 5]

// 位运算
bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap()
bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap()
or := bs1.Or(bs2)

// 编码/解码（默认使用 Base64URL）
encoded := bs.String()
decoded := bitset.NewFromBase64URL(encoded).Unwrap()
```

---

## 📦 完整包生态系统

gust 提供了一套全面的工具包，用于常见的 Go 任务：

| 包 | 描述 | 关键特性 |
|---------|-------------|--------------|
| **`gust/result`** | 类型安全的错误处理 | `Result[T]`, `Map`, `AndThen`, `OrElse` |
| **`gust/option`** | 安全的可选值 | `Option[T]`, `Map`, `Filter`, `AndThen` |
| **`gust/iterator`** | Rust 风格迭代 | 60+ 方法，惰性求值，方法链式调用 |
| **`gust/dict`** | 通用 map 工具 | `Filter`, `Map`, `Keys`, `Values`, `Get` |
| **`gust/vec`** | 通用 slice 工具 | `MapAlone`, `Get`, `Copy`, `Dict` |
| **`gust/conv`** | 类型安全转换 | `BytesToString`, `StringToReadonlyBytes`, 大小写转换, JSON 引用 |
| **`gust/digit`** | 数字转换 | Base 2-62 转换, `FormatByDict`, `ParseByDict` |
| **`gust/random`** | 安全随机字符串 | Base36/Base62 编码, 时间戳嵌入 |
| **`gust/encrypt`** | 加密哈希函数 | MD5, SHA 系列, FNV, CRC, Adler-32, AES 加密 |
| **`gust/bitset`** | 线程安全位集合 | 位运算, 迭代器集成, 多种编码 |
| **`gust/syncutil`** | 并发工具 | `SyncMap`, `Lazy`, mutex 包装器 |
| **`gust/errutil`** | 错误工具 | 堆栈跟踪, panic 恢复, `ErrBox` |
| **`gust/constraints`** | 类型约束 | `Ordering`, `Numeric`, `Digit` |

---

## 🎯 为什么选择 gust？

### 零依赖
gust 具有**零外部依赖**。它只使用 Go 的标准库，保持您的项目精简和安全。

### 生产就绪
- ✅ 全面的测试覆盖
- ✅ 完整的文档和示例
- ✅ 在生产环境中经过验证
- ✅ 积极的维护和支持

### 类型安全
所有操作都是**类型安全**的，具有编译时保证。Go 编译器强制正确使用。

### 性能
gust 使用**零开销抽象**。与传统 Go 代码相比，没有运行时开销。

### Go 1.24+ 集成
与 Go 1.24+ 的标准 `iter.Seq[T]` 迭代器无缝协作，弥合 gust 和标准 Go 之间的差距。

### 社区
- 📖 完整的 API 文档
- 💡 每个功能的丰富示例
- 🐛 活跃的问题追踪
- 💬 社区讨论

---

## 🔗 资源

- 📖 **[完整文档](https://pkg.go.dev/github.com/andeya/gust)** - 完整的 API 参考
- 💡 **[示例](./examples/)** - 按功能组织的综合示例
- 🌐 **[English Documentation](./README.md)** - 英文文档
- 🐛 **[问题追踪](https://github.com/andeya/gust/issues)** - 报告 bug 或请求功能
- 💬 **[讨论](https://github.com/andeya/gust/discussions)** - 提问和分享想法

---

## 📋 要求

- **Go 1.24+**（需要泛型和标准迭代器支持）

---

## 🤝 贡献

我们欢迎贡献！无论您是：

- 🐛 **报告 bug** - 帮助我们改进
- 💡 **建议功能** - 分享您的想法
- 📝 **改进文档** - 让文档更好
- 🔧 **提交 PR** - 贡献代码改进

每个贡献都让 gust 变得更好！

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

本项目采用 **MIT 许可证** - 详见 [LICENSE](./LICENSE) 文件。

---

<div align="center">

**为 Go 社区用心打造 ❤️**

*灵感来自 Rust 的 `Result`、`Option` 和 `Iterator` traits*

[⭐ 在 GitHub 上 Star 我们](https://github.com/andeya/gust) • [📖 文档](https://pkg.go.dev/github.com/andeya/gust) • [🐛 报告 Bug](https://github.com/andeya/gust/issues) • [💡 请求功能](https://github.com/andeya/gust/issues/new)

</div>
