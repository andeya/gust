#!/usr/bin/env python3
import subprocess
import sys
import re
import os

# Get the directory where this script is located
script_dir = os.path.dirname(os.path.abspath(__file__))
# Get the parent directory (gust root)
gust_root = os.path.dirname(script_dir)

# Run go tool cover to get coverage report
result = subprocess.run(
    ['go', 'tool', 'cover', '-func=cover.out'],
    capture_output=True,
    text=True,
    cwd=gust_root
)

if result.returncode != 0:
    print(f"Error running go tool cover: {result.stderr}", file=sys.stderr)
    sys.exit(1)

# Parse coverage output
low_coverage_files = []
for line in result.stdout.split('\n'):
    if not line.strip() or 'total:' in line.lower():
        continue
    
    # Parse line: filename.go:line.function coverage%
    parts = line.rsplit('\t', 2)
    if len(parts) == 3:
        filename = parts[0]
        coverage_str = parts[-1].rstrip('%')
        try:
            coverage = float(coverage_str)
            if 0 < coverage < 95:
                low_coverage_files.append((filename, coverage))
        except ValueError:
            continue

# Sort by coverage (lowest first)
low_coverage_files.sort(key=lambda x: x[1])

# Print results
print("Files with coverage below 95%:\n")
for filename, coverage in low_coverage_files:
    print(f"{filename}: {coverage:.2f}%")

print(f"\nTotal files with low coverage: {len(low_coverage_files)}")
