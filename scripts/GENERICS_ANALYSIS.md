# 泛型实例化分析工具

用于分析和诊断 Go 泛型代码导致的编译器内存暴涨问题。

## 工具说明

### 1. `build_with_generics_analysis.sh` - 构建快速分析脚本
快速构建并显示泛型实例化的关键统计信息。

**使用方法：**
```bash
# 分析当前目录
./build_with_generics_analysis.sh

# 分析指定目录
./build_with_generics_analysis.sh -d /path/to/project
./build_with_generics_analysis.sh --dir ./subpackage

# 查看帮助
./build_with_generics_analysis.sh -h
```

**输出内容：**
- 总实例化次数
- Top 10 最常被实例化的泛型函数
- Top 10 最常见的 Shape 类型
- 按文件统计
- 高频实例化警告

### 2. `test_with_generics_analysis.sh` - 测试快速分析脚本 ⭐ 新增
快速测试并显示泛型实例化的关键统计信息（包括测试代码）。

**使用方法：**
```bash
# 运行测试并分析（当前目录）
./test_with_generics_analysis.sh

# 分析指定目录
./test_with_generics_analysis.sh -d /path/to/project
./test_with_generics_analysis.sh -d ./subpackage -v

# 带测试参数
./test_with_generics_analysis.sh -v
./test_with_generics_analysis.sh ./... -run TestSomething

# 只编译不运行测试
./test_with_generics_analysis.sh -c
./test_with_generics_analysis.sh -d ./subpackage -c

# 查看帮助
./test_with_generics_analysis.sh -h
```

**输出内容：**
- 总实例化次数（包括测试代码）
- Top 10 最常被实例化的泛型函数
- Top 10 最常见的 Shape 类型
- 按文件统计（区分测试文件和源代码）
- 测试文件 vs 源代码文件统计
- 高频实例化警告
- 测试结果

### 3. `analyze_generics.py` - 详细分析工具
提供更详细的统计报告和问题诊断。

**使用方法：**
```bash
# 分析构建代码（当前目录）
python3 analyze_generics.py

# 分析构建代码（指定目录）
python3 analyze_generics.py -d /path/to/project
python3 analyze_generics.py --dir ./subpackage

# 分析测试代码（当前目录）
python3 analyze_generics.py --test
# 或
python3 analyze_generics.py -t

# 分析测试代码（指定目录）
python3 analyze_generics.py -d /path/to/project --test
python3 analyze_generics.py --dir ./subpackage -t

# 保存详细报告到文件
python3 analyze_generics.py --save
python3 analyze_generics.py --test --save  # 测试模式 + 保存
python3 analyze_generics.py -d ./subpackage --test --save

# 查看帮助
python3 analyze_generics.py -h
```

**输出内容：**
- 完整的统计报告
- 百分比分析
- 多类型参数泛型检测
- 测试文件 vs 源代码文件统计（测试模式）
- 潜在问题识别
- 优化建议

### 4. `analyze_generics.sh` - Bash 版本详细分析
Bash 实现的详细分析脚本。

**使用方法：**
```bash
./analyze_generics.sh
```

## 如何找出内存问题源头

### 对于构建代码（go build）

#### 步骤 1: 运行快速分析
```bash
./build_with_generics_analysis.sh
```

#### 步骤 2: 运行详细分析
```bash
python3 analyze_generics.py
```

### 对于测试代码（go test）⭐

#### 步骤 1: 运行测试快速分析
```bash
./test_with_generics_analysis.sh
```

**特别关注：**
- **测试文件中的实例化**：测试代码可能引入额外的泛型实例化
- **测试文件 vs 源代码文件统计**：帮助识别问题来源

#### 步骤 2: 运行测试详细分析
```bash
python3 analyze_generics.py --test
```

**测试模式特有功能：**
- 区分测试文件和源代码文件的实例化
- 分析测试代码引入的额外泛型实例化
- 识别测试中的泛型使用模式

### 通用分析步骤

