#!/bin/bash

# 改进的构建命令，自动分析泛型实例化信息
# 用法: 
#   ./build_with_generics_analysis.sh [选项] [build_args...]
#   选项:
#     -d, --dir DIR    指定要分析的目录（默认：脚本所在目录）
#     -h, --help       显示帮助信息
#   示例:
#     ./build_with_generics_analysis.sh
#     ./build_with_generics_analysis.sh -d /path/to/project
#     ./build_with_generics_analysis.sh -d ./subpackage

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_DIR="$SCRIPT_DIR"
BUILD_ARGS=()

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--dir)
            TARGET_DIR="$2"
            shift 2
            ;;
        -h|--help)
            echo "用法: $0 [选项] [build_args...]"
            echo ""
            echo "选项:"
            echo "  -d, --dir DIR    指定要分析的目录（默认：脚本所在目录）"
            echo "  -h, --help       显示帮助信息"
            echo ""
            echo "示例:"
            echo "  $0"
            echo "  $0 -d /path/to/project"
            echo "  $0 -d ./subpackage"
            exit 0
            ;;
        *)
            BUILD_ARGS+=("$1")
            shift
            ;;
    esac
done

# 切换到目标目录
if [ ! -d "$TARGET_DIR" ]; then
    echo "❌ 错误: 目录不存在: $TARGET_DIR"
    exit 1
fi

cd "$TARGET_DIR"
echo "🔨 开始构建并分析泛型实例化..."
echo "📁 工作目录: $(pwd)"
echo ""

# 运行构建并捕获输出
BUILD_OUTPUT=$(go build -gcflags=-m "${BUILD_ARGS[@]}" 2>&1)
BUILD_EXIT=$?

# 显示构建结果
if [ $BUILD_EXIT -eq 0 ]; then
    echo "✅ 构建成功"
else
    echo "❌ 构建失败"
    echo "$BUILD_OUTPUT"
    exit $BUILD_EXIT
fi

echo ""
echo "📊 正在分析泛型实例化信息..."
echo ""

# 提取泛型实例化统计
SHAPE_LINES=$(echo "$BUILD_OUTPUT" | grep -E "go\.shape" || true)

if [ -z "$SHAPE_LINES" ]; then
    echo "⚠️  未找到泛型实例化信息"
    exit 0
fi

# 快速统计
TOTAL=$(echo "$SHAPE_LINES" | wc -l | tr -d ' ')

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "【快速统计】"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "总泛型实例化次数: $TOTAL"
echo ""

# Top 10 最常被实例化的函数
echo "🔝 Top 10 最常被实例化的泛型函数/类型:"
echo "$SHAPE_LINES" | grep -oE '[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\[go\.shape\.[^\]]+\]' | \
    sed 's/\[go\.shape\.[^\]]*\]//' | \
    sort | uniq -c | sort -rn | head -10 | \
    awk '{printf "  %4d次  %s\n", $1, $2}'
echo ""

# Top 10 Shape 类型
echo "🔝 Top 10 最常见的 Shape 类型:"
echo "$SHAPE_LINES" | grep -oE 'go\.shape\.[^\]]+' | \
    sort | uniq -c | sort -rn | head -10 | \
    awk '{printf "  %4d次  %s\n", $1, $2}'
echo ""

# 按文件统计
echo "🔝 Top 10 按文件统计:"
echo "$SHAPE_LINES" | grep -oE '^[^:]+:[0-9]+:[0-9]+' | \
    cut -d: -f1 | sort | uniq -c | sort -rn | head -10 | \
    awk '{printf "  %4d次  %s\n", $1, $2}'
echo ""

# 多类型参数统计
MULTI_PARAM=$(echo "$SHAPE_LINES" | grep -cE 'go\.shape\.[^,]+,' || echo "0")
if [ "$MULTI_PARAM" -gt 0 ]; then
    echo "⚠️  多类型参数的泛型实例化: $MULTI_PARAM 次"
    echo "    (可能导致组合爆炸，需要特别关注)"
    echo ""
fi

# 找出高频实例化（可能的问题源头）
echo "🚨 高频实例化分析（可能的内存问题源头）:"
HIGH_FREQ=$(echo "$SHAPE_LINES" | grep -oE '[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\[go\.shape\.[^\]]+\]' | \
    sed 's/\[go\.shape\.[^\]]*\]//' | \
    sort | uniq -c | sort -rn | awk '$1 > 50 {print $1, $2}')
if [ -n "$HIGH_FREQ" ]; then
    echo "$HIGH_FREQ" | awk '{printf "  ⚠️  %s: %d次实例化\n", $2, $1}'
else
    echo "  ✅ 未发现异常高频的实例化（所有函数实例化次数 <= 50）"
fi
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💡 提示: 运行 'python3 analyze_generics.py -d \"$TARGET_DIR\"' 获取更详细的分析报告"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

