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

### ✨ Catch 模式：gust 的秘密武器

gust 引入了 **`result.Ret + Unwrap + Catch`** 模式——一种革命性的 Go 错误处理方式：

```go
func fetchUserData(userID int) (r result.Result[string]) {
    defer r.Catch()  // 一行代码处理所有错误！
    user := result.Ret(db.GetUser(userID)).Unwrap()
    profile := result.Ret(api.GetProfile(user.Email)).Unwrap()
    return result.Ok(fmt.Sprintf("%s: %s", user.Name, profile.Bio))
}
```

**一行代码** (`defer r.Catch()`) 消除了**所有** `if err != nil` 检查。错误通过 panic 自动传播，被捕获、转换为 `Result` 并返回。

### ✨ 为什么选择 gust？

| 传统 Go | 使用 gust |
|---------|-----------|
| ❌ 15+ 行错误检查代码 | ✅ 3 行 Catch 模式代码 |
| ❌ 到处都是 `if err != nil` | ✅ 只需一次 `defer r.Catch()` |
| ❌ Nil 指针 panic | ✅ 编译时安全 |
| ❌ 命令式循环 | ✅ 声明式管道 |
| ❌ 难以组合 | ✅ 优雅的方法链式调用 |

---

## 🚀 快速开始

```bash
go get github.com/andeya/gust
```

### 您的第一个 gust 程序（使用 Catch 模式）

```go
package main

import (
    "fmt"
    "github.com/andeya/gust/result"
)

func main() {
    // 使用 Catch 模式 - 错误自动流动！
    processValue := func(value int) (r result.Result[int]) {
        defer r.Catch()
        doubled := value * 2
        if doubled > 20 {
            return result.TryErr[int]("too large")
        }
        return result.Ok(doubled + 5)
    }

    res := processValue(10)
    if res.IsOk() {
        fmt.Println("Success:", res.Unwrap())
    } else {
        fmt.Println("Error:", res.UnwrapErr())
    }
}
```

**输出：** `Success: 25`

---

## 💡 gust 解决的问题

### 之前：传统 Go 代码（15+ 行，4 个错误检查）

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
    if profile == nil {
        return "", fmt.Errorf("profile not found")
    }
    return fmt.Sprintf("%s: %s", user.Name, profile.Bio), nil
}
```

**问题：**
- ❌ 4 个重复的 `if err != nil` 检查
- ❌ 3 个嵌套条件判断
- ❌ 难以测试单个步骤
- ❌ 容易忘记错误处理
- ❌ 15+ 行样板代码

### 之后：使用 gust Catch 模式（8 行，0 个错误检查）

```go
import "github.com/andeya/gust/result"

func fetchUserData(userID int) (r result.Result[string]) {
    defer r.Catch()  // 一行代码处理所有错误！
    user := result.Ret(db.GetUser(userID)).Unwrap()
    if user == nil || user.Email == "" {
        return result.TryErr[string]("invalid user")
    }
    profile := result.Ret(api.GetProfile(user.Email)).Unwrap()
    if profile == nil {
        return result.TryErr[string]("profile not found")
    }
    return result.Ok(fmt.Sprintf("%s: %s", user.Name, profile.Bio))
}
```

**优势：**
- ✅ **一行错误处理** - `defer r.Catch()` 处理一切
- ✅ **线性流程** - 易于从上到下阅读
- ✅ **自动传播** - 错误自动停止执行
- ✅ **可组合** - 每个步骤独立且可测试
- ✅ **类型安全** - 编译器强制正确的错误处理
- ✅ **代码减少 70%** - 从 15+ 行减少到 8 行

---

## 📚 核心功能

### 1. Result<T> - Catch 模式革命

**Catch 模式** (`result.Ret + Unwrap + Catch`) 是 gust 最强大的功能：

```go
import "github.com/andeya/gust/result"

// 之前：传统 Go（多个错误检查）
// func readConfig(filename string) (string, error) {
//     f, err := os.Open(filename)
//     if err != nil {
//         return "", err
//     }
//     defer f.Close()
//     data, err := io.ReadAll(f)
//     if err != nil {
//         return "", err
//     }
//     return string(data), nil
// }

