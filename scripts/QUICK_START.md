# 泛型实例化分析 - 快速开始

## 快速命令参考

### 分析构建代码（go build）
```bash
# 快速分析（当前目录）
./build_with_generics_analysis.sh

# 快速分析（指定目录）
./build_with_generics_analysis.sh -d /path/to/project
./build_with_generics_analysis.sh -d ./subpackage

# 详细分析（当前目录）
python3 analyze_generics.py

# 详细分析（指定目录）
python3 analyze_generics.py -d /path/to/project
python3 analyze_generics.py --dir ./subpackage
```

### 分析测试代码（go test）⭐
```bash
# 快速分析（运行测试，当前目录）
./test_with_generics_analysis.sh

# 快速分析（指定目录）
./test_with_generics_analysis.sh -d /path/to/project
./test_with_generics_analysis.sh -d ./subpackage -v

# 快速分析（只编译不运行）
./test_with_generics_analysis.sh -c
./test_with_generics_analysis.sh -d ./subpackage -c

# 详细分析（当前目录）
python3 analyze_generics.py --test

# 详细分析（指定目录）
python3 analyze_generics.py -d /path/to/project --test
python3 analyze_generics.py --dir ./subpackage -t
```

## 为什么需要分析测试代码？

测试代码可能引入**大量额外的泛型实例化**，因为：

1. **表驱动测试**：使用多种类型组合测试泛型函数
2. **边界测试**：测试各种类型参数的组合
3. **测试辅助函数**：可能使用泛型来减少重复代码

这可能导致 `go test` 编译时内存暴涨，即使 `go build` 没有问题。

## 常见场景

### 场景 1: go build 正常，go test 内存爆炸
```bash
# 1. 先检查构建代码
./build_with_generics_analysis.sh

# 2. 再检查测试代码
./test_with_generics_analysis.sh

# 3. 对比差异，找出测试引入的额外实例化
python3 analyze_generics.py --test
```

### 场景 2: 只想看测试代码的泛型使用
```bash
# 只编译测试，不运行
./test_with_generics_analysis.sh -c

# 或使用详细分析
python3 analyze_generics.py --test --save
```

### 场景 3: 分析指定目录/子包
```bash
# 分析子包
./build_with_generics_analysis.sh -d ./subpackage
./test_with_generics_analysis.sh -d ./subpackage

# 分析其他项目
python3 analyze_generics.py -d /path/to/other/project
python3 analyze_generics.py -d /path/to/other/project --test
```

### 场景 4: 对比优化前后
```bash
# 优化前
./test_with_generics_analysis.sh > before.txt

# 优化代码...

# 优化后
./test_with_generics_analysis.sh > after.txt

# 对比
diff before.txt after.txt
```

## 关键指标

关注这些数字：
- **总实例化次数**：如果 > 1000，可能需要优化
- **高频实例化**：单个函数 > 50 次实例化
- **测试文件实例化**：如果测试文件占比 > 50%，重点优化测试代码
- **多类型参数泛型**：可能导致组合爆炸

## 示例输出解读

```
📁 文件类型统计:
  测试文件中的实例化: 850 次 (70.0%)
  源代码文件中的实例化: 350 次 (30.0%)
```

**解读**：测试代码引入了大量实例化，应该重点优化测试文件。

## 更多信息

查看 `GENERICS_ANALYSIS.md` 获取完整文档。

