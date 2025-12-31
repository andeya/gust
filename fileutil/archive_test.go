package fileutil

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTarGzTo(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create test files
	file1 := filepath.Join(srcDir, "file1.txt")
	file2 := filepath.Join(srcDir, "file2.txt")
	subDir := filepath.Join(srcDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub directory: %v", err)
	}
	file3 := filepath.Join(subDir, "file3.txt")

	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}
	if err := os.WriteFile(file3, []byte("content3"), 0644); err != nil {
		t.Fatalf("Failed to create file3: %v", err)
	}

	// Test archiving without prefix
	var buf bytes.Buffer
	res := TarGzTo(srcDir, &buf, false, nil)
	if res.IsErr() {
		t.Fatalf("TarGzTo should not return error: %v", res.Err())
	}

	// Verify archive contents
	gr, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	files := make(map[string]string)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to read tar entry: %v", err)
		}
		if hdr.Typeflag == tar.TypeReg {
			var content bytes.Buffer
			if _, err := io.Copy(&content, tr); err != nil {
				t.Fatalf("Failed to read file content: %v", err)
			}
			files[hdr.Name] = content.String()
		}
	}

	expectedFiles := map[string]string{
		"file1.txt":     "content1",
		"file2.txt":     "content2",
		"sub/file3.txt": "content3",
	}
	for name, expectedContent := range expectedFiles {
		if content, ok := files[name]; !ok {
			t.Errorf("File %s not found in archive", name)
		} else if content != expectedContent {
			t.Errorf("File %s has wrong content: got %q, want %q", name, content, expectedContent)
		}
	}

	// Test archiving with prefix
	buf.Reset()
	res = TarGzTo(srcDir, &buf, true, nil)
	if res.IsErr() {
		t.Fatalf("TarGzTo should not return error: %v", res.Err())
	}

	// Verify archive contents with prefix
	gr, err = gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	tr = tar.NewReader(gr)
	files = make(map[string]string)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to read tar entry: %v", err)
		}
		if hdr.Typeflag == tar.TypeReg {
			var content bytes.Buffer
			if _, err := io.Copy(&content, tr); err != nil {
				t.Fatalf("Failed to read file content: %v", err)
			}
			files[hdr.Name] = content.String()
		}
	}
	// With includePrefix=true, files should still be in archive
	if len(files) == 0 {
		t.Error("TarGzTo with includePrefix=true should include files")
	}

	// Test archiving single file
	file1Abs, _ := filepath.Abs(file1)
	buf.Reset()
	res = TarGzTo(file1Abs, &buf, false, nil)
	if res.IsErr() {
		t.Fatalf("TarGzTo should not return error for single file: %v", res.Err())
	}

	// Test with ignore elements
	buf.Reset()
	res = TarGzTo(srcDir, &buf, false, nil, "file1.txt", "sub")
	if res.IsErr() {
		t.Fatalf("TarGzTo should not return error: %v", res.Err())
	}

	gr, err = gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	tr = tar.NewReader(gr)
	files = make(map[string]string)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to read tar entry: %v", err)
		}
		if hdr.Typeflag == tar.TypeReg {
			var content bytes.Buffer
			if _, err := io.Copy(&content, tr); err != nil {
				t.Fatalf("Failed to read file content: %v", err)
			}
			files[hdr.Name] = content.String()
		}
	}

	// file1.txt and sub/ should be ignored
	if _, ok := files["file1.txt"]; ok {
		t.Error("file1.txt should be ignored")
	}
	if _, ok := files["sub/file3.txt"]; ok {
		t.Error("sub/file3.txt should be ignored")
	}
	if content, ok := files["file2.txt"]; !ok {
		t.Error("file2.txt should not be ignored")
	} else if content != "content2" {
		t.Errorf("file2.txt has wrong content: got %q, want %q", content, "content2")
	}
}

func TestTarGz(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	dstFile := filepath.Join(tmpDir, "archive.tar.gz")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	file1 := filepath.Join(srcDir, "file1.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	// Test creating archive
	res := TarGz(srcDir, dstFile, false, nil)
	if res.IsErr() {
		t.Fatalf("TarGz should not return error: %v", res.Err())
	}

	// Verify archive file exists
	if _, err := os.Stat(dstFile); err != nil {
		t.Fatalf("TarGz should create archive file: %v", err)
	}

	// Test with log output
	dstFile2 := filepath.Join(tmpDir, "archive2.tar.gz")
	var logOutput []string
	res = TarGz(srcDir, dstFile2, false, func(format string, args ...interface{}) {
		logOutput = append(logOutput, strings.TrimSpace(strings.ReplaceAll(fmt.Sprintf(format, args...), "\n", "")))
	})
	if res.IsErr() {
		t.Fatalf("TarGz should not return error: %v", res.Err())
	}
	if len(logOutput) == 0 {
		t.Error("TarGz should call logOutput function")
	}

	// Test with non-existing source
	res = TarGz(filepath.Join(tmpDir, "nonexistent"), dstFile, false, nil)
	if !res.IsErr() {
		t.Error("TarGz should return error for non-existing source")
	}

	// Test error handling - create archive in non-writable location
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0444); err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755)
	dstFile3 := filepath.Join(readOnlyDir, "archive.tar.gz")
	res = TarGz(srcDir, dstFile3, false, nil)
	if !res.IsErr() {
		t.Error("TarGz should return error when cannot create destination file")
	}
}

