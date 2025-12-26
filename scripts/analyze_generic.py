#!/usr/bin/env python3
"""
泛型实例化分析工具
用于找出导致编译器内存暴涨的泛型代码
"""

import re
import subprocess
import sys
import argparse
import os
from collections import defaultdict, Counter
from typing import Dict, List, Tuple

def run_build(test_mode: bool = False, target_dir: str = ".") -> str:
    """运行构建或测试并获取输出"""
    # 确保目录存在
    if not os.path.isdir(target_dir):
        print(f"❌ 错误: 目录不存在: {target_dir}")
        sys.exit(1)
    
    abs_dir = os.path.abspath(target_dir)
    if test_mode:
        print(f"正在运行 go test -c -gcflags=-m ...")
        print(f"📁 工作目录: {abs_dir}")
        # 先编译测试代码
        compile_result = subprocess.run(
            ["go", "test", "-c", "-gcflags=-m"],
            capture_output=True,
            text=True,
            cwd=abs_dir
        )
        compile_output = compile_result.stdout + compile_result.stderr
        
        # 也运行测试以获取更多信息
        print("正在运行 go test -gcflags=-m ...")
        test_result = subprocess.run(
            ["go", "test", "-gcflags=-m"],
            capture_output=True,
            text=True,
            cwd=abs_dir
        )
        test_output = test_result.stdout + test_result.stderr
        
        return compile_output + "\n" + test_output
    else:
        print(f"正在运行 go build -gcflags=-m ...")
        print(f"📁 工作目录: {abs_dir}")
        result = subprocess.run(
            ["go", "build", "-gcflags=-m"],
            capture_output=True,
            text=True,
            cwd=abs_dir
        )
        return result.stdout + result.stderr

def extract_generic_instances(output: str) -> List[Dict]:
    """提取泛型实例化信息"""
    instances = []
    
    # 匹配模式：函数名[go.shape.类型]
    pattern = r'([a-zA-Z0-9_]+(?:\.[a-zA-Z0-9_]+)*)\[go\.shape\.([^\]]+)\](?::(\d+):(\d+))?'
    
    for line in output.split('\n'):
        matches = re.finditer(pattern, line)
        for match in matches:
            func_name = match.group(1)
            shape_type = match.group(2)
            line_num = match.group(3)
            col_num = match.group(4)
            
            # 提取文件路径
            file_match = re.search(r'([^:]+):(\d+):(\d+)', line)
            file_path = file_match.group(1) if file_match else None
            
            instances.append({
                'func': func_name,
                'shape': shape_type,
                'file': file_path,
                'line': line_num,
                'col': col_num,
                'raw_line': line.strip()
            })
    
    return instances

def analyze_instances(instances: List[Dict]) -> Dict:
    """分析实例化数据"""
    stats = {
        'by_func': Counter(),
        'by_shape': Counter(),
        'by_file': Counter(),
        'by_func_shape': Counter(),
        'multi_param': [],
        'inline_info': [],
        'test_files': Counter(),
        'source_files': Counter()
    }
    
    for inst in instances:
        func = inst['func']
        shape = inst['shape']
        file = inst['file']
        
        stats['by_func'][func] += 1
        stats['by_shape'][shape] += 1
        if file:
            stats['by_file'][file] += 1
            # 区分测试文件和源代码文件
            if file.endswith('_test.go'):
                stats['test_files'][file] += 1
            else:
                stats['source_files'][file] += 1
        stats['by_func_shape'][f"{func}[{shape}]"] += 1
        
        # 检查多类型参数
        if ',' in shape:
            stats['multi_param'].append(inst)
    
    return stats

