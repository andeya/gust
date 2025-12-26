#!/bin/bash

# 泛型实例化分析脚本
# 用于找出导致编译器内存暴涨的泛型代码

echo "正在分析泛型实例化信息..."
echo "=========================================="
echo ""

# 获取构建输出
BUILD_OUTPUT=$(go build -gcflags=-m 2>&1)

# 统计泛型实例化
echo "【泛型实例化统计】"
echo "----------------------------------------"
echo ""

# 提取所有包含 go.shape 的行
SHAPE_LINES=$(echo "$BUILD_OUTPUT" | grep -E "go\.shape")

# 统计每个泛型类型的实例化次数
echo "1. 按泛型类型统计实例化次数："
echo "$SHAPE_LINES" | grep -oE '[a-zA-Z0-9_]+\[go\.shape\.[^\]]+\]' | \
    sed 's/\[go\.shape\.[^\]]*\]//' | \
    sort | uniq -c | sort -rn | head -20
echo ""

# 统计不同 shape 类型的数量
echo "2. Shape 类型统计（不同实例化类型数量）："
echo "$SHAPE_LINES" | grep -oE 'go\.shape\.[^\]]+' | \
    sort | uniq -c | sort -rn | head -20
echo ""

# 统计每个文件的泛型实例化
echo "3. 按文件统计泛型实例化："
echo "$SHAPE_LINES" | grep -oE '^[^:]+:[0-9]+:[0-9]+' | \
    cut -d: -f1 | sort | uniq -c | sort -rn | head -20
echo ""

# 找出最常被实例化的函数/方法
echo "4. 最常被实例化的函数/方法（Top 30）："
echo "$SHAPE_LINES" | grep -oE '[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\[go\.shape\.[^\]]+\]' | \
    sed 's/\[go\.shape\.[^\]]*\]//' | \
    sort | uniq -c | sort -rn | head -30
echo ""

# 统计总的泛型实例化数量
TOTAL_INSTANCES=$(echo "$SHAPE_LINES" | wc -l | tr -d ' ')
UNIQUE_TYPES=$(echo "$SHAPE_LINES" | grep -oE '[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\[go\.shape\.[^\]]+\]' | \
    sed 's/\[go\.shape\.[^\]]*\]//' | sort -u | wc -l | tr -d ' ')

echo "【汇总信息】"
echo "----------------------------------------"
echo "总实例化次数: $TOTAL_INSTANCES"
echo "唯一泛型类型数: $UNIQUE_TYPES"
echo ""

# 找出可能导致内存问题的模式
echo "【潜在问题分析】"
echo "----------------------------------------"
echo ""

# 找出有多个类型参数的泛型（可能导致组合爆炸）
echo "5. 多类型参数的泛型函数（可能导致组合爆炸）："
echo "$SHAPE_LINES" | grep -oE '[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\[go\.shape\.[^,]+,[^]]+\]' | \
    sed 's/\[go\.shape\.[^\]]*\]//' | \
    sort | uniq -c | sort -rn | head -20
echo ""

# 统计内联信息中的泛型
echo "6. 泛型函数内联统计："
INLINE_LINES=$(echo "$BUILD_OUTPUT" | grep "can inline")
GENERIC_INLINE=$(echo "$INLINE_LINES" | grep "go\.shape" | wc -l | tr -d ' ')
TOTAL_INLINE=$(echo "$INLINE_LINES" | wc -l | tr -d ' ')
echo "泛型函数内联数: $GENERIC_INLINE / $TOTAL_INLINE"
echo ""

# 详细输出最频繁的实例化（用于进一步分析）
echo "7. 最频繁的泛型实例化详情（Top 10）："
echo "$SHAPE_LINES" | grep -oE '[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\[go\.shape\.[^\]]+\]' | \
    sort | uniq -c | sort -rn | head -10 | while read count instance; do
    echo "  [$count次] $instance"
    echo "$SHAPE_LINES" | grep "$instance" | head -3 | sed 's/^/    /'
    echo ""
done

echo "=========================================="
echo "分析完成！"
echo ""
echo "提示：如果某个泛型类型被实例化次数过多，"
echo "     可能是导致编译器内存暴涨的原因。"