关注：
- **高频实例化警告**：实例化次数 > 50 的函数
- **多类型参数泛型**：可能导致组合爆炸
- **测试文件中的实例化**（测试模式）：测试代码可能引入大量实例化

重点关注：
1. **Top 20 最常被实例化的泛型函数**
   - 如果某个函数实例化次数异常高（>100），可能是问题源头
   
2. **多类型参数的泛型**
   - 例如 `Func[T1, T2, T3]` 如果有多个类型参数，实例化数量会呈指数增长
   
3. **复杂的 Shape 类型**
   - 包含多个类型参数或复杂类型约束的泛型

4. **测试文件中的泛型使用**（测试模式）
   - 测试代码可能使用多种类型组合测试泛型函数
   - 这可能导致测试编译时内存暴涨

### 步骤 3: 分析具体代码

找到问题函数后，检查：
1. **是否真的需要泛型？**
   - 考虑使用接口或类型断言
   
2. **类型参数是否可以合并？**
   - 减少类型参数数量
   
3. **是否可以使用类型约束？**
   - 使用 `interface` 约束来减少实例化

4. **测试代码中的泛型使用**（测试模式）
   - 检查测试是否使用了过多的类型组合
   - 考虑使用表驱动测试减少泛型实例化

## 示例输出解读

```
🔝 Top 10 最常被实例化的泛型函数/类型:
  150次  iter.Iterator
  120次  iter.Map
   80次  iter.Filter
   60次  gust.Option
```

**解读：**
- `iter.Iterator` 被实例化了 150 次，如果这个数字异常高，可能需要优化
- 检查是否有不必要的类型组合导致实例化过多

## 优化建议

1. **减少类型参数数量**
   ```go
   // ❌ 不好：3个类型参数可能导致组合爆炸
   func Process[T1, T2, T3 any](t1 T1, t2 T2, t3 T3) {}
   
   // ✅ 更好：减少类型参数
   func Process[T any](items []T) {}
   ```

2. **使用类型约束**
   ```go
   // ✅ 使用约束减少实例化
   type Numeric interface {
       int | int64 | float64
   }
   func Sum[T Numeric](items []T) T {}
   ```

3. **避免不必要的泛型**
   ```go
   // ❌ 如果只有一种类型，不需要泛型
   func Process[T int](x T) {}
   
   // ✅ 直接使用具体类型
   func Process(x int) {}
   ```

## 内存问题诊断流程

### 构建代码诊断流程
```
1. 运行 ./build_with_generics_analysis.sh
   ↓
2. 查看高频实例化警告
   ↓
3. 运行 python3 analyze_generics.py 获取详细信息
   ↓
4. 定位到具体的泛型函数
   ↓
5. 分析代码，应用优化建议
   ↓
6. 重新运行分析，验证改进效果
```

### 测试代码诊断流程 ⭐
```
1. 运行 ./test_with_generics_analysis.sh
   ↓
2. 查看测试文件 vs 源代码文件统计
   ↓
3. 查看高频实例化警告（特别关注测试文件）
   ↓
4. 运行 python3 analyze_generics.py --test 获取详细信息
   ↓
5. 定位到具体的泛型函数（区分测试代码和源代码）
   ↓
6. 分析测试代码，优化测试中的泛型使用
   ↓
7. 重新运行分析，验证改进效果
```

### 完整诊断流程（推荐）
```
1. 先分析构建代码: ./build_with_generics_analysis.sh
   ↓
2. 再分析测试代码: ./test_with_generics_analysis.sh
   ↓
3. 对比两者差异，找出测试引入的额外实例化
   ↓
4. 使用详细分析工具深入调查
   ↓
5. 优化问题代码
   ↓
6. 重新验证
```

## 注意事项

- 这些工具基于 `go build -gcflags=-m` 的输出
- `go.shape.*` 是 Go 编译器内部使用的形状类型
- 实例化次数高不一定就是问题，需要结合实际情况判断
- 建议在优化前后对比分析结果

