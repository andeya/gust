# gust 🌬️

[![tag](https://img.shields.io/github/tag/andeya/gust.svg)](https://github.com/andeya/gust/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-%23007d9c)
[![GoDoc](https://godoc.org/github.com/andeya/gust?status.svg)](https://pkg.go.dev/github.com/andeya/gust)
![Build Status](https://github.com/andeya/gust/actions/workflows/go-ci.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/andeya/gust)](https://goreportcard.com/report/github.com/andeya/gust)
[![Coverage](https://img.shields.io/codecov/c/github.com/andeya/gust)](https://codecov.io/gh/andeya/gust)
[![License](https://img.shields.io/github/license/andeya/gust)](./LICENSE)

**将 Rust 的优雅带入 Go** - 一个强大的库，让错误处理、可选值和迭代在 Go 中变得像在 Rust 中一样优雅和安全。

> 🎯 **零依赖** • 🚀 **生产就绪** • 📚 **文档完善** • ✨ **类型安全**

**语言:** [English](./README.md) | [中文](./README_ZH.md)

## ✨ 为什么选择 gust？

厌倦了到处写 `if err != nil`？受够了 nil 指针 panic？想要在 Go 中使用 Rust 风格的迭代器链？

**gust** 将 Rust 的最佳模式带入 Go，使您的代码：
- 🛡️ **更安全** - 不再有 nil 指针 panic
- 🎯 **更简洁** - 优雅地链式操作
- 🚀 **更具表现力** - 表达你的意图，而不是样板代码

### 从命令式到声明式

gust 帮助您从**命令式**（关注*如何*）转向**声明式**（关注*什么*）编程：

![声明式 vs 命令式](./doc/declarative_vs_imperative.jpg)

使用 gust，您描述的是**想要实现什么**，而不是**如何一步步实现**。这使得您的代码更易读、更易维护，且更不容易出错。

### 使用 gust 之前（传统 Go）

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
    result := fmt.Sprintf("%s: %s", user.Name, profile.Bio)
    return result, nil
}
```

### 使用 gust 之后（优雅且安全）

```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/ret"

func fetchUserData(userID int) gust.Result[string] {
    return ret.AndThen(gust.Ret(getUser(userID)), func(user *User) gust.Result[string] {
        if user == nil || user.Email == "" {
            return gust.Err[string]("invalid user")
        }
        return ret.Map(gust.Ret(getProfile(user.Email)), func(profile *Profile) string {
            return fmt.Sprintf("%s: %s", user.Name, profile.Bio)
        })
    })
}

// 查看 examples/ 中的 ExampleResult_fetchUserData 获取完整可运行示例
```

**改变了什么？**
- ✅ **没有错误样板代码** - 错误自然地在链中流动
- ✅ **没有嵌套 if-else** - 线性流程，易于阅读
- ✅ **自动传播** - 错误自动停止链的执行
- ✅ **可组合** - 每个步骤都是独立且可测试的
- ✅ **类型安全** - 编译器强制正确的错误处理

## 🚀 快速开始

```bash
go get github.com/andeya/gust
```

## 📚 核心功能

### 1. Result<T> - 优雅的错误处理

用可链式调用的 `Result[T]` 替换 `(T, error)`：

```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/ret"

// 链式操作可能失败的操作
result := gust.Ok(10).
    Map(func(x int) int { return x * 2 }).
    AndThen(func(x int) gust.Result[int] {
        if x > 15 {
            return gust.Err[int]("too large")
        }
        return gust.Ok(x + 5)
    }).
    OrElse(func(err error) gust.Result[int] {
        fmt.Println("Error handled:", err)
        return gust.Ok(0) // 回退值
    })

fmt.Println("Final value:", result.Unwrap())
// Output: Error handled: too large
// Final value: 0
```

**核心优势：**
- ✅ 不再需要 `if err != nil` 样板代码
- ✅ 自动错误传播
- ✅ 优雅地链式多个操作
- ✅ 类型安全的错误处理

### 2. Option<T> - 不再有 Nil Panic

用安全的 `Option[T]` 替换 `*T` 和 `(T, bool)`：

```go
// 安全的除法，无需 nil 检查
divide := func(a, b float64) gust.Option[float64] {
    if b == 0 {
        return gust.None[float64]()
    }
    return gust.Some(a / b)
}

result := divide(10, 2).
    Map(func(x float64) float64 { return x * 2 }).
    UnwrapOr(0)

fmt.Println(result) // 10
```

**核心优势：**
- ✅ 消除 nil 指针 panic
- ✅ 明确的可选值
- ✅ 安全地链式操作
- ✅ 编译器强制安全

### 3. Iterator - Go 中的 Rust 风格迭代

完整的 Rust Iterator trait 实现，支持方法链：

```go
import "github.com/andeya/gust/iter"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

sum := iter.FromSlice(numbers).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Take(3).
    Fold(0, func(acc int, x int) int {
        return acc + x
    })

fmt.Println(sum) // 56 (4 + 16 + 36)
```

**可用方法：**
- **适配器**: `Map`, `Filter`, `Chain`, `Zip`, `Enumerate`, `Skip`, `Take`, `StepBy`, `FlatMap`, `Flatten`
- **消费者**: `Fold`, `Reduce`, `Collect`, `Count`, `All`, `Any`, `Find`, `Sum`, `Product`, `Partition`
- **高级**: `Scan`, `Intersperse`, `Peekable`, `ArrayChunks`, `FindMap`, `MapWhile`
- **双端**: `NextBack`, `Rfold`, `TryRfold`, `Rfind`
- 还有 60+ 个来自 Rust Iterator trait 的方法！

**注意：** 对于类型转换操作（例如，从 `string` 到 `int` 的 `Map`），请使用函数式 API：
```go
iter.Map(iter.FromSlice(strings), func(s string) int { return len(s) })
```

对于相同类型的操作，您可以使用方法链：
```go
iter.FromSlice(numbers).Filter(func(x int) bool { return x > 0 }).Map(func(x int) int { return x * 2 })
```

**核心优势：**
- ✅ Rust 风格的方法链
- ✅ 惰性求值
- ✅ 类型安全的转换
- ✅ 尽可能零拷贝

#### Go 标准迭代器集成

gust 迭代器与 Go 1.23+ 标准迭代器无缝集成：

**将 gust Iterator 转换为 Go 的 `iter.Seq[T]`：**
```go
import "github.com/andeya/gust/iter"

numbers := []int{1, 2, 3, 4, 5}
gustIter := iter.FromSlice(numbers).Filter(func(x int) bool { return x%2 == 0 })

// 在 Go 标准的 for-range 循环中使用
for v := range gustIter.Seq() {
    fmt.Println(v) // 输出 2, 4
}
```

**将 Go 的 `iter.Seq[T]` 转换为 gust Iterator：**
```go
import "github.com/andeya/gust/iter"

// 创建 Go 标准迭代器序列
goSeq := func(yield func(int) bool) {
    for i := 0; i < 5; i++ {
        if !yield(i) {
            return
        }
    }
}

// 转换为 gust Iterator 并使用 gust 方法
gustIter, deferStop := iter.FromSeq(goSeq)
defer deferStop()
result := gustIter.Map(func(x int) int { return x * 2 }).Collect()
fmt.Println(result) // [0 2 4 6 8]
```

### 4. 双端迭代器

从两端迭代：

```go
import "github.com/andeya/gust/iter"

numbers := []int{1, 2, 3, 4, 5}
deIter := iter.FromSlice(numbers).MustToDoubleEnded()

// 从前端迭代
if val := deIter.Next(); val.IsSome() {
    fmt.Println("Front:", val.Unwrap()) // Front: 1
}

// 从后端迭代
if val := deIter.NextBack(); val.IsSome() {
    fmt.Println("Back:", val.Unwrap()) // Back: 5
}
```

## 📖 示例

### 解析和过滤错误处理

```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/iter"
import "strconv"

// 将字符串解析为整数，自动过滤错误
numbers := []string{"1", "2", "three", "4", "five"}

results := iter.FilterMap(
    iter.RetMap(iter.FromSlice(numbers), strconv.Atoi),
    gust.Result[int].Ok,
).
    Collect()

fmt.Println("Parsed numbers:", results)
// Output: Parsed numbers: [1 2 4]
```

### 真实世界的数据管道

```go
// 处理用户输入：解析、验证、转换、限制
input := []string{"10", "20", "invalid", "30", "0", "40"}

results := iter.FilterMap(
    iter.RetMap(iter.FromSlice(input), strconv.Atoi),
    gust.Result[int].Ok,
).
    Filter(func(x int) bool { return x > 0 }).
    Map(func(x int) int { return x * 2 }).
    Take(3).
    Collect()

fmt.Println(results) // [20 40 60]
```

### Option 链式操作

```go
// 在可选值上链式操作并过滤
result := gust.Some(5).
    Map(func(x int) int { return x * 2 }).
    Filter(func(x int) bool { return x > 8 }).
    XMap(func(x int) any {
        return fmt.Sprintf("Value: %d", x)
    }).
    UnwrapOr("No value")

fmt.Println(result) // "Value: 10"
```

### 数据分区

```go
// 将数字分为偶数和奇数
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

evens, odds := iter.FromSlice(numbers).
    Partition(func(x int) bool { return x%2 == 0 })

fmt.Println("Evens:", evens) // [2 4 6 8 10]
fmt.Println("Odds:", odds)   // [1 3 5 7 9]
```

## 📦 附加包

gust 提供了多个工具包来扩展其功能：

- **`gust/dict`** - 通用 map 工具（Filter, Map, Keys, Values 等）
- **`gust/vec`** - 通用 slice 工具
- **`gust/valconv`** - 类型安全的值转换
- **`gust/digit`** - 数字转换工具
- **`gust/opt`** - `Option[T]` 辅助函数（Map, AndThen, Zip, Unzip, Assert 等）
- **`gust/ret`** - `Result[T]` 辅助函数（Map, AndThen, Assert, Flatten 等）
- **`gust/iter`** - Rust 风格迭代器实现（参见上面的[迭代器部分](#3-iterator---go-中的-rust-风格迭代)）

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

**Opt 工具：**
```go
import "github.com/andeya/gust/opt"

some := gust.Some(5)
doubled := opt.Map(some, func(x int) int { return x * 2 })
```

**Ret 工具：**
```go
import "github.com/andeya/gust/ret"

result := gust.Ok(10)
doubled := ret.Map(result, func(x int) int { return x * 2 })
```

更多详细信息，请参阅[完整文档](https://pkg.go.dev/github.com/andeya/gust)和[示例](./examples/)。

### 详细示例

#### Dict 工具

```go
import "github.com/andeya/gust/dict"

m := map[string]int{"a": 1, "b": 2, "c": 3}

// 使用 Option 获取
value := dict.Get(m, "b")
fmt.Println(value.UnwrapOr(0)) // 2

// 过滤 map
filtered := dict.Filter(m, func(k string, v int) bool {
    return v > 1
})
fmt.Println(filtered) // map[b:2 c:3]

// 映射值
mapped := dict.MapValue(m, func(k string, v int) int {
    return v * 2
})
fmt.Println(mapped) // map[a:2 b:4 c:6]
```

#### Vec 工具
```go
import "github.com/andeya/gust/vec"

numbers := []int{1, 2, 3, 4, 5}
doubled := vec.MapAlone(numbers, func(x int) int { return x * 2 })
fmt.Println(doubled) // [2 4 6 8 10]
```

#### Opt 工具
```go
import "github.com/andeya/gust/opt"

some := gust.Some(5)
doubled := opt.Map(some, func(x int) int { return x * 2 })
zipped := opt.Zip(gust.Some(1), gust.Some("hello"))
```

#### Ret 工具
```go
import "github.com/andeya/gust/ret"

result := gust.Ok(10)
doubled := ret.Map(result, func(x int) int { return x * 2 })
chained := ret.AndThen(gust.Ok(5), func(x int) gust.Result[int] {
    return gust.Ok(x * 2)
})
```

## 🔗 资源

- 📖 [完整文档](https://pkg.go.dev/github.com/andeya/gust) - 完整的 API 参考
- 💡 [示例](./examples/) - 按功能组织的综合示例
- 🌐 [English Documentation](./README.md) - English documentation
- 🐛 [问题追踪](https://github.com/andeya/gust/issues) - 报告 bug 或请求功能
- 💬 [讨论](https://github.com/andeya/gust/discussions) - 提问和分享想法

## 📋 要求

需要 **Go 1.23+**（支持泛型和标准迭代器）

## 🤝 贡献

欢迎贡献！无论是：
- 🐛 报告 bug
- 💡 建议新功能
- 📝 改进文档
- 🔧 提交 pull request

每一个贡献都让 gust 变得更好！请随时提交 Pull Request 或打开 issue。

## 📄 许可证

本项目采用 MIT 许可证（MIT License）。

---

**为 Go 社区用心制作 ❤️**

*灵感来自 Rust 的 `Result`、`Option` 和 `Iterator` traits*

