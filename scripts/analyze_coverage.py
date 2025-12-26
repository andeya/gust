#!/usr/bin/env python3
"""
分析覆盖率报告，找出低于95%的文件和未覆盖的行
"""

import re
from collections import defaultdict
from typing import Dict, List, Tuple

def parse_cover_out(file_path: str) -> Dict[str, Dict]:
    """解析 cover.out 文件"""
    file_stats = defaultdict(lambda: {
        'total_statements': 0,
        'covered_statements': 0,
        'uncovered_lines': [],
        'low_coverage_lines': []  # 执行次数很少的行
    })
    
    with open(file_path, 'r') as f:
        for line in f:
            line = line.strip()
            if not line or line == 'mode: atomic':
                continue
            
            # 格式: file.go:startLine.startCol,endLine.endCol statements count
            match = re.match(r'([^:]+):(\d+)\.(\d+),(\d+)\.(\d+)\s+(\d+)\s+(\d+)', line)
            if not match:
                continue
            
            file_path = match.group(1)
            start_line = int(match.group(2))
            end_line = int(match.group(4))
            statements = int(match.group(6))
            count = int(match.group(7))
            
            stats = file_stats[file_path]
            stats['total_statements'] += statements
            
            if count > 0:
                stats['covered_statements'] += statements
            else:
                # 记录未覆盖的行范围
                for line_num in range(start_line, end_line + 1):
                    if line_num not in stats['uncovered_lines']:
                        stats['uncovered_lines'].append(line_num)
            
            # 记录执行次数很少的行（小于等于1次）
            if count <= 1 and count > 0:
                for line_num in range(start_line, end_line + 1):
                    if line_num not in stats['low_coverage_lines']:
                        stats['low_coverage_lines'].append(line_num)
    
    return file_stats

def calculate_coverage(stats: Dict) -> float:
    """计算覆盖率百分比"""
    if stats['total_statements'] == 0:
        return 100.0
    return (stats['covered_statements'] / stats['total_statements']) * 100.0

def main():
    cover_file = 'cover.out'
    
    print("=" * 80)
    print("覆盖率分析报告")
    print("=" * 80)
    print()
    
    file_stats = parse_cover_out(cover_file)
    
    # 按覆盖率排序
    files_by_coverage = []
    for file_path, stats in file_stats.items():
        coverage = calculate_coverage(stats)
        files_by_coverage.append((file_path, coverage, stats))
    
    files_by_coverage.sort(key=lambda x: x[1])
    
    # 找出低于95%的文件
    low_coverage_files = [(f, c, s) for f, c, s in files_by_coverage if c < 95.0 and c > 0]
    
    print(f"📊 总文件数: {len(files_by_coverage)}")
    print(f"⚠️  低于95%覆盖率的文件数: {len(low_coverage_files)}")
    print()
    
    if low_coverage_files:
        print("=" * 80)
        print("低于95%覆盖率的文件详情:")
        print("=" * 80)
        print()
        
        for file_path, coverage, stats in low_coverage_files:
            print(f"📄 {file_path}")
            print(f"   覆盖率: {coverage:.2f}% ({stats['covered_statements']}/{stats['total_statements']} 语句)")
            
            if stats['uncovered_lines']:
                uncovered = sorted(stats['uncovered_lines'])
                # 合并连续的行号
                ranges = []
                start = uncovered[0]
                end = uncovered[0]
                for line in uncovered[1:]:
                    if line == end + 1:
                        end = line
                    else:
                        if start == end:
                            ranges.append(str(start))
                        else:
                            ranges.append(f"{start}-{end}")
                        start = line
                        end = line
                if start == end:
                    ranges.append(str(start))
                else:
                    ranges.append(f"{start}-{end}")
                
                print(f"   ❌ 未覆盖的行: {', '.join(ranges[:20])}{' ...' if len(ranges) > 20 else ''}")
            
            if stats['low_coverage_lines']:
                low_cov = sorted(stats['low_coverage_lines'])
                ranges = []
                start = low_cov[0]
                end = low_cov[0]
                for line in low_cov[1:]:
                    if line == end + 1:
                        end = line
                    else:
                        if start == end:
                            ranges.append(str(start))
                        else:
                            ranges.append(f"{start}-{end}")
                        start = line
                        end = line
                if start == end:
                    ranges.append(str(start))
                else:
                    ranges.append(f"{start}-{end}")
                
                print(f"   ⚠️  低覆盖率行(≤1次): {', '.join(ranges[:10])}{' ...' if len(ranges) > 10 else ''}")
            
            print()
    
    # 显示总体统计
    total_stats = {
        'total_statements': sum(s['total_statements'] for _, _, s in files_by_coverage),
        'covered_statements': sum(s['covered_statements'] for _, _, s in files_by_coverage),
    }
    overall_coverage = calculate_coverage(total_stats)
    print("=" * 80)
    print(f"📈 总体覆盖率: {overall_coverage:.2f}%")
    print(f"   总语句数: {total_stats['total_statements']}")
    print(f"   已覆盖: {total_stats['covered_statements']}")
    print(f"   未覆盖: {total_stats['total_statements'] - total_stats['covered_statements']}")
    print("=" * 80)

if __name__ == '__main__':
    main()