def print_statistics(stats: Dict, instances: List[Dict]):
    """打印统计信息"""
    print("\n" + "="*80)
    print("【泛型实例化统计报告】")
    print("="*80)
    
    total = len(instances)
    unique_funcs = len(stats['by_func'])
    unique_shapes = len(stats['by_shape'])
    
    print(f"\n📊 总体统计:")
    print(f"  总实例化次数: {total}")
    print(f"  唯一泛型函数/类型数: {unique_funcs}")
    print(f"  唯一 Shape 类型数: {unique_shapes}")
    
    print(f"\n🔝 Top 20 最常被实例化的泛型函数/类型:")
    print("-" * 80)
    for func, count in stats['by_func'].most_common(20):
        percentage = (count / total * 100) if total > 0 else 0
        print(f"  {count:4d}次 ({percentage:5.1f}%)  {func}")
    
    print(f"\n🔝 Top 20 最常见的 Shape 类型:")
    print("-" * 80)
    for shape, count in stats['by_shape'].most_common(20):
        percentage = (count / total * 100) if total > 0 else 0
        print(f"  {count:4d}次 ({percentage:5.1f}%)  {shape}")
    
    print(f"\n🔝 Top 20 按文件统计:")
    print("-" * 80)
    for file, count in stats['by_file'].most_common(20):
        percentage = (count / total * 100) if total > 0 else 0
        file_type = "🧪测试" if file.endswith('_test.go') else "📄源码"
        print(f"  {count:4d}次 ({percentage:5.1f}%)  [{file_type}] {file}")
    
    # 显示测试文件 vs 源代码文件统计
    test_total = sum(stats['test_files'].values())
    source_total = sum(stats['source_files'].values())
    if test_total > 0 or source_total > 0:
        print(f"\n📁 文件类型统计:")
        print("-" * 80)
        print(f"  测试文件中的实例化: {test_total} 次 ({test_total/total*100 if total > 0 else 0:.1f}%)")
        print(f"  源代码文件中的实例化: {source_total} 次 ({source_total/total*100 if total > 0 else 0:.1f}%)")
    
    print(f"\n⚠️  多类型参数的泛型（可能导致组合爆炸）:")
    print("-" * 80)
    multi_param_funcs = Counter()
    for inst in stats['multi_param']:
        multi_param_funcs[inst['func']] += 1
    
    if multi_param_funcs:
        for func, count in multi_param_funcs.most_common(20):
            print(f"  {count:4d}次  {func}")
    else:
        print("  未发现多类型参数的泛型")
    
    print(f"\n🔍 最频繁的泛型实例化组合（Top 15）:")
    print("-" * 80)
    for combo, count in stats['by_func_shape'].most_common(15):
        percentage = (count / total * 100) if total > 0 else 0
        print(f"  {count:4d}次 ({percentage:5.1f}%)  {combo}")
    
    # 找出可能导致内存问题的模式
    print(f"\n🚨 潜在问题分析:")
    print("-" * 80)
    
    # 找出实例化次数异常高的函数
    high_instances = [(f, c) for f, c in stats['by_func'].items() if c > 50]
    if high_instances:
        print("\n  实例化次数 > 50 的函数（可能是内存问题源头）:")
        for func, count in sorted(high_instances, key=lambda x: x[1], reverse=True):
            print(f"    ⚠️  {func}: {count}次")
    
    # 分析复杂 shape 类型
    complex_shapes = [(s, c) for s, c in stats['by_shape'].items() if ',' in s or len(s) > 50]
    if complex_shapes:
        print("\n  复杂的 Shape 类型（可能导致编译变慢）:")
        for shape, count in sorted(complex_shapes, key=lambda x: x[1], reverse=True)[:10]:
            print(f"    ⚠️  {shape[:80]}: {count}次")
    
    print("\n" + "="*80)
    print("💡 建议:")
    print("  1. 关注实例化次数 > 50 的泛型函数")
    print("  2. 检查多类型参数的泛型是否有不必要的组合")
    print("  3. 考虑使用类型约束来减少实例化数量")
    print("  4. 对于频繁实例化的泛型，考虑使用接口或代码生成")
    print("="*80)

def main():
    parser = argparse.ArgumentParser(
        description='泛型实例化分析工具 - 用于找出导致编译器内存暴涨的泛型代码',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  %(prog)s                    # 分析当前目录的构建代码
  %(prog)s --test              # 分析当前目录的测试代码
  %(prog)s -d /path/to/project # 分析指定目录的构建代码
  %(prog)s -d ./subpackage --test  # 分析指定目录的测试代码
  %(prog)s --test --save       # 分析测试代码并保存报告
        """
    )
    parser.add_argument(
        '-d', '--dir',
        default='.',
        help='指定要分析的目录（默认：当前目录）'
    )
    parser.add_argument(
        '-t', '--test',
        action='store_true',
        help='分析测试代码（包括测试文件）'
    )
    parser.add_argument(
        '--save',
        action='store_true',
        help='保存详细报告到文件'
    )
    
    args = parser.parse_args()
    
    test_mode = args.test
    target_dir = args.dir
    
    if test_mode:
        print("泛型实例化分析工具 - 测试模式")
    else:
        print("泛型实例化分析工具 - 构建模式")
    print("="*80)
    
    # 运行构建或测试
    output = run_build(test_mode=test_mode, target_dir=target_dir)
    
    # 提取实例化信息
    print("\n正在提取泛型实例化信息...")
    instances = extract_generic_instances(output)
    
    if not instances:
        print("⚠️  未找到泛型实例化信息")
        mode_str = "测试" if test_mode else "构建"
        print(f"\n提示：确保代码中使用了泛型，并且使用 -gcflags=-m 标志进行{mode_str}")
        sys.exit(1)
    
    # 分析数据
    print("正在分析数据...")
    stats = analyze_instances(instances)
    
    # 打印统计
    print_statistics(stats, instances)
    
    # 可选：保存详细报告
    if args.save:
        report_file = 'generics_test_report.txt' if test_mode else 'generics_report.txt'
        # 如果指定了目录，将报告保存到该目录
        if target_dir != '.':
            report_file = os.path.join(target_dir, report_file)
        with open(report_file, 'w', encoding='utf-8') as f:
            mode_str = "测试" if test_mode else "构建"
            f.write(f"泛型实例化详细报告 ({mode_str}模式)\n")
            f.write(f"工作目录: {os.path.abspath(target_dir)}\n")
            f.write("="*80 + "\n\n")
            for inst in instances:
                f.write(f"{inst['raw_line']}\n")
        print(f"\n📄 详细报告已保存到 {os.path.abspath(report_file)}")

if __name__ == '__main__':
    main()