// 之后：gust Catch 模式（线性流程，无错误检查）
func readConfig(filename string) (r result.Result[string]) {
    defer r.Catch()  // 一行代码处理所有错误！
    data := result.Ret(os.ReadFile(filename)).Unwrap()
    return result.Ok(string(data))
}
```

**关键方法：**
- `result.Ret(T, error)` - 将 `(T, error)` 转换为 `Result[T]`
- `Unwrap()` - 提取值（如果错误则 panic，被 `Catch` 捕获）
- `defer r.Catch()` - 捕获所有 panic 并转换为 `Result` 错误
- `Map` - 如果 Ok 则转换值
- `AndThen` - 链式调用返回 Result 的操作
- `UnwrapOr` - 安全提取值（带默认值，**永不 panic**）

**实际应用场景：**
- API 调用链
- 数据库操作
- 文件 I/O 操作
- 数据验证管道

### 2. Option<T> - 告别 Nil Panic

用安全的 `Option[T]` 替换 `*T` 和 `(T, bool)`，防止 nil 指针 panic：

```go
import "github.com/andeya/gust/option"

// 之前：传统 Go（到处都是 nil 检查）
// func divide(a, b float64) *float64 {
//     if b == 0 {
//         return nil
//     }
//     result := a / b
//     return &result
// }
// result := divide(10, 2)
// if result != nil {
//     fmt.Println(*result * 2)  // 存在 nil 指针 panic 风险
// }

// 之后：gust Option（类型安全，无 nil panic）
divide := func(a, b float64) option.Option[float64] {
    if b == 0 {
        return option.None[float64]()
    }
    return option.Some(a / b)
}

quotient := divide(10, 2).
    Map(func(x float64) float64 { return x * 2 }).
    UnwrapOr(0)  // 安全：永不 panic

fmt.Println(quotient) // 10
```

**关键方法：**
- `Map` - 如果 Some 则转换值
- `AndThen` - 链式调用返回 Option 的操作
- `Filter` - 条件过滤值
- `UnwrapOr` - 安全提取值（带默认值，**永不 panic**）

**实际应用场景：**
- 配置读取
- 可选函数参数
- Map 查找
- JSON 反序列化

### 3. Iterator - Rust 风格的迭代

完整的 Rust Iterator trait 实现，提供 **60+ 方法**用于声明式数据处理：

```go
import "github.com/andeya/gust/iterator"

// 之前：传统 Go（嵌套循环，手动错误处理）
// func processNumbers(input []string) ([]int, error) {
//     var results []int
//     for _, s := range input {
//         n, err := strconv.Atoi(s)
//         if err != nil {
//             continue
//         }
//         if n > 0 {
//             results = append(results, n*2)
//         }
//     }
//     return results, nil
// }

// 之后：gust Iterator（声明式，类型安全，代码减少 70%）
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

**亮点：**
- 🚀 **60+ 方法**来自 Rust 的 Iterator trait
- 🔄 **惰性求值** - 按需计算
- 🔗 **方法链式调用** - 优雅地组合复杂操作
- 🔌 **Go 1.24+ 集成** - 与标准 `iter.Seq[T]` 协作
- 🎯 **类型安全** - 编译时保证
- ⚡ **零开销抽象** - 无性能开销

**方法分类：**
- **构造函数**: `FromSlice`, `FromRange`, `FromFunc`, `Empty`, `Once`, `Repeat`
- **BitSet 迭代器**: `FromBitSet`, `FromBitSetOnes`, `FromBitSetZeros`
- **Go 集成**: `FromSeq`, `Seq`, `Pull` (Go 1.24+ 标准迭代器)
- **基本适配器**: `Map`, `Filter`, `Chain`, `Zip`, `Enumerate`
- **过滤**: `Skip`, `Take`, `StepBy`, `SkipWhile`, `TakeWhile`
- **转换**: `MapWhile`, `Scan`, `FlatMap`, `Flatten`
- **分块**: `MapWindows`, `ArrayChunks`, `ChunkBy`
- **消费者**: `Collect`, `Fold`, `Reduce`, `Count`, `Sum`, `Product`, `Partition`
- **搜索**: `Find`, `FindMap`, `Position`, `All`, `Any`
- **最值**: `Max`, `Min`, `MaxBy`, `MinBy`, `MaxByKey`, `MinByKey`
- **双端**: `NextBack`, `Rfold`, `Rfind`, `NthBack`

