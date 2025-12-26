package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output> <input1> [input2] ...\n", os.Args[0])
		os.Exit(1)
	}

	outputFile := os.Args[1]
	inputFiles := os.Args[2:]

	// Coverage entry stores count and statements
	type coverageEntry struct {
		count      int
		statements int
	}

	// Map to store coverage data: file:line -> coverageEntry
	coverage := make(map[string]coverageEntry)
	mode := ""

	// Read all input files
	for _, inputFile := range inputFiles {
		func() {
			file, err := os.Open(inputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", inputFile, err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			firstLine := true
			for scanner.Scan() {
				line := scanner.Text()
				if firstLine {
					// First line is mode
					if strings.HasPrefix(line, "mode:") {
						mode = line
					}
					firstLine = false
					continue
				}

				// Parse coverage line: file.go:start.startCol,end.endCol count statements
				parts := strings.Fields(line)
				if len(parts) < 2 {
					continue
				}

				key := parts[0] // file.go:start.startCol,end.endCol
				var count, statements int
				if len(parts) >= 2 {
					fmt.Sscanf(parts[1], "%d", &count)
				}
				if len(parts) >= 3 {
					fmt.Sscanf(parts[2], "%d", &statements)
				}

				// Take maximum count for the same line, preserve statements
				if existing, exists := coverage[key]; !exists || count > existing.count {
					coverage[key] = coverageEntry{
						count:      count,
						statements: statements,
					}
				} else if count == existing.count && statements > existing.statements {
					// If count is same, take max statements
					coverage[key] = coverageEntry{
						count:      count,
						statements: statements,
					}
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", inputFile, err)
			}
		}()
	}

	// Write merged coverage file
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// Write mode line
	if mode != "" {
		fmt.Fprintln(writer, mode)
	} else {
		fmt.Fprintln(writer, "mode: atomic")
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(coverage))
	for k := range coverage {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Write coverage lines
	for _, key := range keys {
		entry := coverage[key]
		fmt.Fprintf(writer, "%s %d %d\n", key, entry.count, entry.statements)
	}

	fmt.Printf("Successfully merged %d coverage files into %s\n", len(inputFiles), outputFile)
}
