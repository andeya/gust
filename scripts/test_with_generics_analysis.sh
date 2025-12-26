#!/bin/bash

# 改进的测试命令，自动分析泛型实例化信息
# 用法: 
#   ./test_with_generics_analysis.sh [选项] [test_args...]
#   选项:
#     -d, --dir DIR    指定要分析的目录（默认：脚本所在目录）
#     -h, --help       显示帮助信息
#   示例:
#     ./test_with_generics_analysis.sh
#     ./test_with_generics_analysis.sh -d /path/to/project
#     ./test_with_generics_analysis.sh -d ./subpackage -v
#     ./test_with_generics_analysis.sh ./... -run TestSomething

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_DIR="$SCRIPT_DIR"
TEST_ARGS=()

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--dir)
            TARGET_DIR="$2"
            shift 2
            ;;
        -h|--help)
            echo "用法: $0 [选项] [test_args...]"
            echo ""
            echo "选项:"
            echo "  -d, --dir DIR    指定要分析的目录（默认：脚本所在目录）"
            echo "  -h, --help       显示帮助信息"
            echo ""
            echo "示例:"
            echo "  $0"
            echo "  $0 -d /path/to/project"
            echo "  $0 -d ./subpackage -v"
            echo "  $0 ./... -run TestSomething"
            echo ""
            echo "注意: -d 参数必须在其他参数之前指定"
            exit 0
            ;;
        *)
            TEST_ARGS+=("$1")
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
echo "🧪 开始测试并分析泛型实例化..."
echo "📁 工作目录: $(pwd)"
echo ""

# 运行测试编译并捕获输出
# 使用 -c 标志只编译不运行，这样可以捕获所有编译输出
echo "正在编译测试代码..."
TEST_OUTPUT=$(go test -c -gcflags=-m "${TEST_ARGS[@]}" 2>&1 || true)

# 也尝试运行测试（如果用户想要运行测试）
RUN_TESTS=true
if [[ " ${TEST_ARGS[@]} " =~ " -c " ]] || [[ " ${TEST_ARGS[@]} " =~ " --compile-only " ]]; then
    RUN_TESTS=false
fi

if [ "$RUN_TESTS" = true ]; then
    echo ""
    echo "正在运行测试..."
    TEST_RESULT=$(go test -gcflags=-m "${TEST_ARGS[@]}" 2>&1)
    TEST_EXIT=$?
    
    # 合并输出
    COMBINED_OUTPUT="$TEST_OUTPUT"$'\n'"$TEST_RESULT"
else
    COMBINED_OUTPUT="$TEST_OUTPUT"
    TEST_EXIT=0
fi

# 提取泛型实例化统计
SHAPE_LINES=$(echo "$COMBINED_OUTPUT" | grep -E "go\.shape" || true)

if [ -z "$SHAPE_LINES" ]; then
    echo ""
    echo "⚠️  未找到泛型实例化信息"
    if [ "$RUN_TESTS" = true ] && [ $TEST_EXIT -ne 0 ]; then
        echo ""
        echo "测试输出："
        echo "$TEST_RESULT"
        exit $TEST_EXIT
    fi
    exit 0
fi

echo ""
echo "📊 正在分析泛型实例化信息..."
echo ""

# 快速统计
TOTAL=$(echo "$SHAPE_LINES" | wc -l | tr -d ' ')

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "【测试代码泛型实例化统计】"
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

# 区分测试文件和源代码文件
echo "📁 测试文件 vs 源代码文件统计:"
TEST_FILES=$(echo "$SHAPE_LINES" | grep -oE '^[^:]+:[0-9]+:[0-9]+' | \
    cut -d: -f1 | grep -E '_test\.go$' | wc -l | tr -d ' ')
SOURCE_FILES=$(echo "$SHAPE_LINES" | grep -oE '^[^:]+:[0-9]+:[0-9]+' | \
    cut -d: -f1 | grep -vE '_test\.go$' | wc -l | tr -d ' ')
echo "  测试文件中的实例化: $TEST_FILES 次"
echo "  源代码文件中的实例化: $SOURCE_FILES 次"
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

# 显示测试结果（如果运行了测试）
if [ "$RUN_TESTS" = true ]; then
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "【测试结果】"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    # 提取测试结果部分（去除编译输出）
    echo "$TEST_RESULT" | grep -v "go\.shape" | tail -20
    echo ""
    
    if [ $TEST_EXIT -eq 0 ]; then
        echo "✅ 测试通过"
    else
        echo "❌ 测试失败"
    fi
    echo ""
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💡 提示:"
echo "  - 运行 'python3 analyze_generics.py -d \"$TARGET_DIR\" --test' 获取更详细的分析报告"
echo "  - 使用 'go test -c -gcflags=-m' 可以只编译不运行测试"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 返回测试退出码
if [ "$RUN_TESTS" = true ]; then
    exit $TEST_EXIT
fi