---

## 🌟 实际案例

### 案例 1：数据处理管道（Iterator + Result）

**之前：传统 Go**（嵌套循环 + 错误处理，15+ 行）

```go
func processUserInput(input []string) ([]int, error) {
    var results []int
    for _, s := range input {
        n, err := strconv.Atoi(s)
        if err != nil {
            continue
        }
        if n > 0 {
            results = append(results, n*2)
        }
    }
    if len(results) == 0 {
        return nil, fmt.Errorf("no valid numbers")
    }
    return results, nil
}
```

**之后：gust Iterator + Result**（声明式，类型安全，8 行）

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

**结果：** 代码减少 70%，类型安全，声明式

### 案例 2：API 调用链（Catch 模式）

**之前：传统 Go**（15+ 行，4 个错误检查）

```go
func fetchUserProfile(userID int) (string, error) {
    user, err := db.GetUser(userID)
    if err != nil {
        return "", fmt.Errorf("db error: %w", err)
    }
    if user == nil || user.Email == "" {
        return "", fmt.Errorf("invalid user")
    }
    profile, err := api.GetProfile(user.Email)
    if err != nil {
        return "", fmt.Errorf("api error: %w", err)
    }
    if profile == nil {
        return "", fmt.Errorf("profile not found")
    }
    return fmt.Sprintf("%s: %s", user.Name, profile.Bio), nil
}
```

**之后：gust Catch 模式**（8 行，0 个错误检查）

```go
import "github.com/andeya/gust/result"

func fetchUserProfile(userID int) (r result.Result[string]) {
    defer r.Catch()  // 一行代码处理所有错误！
    user := result.Ret(db.GetUser(userID)).Unwrap()
    if user == nil || user.Email == "" {
        return result.TryErr[string]("invalid user")
    }
    profile := result.Ret(api.GetProfile(user.Email)).Unwrap()
    if profile == nil {
        return result.TryErr[string]("profile not found")
    }
    return result.Ok(fmt.Sprintf("%s: %s", user.Name, profile.Bio))
}

// 使用
profileRes := fetchUserProfile(123)
if profileRes.IsOk() {
    fmt.Println(profileRes.Unwrap())
} else {
    fmt.Println("Error:", profileRes.UnwrapErr())
}
```

**结果：** 代码减少 70%，线性流程，自动错误传播

### 案例 3：文件系统操作（Catch 模式）

**之前：传统 Go**（多个错误检查，嵌套条件）

```go
func copyDirectory(src, dst string) error {
    info, err := os.Stat(src)
    if err != nil {
        return err
    }
    if err = os.MkdirAll(dst, info.Mode()); err != nil {
        return err
    }
    entries, err := os.ReadDir(src)
    if err != nil {
        return err
    }
    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        dstPath := filepath.Join(dst, entry.Name())
        if entry.IsDir() {
            if err = copyDirectory(srcPath, dstPath); err != nil {
                return err
            }
        } else {
            if err = copyFile(srcPath, dstPath); err != nil {
                return err
            }
        }
    }
    return nil
}
```

**之后：gust Catch 模式**（线性流程，单一错误处理器）

```go
import (
    "github.com/andeya/gust/fileutil"
    "github.com/andeya/gust/result"
    "os"
    "path/filepath"
)

func copyDirectory(src, dst string) (r result.VoidResult) {
    defer r.Catch()  // 一行代码处理所有错误！
    info := result.Ret(os.Stat(src)).Unwrap()
    result.RetVoid(os.MkdirAll(dst, info.Mode())).Unwrap()
    entries := result.Ret(os.ReadDir(src)).Unwrap()
    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        dstPath := filepath.Join(dst, entry.Name())
        if entry.IsDir() {
            copyDirectory(srcPath, dstPath).Unwrap()
        } else {
            fileutil.CopyFile(srcPath, dstPath).Unwrap()
        }
    }
    return result.OkVoid()
}
```