func TestTarGzToWithLogOutput(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	file1 := filepath.Join(srcDir, "file1.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	var logOutput []string
	var buf bytes.Buffer
	res := TarGzTo(srcDir, &buf, false, func(format string, args ...interface{}) {
		logOutput = append(logOutput, fmt.Sprintf(format, args...))
	})
	if res.IsErr() {
		t.Fatalf("TarGzTo should not return error: %v", res.Err())
	}
	if len(logOutput) == 0 {
		t.Error("TarGzTo should call logOutput function")
	}
}

func TestTarGzToErrorHandling(t *testing.T) {
	// Test with non-existing source
	var buf bytes.Buffer
	res := TarGzTo("/nonexistent/path", &buf, false, nil)
	if !res.IsErr() {
		t.Error("TarGzTo should return error for non-existing source")
	}

	// Test with invalid source path
	res = TarGzTo(string([]byte{0}), &buf, false, nil)
	if !res.IsErr() {
		t.Error("TarGzTo should return error for invalid source path")
	}
}

func TestTarGzToIgnorePatternVariations(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create files with various patterns
	file1 := filepath.Join(srcDir, "ignore.txt")
	file2 := filepath.Join(srcDir, "keep.txt")
	subDir := filepath.Join(srcDir, "ignore_dir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub directory: %v", err)
	}
	file3 := filepath.Join(subDir, "file.txt")
	file4 := filepath.Join(srcDir, "ignore_dir", "file.txt")

	if err := os.WriteFile(file1, []byte("ignore1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("keep"), 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}
	if err := os.WriteFile(file3, []byte("ignore3"), 0644); err != nil {
		t.Fatalf("Failed to create file3: %v", err)
	}
	if err := os.WriteFile(file4, []byte("ignore4"), 0644); err != nil {
		t.Fatalf("Failed to create file4: %v", err)
	}

	var buf bytes.Buffer
	separator := string(filepath.Separator)
	// Test with ignore pattern that matches prefix
	res := TarGzTo(srcDir, &buf, false, nil, "ignore_dir")
	if res.IsErr() {
		t.Fatalf("TarGzTo should not return error: %v", res.Err())
	}

	gr, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	files := make(map[string]bool)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to read tar entry: %v", err)
		}
		if hdr.Typeflag == tar.TypeReg {
			files[hdr.Name] = true
		}
	}

	// Test various ignore patterns
	testCases := []struct {
		name           string
		ignoreElem     []string
		shouldExist    []string
		shouldNotExist []string
	}{
		{
			name:           "prefix match",
			ignoreElem:     []string{"ignore_dir"},
			shouldExist:    []string{"keep.txt"},
			shouldNotExist: []string{"ignore_dir" + separator + "file.txt"},
		},
		{
			name:           "exact match",
			ignoreElem:     []string{"ignore.txt"},
			shouldExist:    []string{"keep.txt"},
			shouldNotExist: []string{"ignore.txt"},
		},
		{
			name:           "suffix match",
			ignoreElem:     []string{separator + "ignore.txt"},
			shouldExist:    []string{"keep.txt"},
			shouldNotExist: []string{"ignore.txt"},
		},
		{
			name:           "contains match",
			ignoreElem:     []string{"ignore_dir"},
			shouldExist:    []string{"keep.txt"},
			shouldNotExist: []string{"ignore_dir" + separator + "file.txt"},
		},
	}

	for _, tc := range testCases {
		buf.Reset()
		res := TarGzTo(srcDir, &buf, false, nil, tc.ignoreElem...)
		if res.IsErr() {
			t.Fatalf("TarGzTo should not return error for %s: %v", tc.name, res.Err())
		}

		gr, err = gzip.NewReader(&buf)
		if err != nil {
			t.Fatalf("Failed to create gzip reader: %v", err)
		}
		tr = tar.NewReader(gr)
		files = make(map[string]bool)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("Failed to read tar entry: %v", err)
			}
			if hdr.Typeflag == tar.TypeReg {
				files[hdr.Name] = true
			}
		}
		gr.Close()

		for _, name := range tc.shouldExist {
			if !files[name] {
				t.Errorf("File %s should exist in archive for %s", name, tc.name)
			}
		}
		for _, name := range tc.shouldNotExist {
			if files[name] {
				t.Errorf("File %s should not exist in archive for %s", name, tc.name)
			}
		}
	}
}

func TestTarGzToIgnorePatterns(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create files including .DS_Store (should be auto-ignored)
	file1 := filepath.Join(srcDir, "file1.txt")
	dsStore := filepath.Join(srcDir, ".DS_Store")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	if err := os.WriteFile(dsStore, []byte("ds_store"), 0644); err != nil {
		t.Fatalf("Failed to create .DS_Store: %v", err)
	}

	var buf bytes.Buffer
	res := TarGzTo(srcDir, &buf, false, nil)
	if res.IsErr() {
		t.Fatalf("TarGzTo should not return error: %v", res.Err())
	}

	// Verify .DS_Store is ignored
	gr, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	files := make(map[string]bool)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to read tar entry: %v", err)
		}
		if hdr.Typeflag == tar.TypeReg {
			files[hdr.Name] = true
		}
	}

	if files[".DS_Store"] {
		t.Error(".DS_Store should be automatically ignored")
	}
	if !files["file1.txt"] {
		t.Error("file1.txt should not be ignored")
	}
}