**结果：** 线性代码流程，自动错误传播，代码减少 70%

### 案例 4：配置管理（Option）

**之前：传统 Go**（nil 检查，错误处理）

```go
type Config struct {
    APIKey *string
    Port   int
}

func loadConfig() (Config, error) {
    apiKeyEnv := os.Getenv("API_KEY")
    var apiKey *string
    if apiKeyEnv != "" {
        apiKey = &apiKeyEnv
    }
    portStr := os.Getenv("PORT")
    port := 8080
    if portStr != "" {
        p, err := strconv.Atoi(portStr)
        if err != nil {
            return Config{}, err
        }
        port = p
    }
    return Config{APIKey: apiKey, Port: port}, nil
}
```

**之后：gust Option**（类型安全，无 nil 检查）

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
port := config.Port.UnwrapOr(8080)   // 如果未设置，默认为 8080
apiKey := config.APIKey.UnwrapOr("") // 如果未设置，默认为空字符串
```

**结果：** 类型安全，无 nil 检查，优雅的默认值

---

## 📦 完整包生态系统

gust 为常见的 Go 任务提供了一套全面的工具包：

| 包 | 描述 | 关键功能 |
|---------|-------------|--------------|
| **`gust/result`** | 类型安全的错误处理 | `Result[T]`, Catch 模式, `Map`, `AndThen` |
| **`gust/option`** | 安全的可选值 | `Option[T]`, `Map`, `Filter`, `AndThen` |
| **`gust/iterator`** | Rust 风格的迭代 | 60+ 方法，惰性求值，方法链式调用 |
| **`gust/dict`** | 泛型 map 工具 | `Filter`, `Map`, `Keys`, `Values`, `Get` |
| **`gust/vec`** | 泛型 slice 工具 | `MapAlone`, `Get`, `Copy`, `Dict` |
| **`gust/conv`** | 类型安全转换 | `BytesToString`, `StringToReadonlyBytes`, 大小写转换, JSON 引用 |
| **`gust/digit`** | 数字转换 | Base 2-62 转换, `FormatByDict`, `ParseByDict` |
| **`gust/random`** | 安全随机字符串 | Base36/Base62 编码, 时间戳嵌入 |
| **`gust/encrypt`** | 加密哈希函数 | MD5, SHA 系列, FNV, CRC, Adler-32, AES 加密 |
| **`gust/bitset`** | 线程安全位集合 | 位运算, 迭代器集成, 多种编码 |
| **`gust/syncutil`** | 并发工具 | `SyncMap`, `Lazy`, mutex 包装器 |
| **`gust/errutil`** | 错误工具 | 堆栈跟踪, panic 恢复, `ErrBox` |
| **`gust/constraints`** | 类型约束 | `Ordering`, `Numeric`, `Digit` |
| **`gust/fileutil`** | 文件操作 | 路径操作, 文件 I/O, 目录操作, tar.gz 归档 |
| **`gust/coarsetime`** | 快速粗粒度时间 | 实时时间 & 单调时间, 可配置精度, 比 `time.Now()` 快 30 倍 |
| **`gust/shutdown`** | 优雅关闭与重启 | 信号处理, 清理钩子, 优雅进程重启 (Unix) |

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

每一个贡献都让 gust 变得更好！

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

*受 Rust 的 `Result`、`Option` 和 `Iterator` traits 启发*

[⭐ 在 GitHub 上给我们点星](https://github.com/andeya/gust) • [📖 文档](https://pkg.go.dev/github.com/andeya/gust) • [🐛 报告 Bug](https://github.com/andeya/gust/issues) • [💡 请求功能](https://github.com/andeya/gust/issues/new)

</div>
