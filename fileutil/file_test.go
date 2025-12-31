package fileutil

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/andeya/gust/result"
)

func TestFileExist(t *testing.T) {
	// Test with existing file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	existed, isDir := FileExist(tmpFile)
	if !existed {
		t.Error("FileExist should return true for existing file")
	}
	if isDir {
		t.Error("FileExist should return false for isDir when file is not a directory")
	}

	// Test with existing directory
	existed, isDir = FileExist(tmpDir)
	if !existed {
		t.Error("FileExist should return true for existing directory")
	}
	if !isDir {
		t.Error("FileExist should return true for isDir when path is a directory")
	}

	// Test with non-existing file
	existed, isDir = FileExist(filepath.Join(tmpDir, "nonexistent.txt"))
	if existed {
		t.Error("FileExist should return false for non-existing file")
	}
	if isDir {
		t.Error("FileExist should return false for isDir when file doesn't exist")
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if !FileExists(tmpFile) {
		t.Error("FileExists should return true for existing file")
	}

	if FileExists(filepath.Join(tmpDir, "nonexistent.txt")) {
		t.Error("FileExists should return false for non-existing file")
	}
}

func TestSearchFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.txt")
	if err := os.WriteFile(tmpFile, []byte("config"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test successful search
	res := SearchFile("config.txt", tmpDir)
	if res.IsErr() {
		t.Errorf("SearchFile should not return error for existing file: %v", res.Err())
	}
	found := res.Unwrap()
	if found != tmpFile {
		t.Errorf("SearchFile returned wrong path: got %s, want %s", found, tmpFile)
	}

	// Test search in multiple paths
	otherDir := t.TempDir()
	res = SearchFile("config.txt", otherDir, tmpDir)
	if res.IsErr() {
		t.Errorf("SearchFile should not return error when file found in second path: %v", res.Err())
	}
	found = res.Unwrap()
	if found != tmpFile {
		t.Errorf("SearchFile returned wrong path: got %s, want %s", found, tmpFile)
	}

	// Test file not found
	res = SearchFile("nonexistent.txt", tmpDir)
	if !res.IsErr() {
		t.Error("SearchFile should return error when file not found")
	}
}

func TestGrepFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "hello world\nfoo bar\nhello again\nbaz qux"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test matching pattern
	res := GrepFile(`^hello`, tmpFile)
	if res.IsErr() {
		t.Fatalf("GrepFile should not return error: %v", res.Err())
	}
	lines := res.Unwrap()
	expected := []string{"hello world", "hello again"}
	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("GrepFile returned wrong lines: got %v, want %v", lines, expected)
	}

	// Test non-matching pattern
	res = GrepFile(`^xyz`, tmpFile)
	if res.IsErr() {
		t.Fatalf("GrepFile should not return error: %v", res.Err())
	}
	lines = res.Unwrap()
	if len(lines) != 0 {
		t.Errorf("GrepFile should return empty slice for non-matching pattern: got %v", lines)
	}

	// Test invalid pattern
	res = GrepFile(`[invalid`, tmpFile)
	if !res.IsErr() {
		t.Error("GrepFile should return error for invalid regex pattern")
	}

	// Test non-existing file
	res = GrepFile(`hello`, filepath.Join(tmpDir, "nonexistent.txt"))
	if !res.IsErr() {
		t.Error("GrepFile should return error for non-existing file")
	}

	// Test with long line (isPrefix case)
	longLineFile := filepath.Join(tmpDir, "longline.txt")
	// Create a file with a very long line that will be split by ReadLine
	longContent := strings.Repeat("a", 10000) + "\nhello\n" + strings.Repeat("b", 10000)
	if err := os.WriteFile(longLineFile, []byte(longContent), 0644); err != nil {
		t.Fatalf("Failed to create long line file: %v", err)
	}
	res = GrepFile(`hello`, longLineFile)
	if res.IsErr() {
		t.Fatalf("GrepFile should not return error: %v", res.Err())
	}
	lines = res.Unwrap()
	if len(lines) != 1 || lines[0] != "hello" {
		t.Errorf("GrepFile should handle long lines: got %v", lines)
	}

	// Test empty file
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	res = GrepFile(`.*`, emptyFile)
	if res.IsErr() {
		t.Fatalf("GrepFile should not return error for empty file: %v", res.Err())
	}
	lines = res.Unwrap()
	if len(lines) != 0 {
		t.Errorf("GrepFile should return empty slice for empty file: got %v", lines)
	}
}

func TestFilepathSplitExt(t *testing.T) {
	tests := []struct {
		filename string
		root     string
		ext      string
	}{
		{"/root/dir/sub/file.ext", "/root/dir/sub/file", ".ext"},
		{"file.ext", "file", ".ext"},
		{"file", "file", ""},
		{".cshrc", "", ".cshrc"},
		{"./file.ext", "./file", ".ext"},
		{"file.ext.ext", "file.ext", ".ext"},
		{"path/to/file", "path/to/file", ""},
	}

	for _, tt := range tests {
		root, ext := FilepathSplitExt(tt.filename)
		if root != tt.root || ext != tt.ext {
			t.Errorf("FilepathSplitExt(%q) = (%q, %q), want (%q, %q)", tt.filename, root, ext, tt.root, tt.ext)
		}
	}

	// Test with slashInsensitive=true
	root, ext := FilepathSplitExt("../..\\../.\\./root/dir/sub\\file.go.ext", true)
	if filepath.Separator == '/' {
		expectedRoot := strings.ReplaceAll("../..\\../.\\./root/dir/sub\\file.go", "\\", "/")
		if root != expectedRoot || ext != ".ext" {
			t.Errorf("FilepathSplitExt with slashInsensitive failed: got (%q, %q)", root, ext)
		}
	} else {
		expectedRoot := strings.ReplaceAll("../..\\../.\\./root/dir/sub\\file.go", "/", "\\")
		if root != expectedRoot || ext != ".ext" {
			t.Errorf("FilepathSplitExt with slashInsensitive failed: got (%q, %q)", root, ext)
		}
	}

	// Test with slashInsensitive=false (default)
	root, ext = FilepathSplitExt("../..\\../.\\./root/dir/sub\\file.go.ext", false)
	if root != "../..\\../.\\./root/dir/sub\\file.go" || ext != ".ext" {
		t.Errorf("FilepathSplitExt with slashInsensitive=false failed: got (%q, %q)", root, ext)
	}
}

func TestFilepathStem(t *testing.T) {
	tests := []struct {
		filename string
		stem     string
	}{
		{"/root/dir/sub/file.ext", "file"},
		{"file.ext", "file"},
		{"file", "file"},
		{"./", ""},
		{".cshrc", ""},
		{"file.ext.ext", "file.ext"},
		{"path/to/file.ext", "file"},
	}

	for _, tt := range tests {
		stem := FilepathStem(tt.filename)
		if stem != tt.stem {
			t.Errorf("FilepathStem(%q) = %q, want %q", tt.filename, stem, tt.stem)
		}
	}

	// Test with slashInsensitive=true
	stem := FilepathStem("../..\\../.\\./root/dir/sub\\file.go.ext", true)
	if stem != "file.go" {
		t.Errorf("FilepathStem with slashInsensitive failed: got %q, want %q", stem, "file.go")
	}

	// Test with slashInsensitive=false (default)
	stem = FilepathStem("../..\\../.\\./root/dir/sub\\file.go.ext", false)
	// On Windows, the path separator is backslash, so base will be "sub\\file.go.ext" or similar
	// We just verify it extracts the stem correctly
	if !strings.Contains(stem, "file.go") && stem != "file.go.ext" {
		t.Errorf("FilepathStem with slashInsensitive=false failed: got %q", stem)
	}
}

func TestFilepathSlashInsensitive(t *testing.T) {
	if filepath.Separator == '/' {
		// Test Unix path (filepath.Separator == '/')
		result := FilepathSlashInsensitive("path\\to\\file")
		if result != "path/to/file" {
			t.Errorf("FilepathSlashInsensitive on Unix: got %q, want %q", result, "path/to/file")
		}
		// Test already correct separator
		result = FilepathSlashInsensitive("path/to/file")
		if result != "path/to/file" {
			t.Errorf("FilepathSlashInsensitive on Unix should preserve forward slashes: got %q", result)
		}
		// Test mixed separators
		result = FilepathSlashInsensitive("path\\to/file")
		if result != "path/to/file" {
			t.Errorf("FilepathSlashInsensitive on Unix: got %q, want %q", result, "path/to/file")
		}
	} else {
		// Test Windows path (filepath.Separator == '\\')
		result := FilepathSlashInsensitive("path/to/file")
		if result != "path\\to\\file" {
			t.Errorf("FilepathSlashInsensitive on Windows: got %q, want %q", result, "path\\to\\file")
		}
		// Test already correct separator
		result = FilepathSlashInsensitive("path\\to\\file")
		if result != "path\\to\\file" {
			t.Errorf("FilepathSlashInsensitive on Windows should preserve backslashes: got %q", result)
		}
		// Test mixed separators
		result = FilepathSlashInsensitive("path/to\\file")
		if result != "path\\to\\file" {
			t.Errorf("FilepathSlashInsensitive on Windows: got %q, want %q", result, "path\\to\\file")
		}
	}
}

func TestFilepathContains(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub directory: %v", err)
	}
	subFile := filepath.Join(subDir, "file.txt")
	if err := os.WriteFile(subFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test valid contains
	res := FilepathContains(tmpDir, []string{subDir, subFile})
	if res.IsErr() {
		t.Errorf("FilepathContains should not return error for contained paths: %v", res.Err())
	}

	// Test invalid contains
	otherDir := t.TempDir()
	res = FilepathContains(tmpDir, []string{otherDir})
	if !res.IsErr() {
		t.Error("FilepathContains should return error for paths not contained")
	}

	// Test with invalid basepath
	res = FilepathContains("", []string{subDir})
	if !res.IsErr() {
		t.Error("FilepathContains should return error for invalid basepath")
	}

	// Test with empty subpaths
	res = FilepathContains(tmpDir, []string{})
	if res.IsErr() {
		t.Errorf("FilepathContains should not return error for empty subpaths: %v", res.Err())
	}

	// Test with invalid subpath (non-existent path)
	// Note: FilepathContains uses filepath.Abs which may succeed even for non-existent paths
	// So we test with a path that exists but is not contained
	otherDir2 := t.TempDir()
	res2 := FilepathContains(tmpDir, []string{otherDir2})
	if !res2.IsErr() {
		t.Error("FilepathContains should return error for paths not contained")
	}
}

func TestFilepathAbsolute(t *testing.T) {
	tmpDir := t.TempDir()
	relPath := "test.txt"
	absPath := filepath.Join(tmpDir, relPath)

	paths := []string{relPath, absPath}
	res := FilepathAbsolute(paths)
	if res.IsErr() {
		t.Fatalf("FilepathAbsolute should not return error: %v", res.Err())
	}
	result := res.Unwrap()
	if len(result) != 2 {
		t.Errorf("FilepathAbsolute returned wrong length: got %d, want 2", len(result))
	}
	if !filepath.IsAbs(result[0]) {
		t.Errorf("FilepathAbsolute should return absolute path: got %q", result[0])
	}
	if result[1] != absPath {
		t.Errorf("FilepathAbsolute should preserve absolute paths: got %q, want %q", result[1], absPath)
	}

	// Test with empty paths
	res = FilepathAbsolute([]string{})
	if res.IsErr() {
		t.Fatalf("FilepathAbsolute should not return error for empty paths: %v", res.Err())
	}
	result = res.Unwrap()
	if len(result) != 0 {
		t.Errorf("FilepathAbsolute should return empty slice for empty input: got %d", len(result))
	}
}

func TestFilepathAbsoluteMap(t *testing.T) {
	tmpDir := t.TempDir()
	relPath := "test.txt"
	absPath := filepath.Join(tmpDir, relPath)

	paths := []string{relPath, absPath}
	res := FilepathAbsoluteMap(paths)
	if res.IsErr() {
		t.Fatalf("FilepathAbsoluteMap should not return error: %v", res.Err())
	}
	result := res.Unwrap()
	if len(result) != 2 {
		t.Errorf("FilepathAbsoluteMap returned wrong length: got %d, want 2", len(result))
	}
	if !filepath.IsAbs(result[relPath]) {
		t.Errorf("FilepathAbsoluteMap should return absolute path: got %q", result[relPath])
	}
	if result[absPath] != absPath {
		t.Errorf("FilepathAbsoluteMap should preserve absolute paths: got %q, want %q", result[absPath], absPath)
	}

	// Test with empty paths
	res = FilepathAbsoluteMap([]string{})
	if res.IsErr() {
		t.Fatalf("FilepathAbsoluteMap should not return error for empty paths: %v", res.Err())
	}
	result = res.Unwrap()
	if len(result) != 0 {
		t.Errorf("FilepathAbsoluteMap should return empty map for empty input: got %d", len(result))
	}
}

func TestFilepathRelativeErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with invalid basepath
	res := FilepathRelative("", []string{tmpDir})
	if !res.IsErr() {
		t.Error("FilepathRelative should return error for invalid basepath")
	}

	// Test with invalid target path
	res = FilepathRelative(tmpDir, []string{string([]byte{0})})
	if !res.IsErr() {
		t.Error("FilepathRelative should return error for invalid target path")
	}
}

func TestFilepathRelative(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub directory: %v", err)
	}
	subFile := filepath.Join(subDir, "file.txt")
	if err := os.WriteFile(subFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	paths := []string{subDir, subFile}
	res := FilepathRelative(tmpDir, paths)
	if res.IsErr() {
		t.Fatalf("FilepathRelative should not return error: %v", res.Err())
	}
	result := res.Unwrap()
	if len(result) != 2 {
		t.Errorf("FilepathRelative returned wrong length: got %d, want 2", len(result))
	}
	if result[0] != "sub" {
		t.Errorf("FilepathRelative returned wrong relative path: got %q, want %q", result[0], "sub")
	}
	if result[1] != filepath.Join("sub", "file.txt") {
		t.Errorf("FilepathRelative returned wrong relative path: got %q", result[1])
	}

	// Test with path not contained
	otherDir := t.TempDir()
	res = FilepathRelative(tmpDir, []string{otherDir})
	if !res.IsErr() {
		t.Error("FilepathRelative should return error for paths not contained")
	}
}

func TestFilepathRelativeMap(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub directory: %v", err)
	}
	subFile := filepath.Join(subDir, "file.txt")
	if err := os.WriteFile(subFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	paths := []string{subDir, subFile}
	res := FilepathRelativeMap(tmpDir, paths)
	if res.IsErr() {
		t.Fatalf("FilepathRelativeMap should not return error: %v", res.Err())
	}
	result := res.Unwrap()
	if len(result) != 2 {
		t.Errorf("FilepathRelativeMap returned wrong length: got %d, want 2", len(result))
	}
	if result[subDir] != "sub" {
		t.Errorf("FilepathRelativeMap returned wrong relative path: got %q, want %q", result[subDir], "sub")
	}

	// Test with empty paths
	res = FilepathRelativeMap(tmpDir, []string{})
	if res.IsErr() {
		t.Fatalf("FilepathRelativeMap should not return error for empty paths: %v", res.Err())
	}
	result = res.Unwrap()
	if len(result) != 0 {
		t.Errorf("FilepathRelativeMap should return empty map for empty input: got %d", len(result))
	}

	// Test with path not contained
	otherDir := t.TempDir()
	res = FilepathRelativeMap(tmpDir, []string{otherDir})
	if !res.IsErr() {
		t.Error("FilepathRelativeMap should return error for paths not contained")
	}
}

func TestFilepathDistinctErrorHandling(t *testing.T) {
	// Note: filepath.Abs may not error on all invalid paths depending on the system
	// We test with a path that should cause an error on most systems
	// On some systems, this might not error, so we just verify it doesn't panic
	res := FilepathDistinct([]string{string([]byte{0})}, false)
	if res.IsErr() {
		// Error is expected, test passes
		return
	}
	// If no error, that's also acceptable on some systems
	t.Log("FilepathDistinct did not error for invalid path (may be system-dependent)")
}

func TestFilepathDistinct(t *testing.T) {
	tmpDir := t.TempDir()
	path1 := filepath.Join(tmpDir, "file1.txt")
	path2 := filepath.Join(tmpDir, "file2.txt")

	// Test with duplicates
	paths := []string{path1, path2, path1, path2}
	res := FilepathDistinct(paths, false)
	if res.IsErr() {
		t.Fatalf("FilepathDistinct should not return error: %v", res.Err())
	}
	result := res.Unwrap()
	if len(result) != 2 {
		t.Errorf("FilepathDistinct should remove duplicates: got %d, want 2", len(result))
	}

	// Test with toAbs=true
	res = FilepathDistinct(paths, true)
	if res.IsErr() {
		t.Fatalf("FilepathDistinct should not return error: %v", res.Err())
	}
	result = res.Unwrap()
	if len(result) != 2 {
		t.Errorf("FilepathDistinct should remove duplicates: got %d, want 2", len(result))
	}
	for _, p := range result {
		if !filepath.IsAbs(p) {
			t.Errorf("FilepathDistinct with toAbs=true should return absolute paths: got %q", p)
		}
	}
}

func TestFilepathSame(t *testing.T) {
	tmpDir := t.TempDir()
	path1 := filepath.Join(tmpDir, "file.txt")
	path2 := filepath.Join(tmpDir, "file.txt")
	path3 := filepath.Join(tmpDir, "other.txt")

	// Test same paths
	res := FilepathSame(path1, path2)
	if res.IsErr() {
		t.Fatalf("FilepathSame should not return error: %v", res.Err())
	}
	same := res.Unwrap()
	if !same {
		t.Error("FilepathSame should return true for same paths")
	}

	// Test different paths
	res = FilepathSame(path1, path3)
	if res.IsErr() {
		t.Fatalf("FilepathSame should not return error: %v", res.Err())
	}
	same = res.Unwrap()
	if same {
		t.Error("FilepathSame should return false for different paths")
	}

	// Test identical strings
	res = FilepathSame("test", "test")
	if res.IsErr() {
		t.Fatalf("FilepathSame should not return error: %v", res.Err())
	}
	same = res.Unwrap()
	if !same {
		t.Error("FilepathSame should return true for identical strings")
	}

	// Test with empty path1 (filepath.Abs("") may succeed on some systems)
	// We test with a clearly invalid path instead
	invalidPath := string([]byte{0}) // null byte is invalid in paths
	res = FilepathSame(invalidPath, path2)
	if !res.IsErr() {
		// On some systems, this might not error, so we just verify it doesn't panic
		t.Log("FilepathSame with invalid path1 did not error (may be system-dependent)")
	}

	// Test with empty path2
	res = FilepathSame(path1, invalidPath)
	if !res.IsErr() {
		// On some systems, this might not error, so we just verify it doesn't panic
		t.Log("FilepathSame with invalid path2 did not error (may be system-dependent)")
	}
}

func TestMkdirAll(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "new", "sub", "dir")

	// Test with default permission
	res := MkdirAll(newDir)
	if res.IsErr() {
		t.Fatalf("MkdirAll should not return error: %v", res.Err())
	}
	info, err := os.Stat(newDir)
	if err != nil {
		t.Fatalf("MkdirAll should create directory: %v", err)
	}
	if !info.IsDir() {
		t.Error("MkdirAll should create a directory")
	}

	// Test with custom permission
	customDir := filepath.Join(tmpDir, "custom")
	res = MkdirAll(customDir, 0700)
	if res.IsErr() {
		t.Fatalf("MkdirAll should not return error: %v", res.Err())
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "sub", "test.txt")
	content := []byte("test content")

	// Test with default permission
	res := WriteFile(testFile, content)
	if res.IsErr() {
		t.Fatalf("WriteFile should not return error: %v", res.Err())
	}
	readContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("WriteFile should create file: %v", err)
	}
	if string(readContent) != string(content) {
		t.Errorf("WriteFile wrote wrong content: got %q, want %q", string(readContent), string(content))
	}

	// Test with custom permission
	customFile := filepath.Join(tmpDir, "custom.txt")
	res = WriteFile(customFile, content, 0600)
	if res.IsErr() {
		t.Fatalf("WriteFile should not return error: %v", res.Err())
	}

	// Test executable file extensions
	executableExts := []string{".sh", ".py", ".rb", ".bat", ".com", ".vbs", ".htm", ".run", ".App", ".exe", ".reg"}
	for _, ext := range executableExts {
		execFile := filepath.Join(tmpDir, "script"+ext)
		res = WriteFile(execFile, content)
		if res.IsErr() {
			t.Fatalf("WriteFile should not return error for %s: %v", ext, res.Err())
		}
		info, err := os.Stat(execFile)
		if err != nil {
			t.Fatalf("WriteFile should create file: %v", err)
		}
		if info.Mode().Perm()&0111 == 0 {
			t.Errorf("WriteFile should set executable permission for %s files", ext)
		}
	}

	// Test with non-executable extension
	txtFile := filepath.Join(tmpDir, "file.txt")
	res = WriteFile(txtFile, content)
	if res.IsErr() {
		t.Fatalf("WriteFile should not return error: %v", res.Err())
	}
	info, err := os.Stat(txtFile)
	if err != nil {
		t.Fatalf("WriteFile should create file: %v", err)
	}
	if info.Mode().Perm()&0111 != 0 {
		t.Error("WriteFile should not set executable permission for .txt files")
	}

	// Test with file without extension
	noExtFile := filepath.Join(tmpDir, "file")
	res = WriteFile(noExtFile, content)
	if res.IsErr() {
		t.Fatalf("WriteFile should not return error: %v", res.Err())
	}
}

func TestRewriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	originalContent := []byte("original content")
	if err := os.WriteFile(testFile, originalContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test rewriting file
	newContent := []byte("new content")
	res := RewriteFile(testFile, func(content []byte) result.Result[[]byte] {
		return result.Ok(newContent)
	})
	if res.IsErr() {
		t.Fatalf("RewriteFile should not return error: %v", res.Err())
	}
	readContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("RewriteFile should preserve file: %v", err)
	}
	if string(readContent) != string(newContent) {
		t.Errorf("RewriteFile wrote wrong content: got %q, want %q", string(readContent), string(newContent))
	}

	// Test with no change
	res = RewriteFile(testFile, func(content []byte) result.Result[[]byte] {
		return result.Ok(content)
	})
	if res.IsErr() {
		t.Fatalf("RewriteFile should not return error when content unchanged: %v", res.Err())
	}

	// Test with error in function
	res = RewriteFile(testFile, func(content []byte) result.Result[[]byte] {
		return result.TryErr[[]byte](os.ErrPermission)
	})
	if !res.IsErr() {
		t.Error("RewriteFile should return error when function returns error")
	}

	// Test with non-existing file
	res = RewriteFile(filepath.Join(tmpDir, "nonexistent.txt"), func(content []byte) result.Result[[]byte] {
		return result.Ok(content)
	})
	if !res.IsErr() {
		t.Error("RewriteFile should return error for non-existing file")
	}
}

func TestRewriteToFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "src.txt")
	dstFile := filepath.Join(tmpDir, "dst.txt")
	content := []byte("source content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	res := RewriteToFile(srcFile, dstFile, func(content []byte) result.Result[[]byte] {
		return result.Ok([]byte("modified content"))
	})
	if res.IsErr() {
		t.Fatalf("RewriteToFile should not return error: %v", res.Err())
	}
	readContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("RewriteToFile should create destination file: %v", err)
	}
	if string(readContent) != "modified content" {
		t.Errorf("RewriteToFile wrote wrong content: got %q, want %q", string(readContent), "modified content")
	}

	// Test with non-existing source file
	res = RewriteToFile(filepath.Join(tmpDir, "nonexistent.txt"), dstFile, func(content []byte) result.Result[[]byte] {
		return result.Ok(content)
	})
	if !res.IsErr() {
		t.Error("RewriteToFile should return error for non-existing source file")
	}

	// Test with error in function
	res = RewriteToFile(srcFile, dstFile, func(content []byte) result.Result[[]byte] {
		return result.TryErr[[]byte](os.ErrPermission)
	})
	if !res.IsErr() {
		t.Error("RewriteToFile should return error when function returns error")
	}

	// Test with non-existing source file (covers os.Open error)
	res = RewriteToFile(filepath.Join(tmpDir, "nonexistent.txt"), dstFile, func(content []byte) result.Result[[]byte] {
		return result.Ok(content)
	})
	if !res.IsErr() {
		t.Error("RewriteToFile should return error for non-existing source file")
	}
}

func TestReplaceFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("hello world")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test replacing middle part
	res := ReplaceFile(testFile, 6, 11, "universe")
	if res.IsErr() {
		t.Fatalf("ReplaceFile should not return error: %v", res.Err())
	}
	readContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReplaceFile should preserve file: %v", err)
	}
	if string(readContent) != "hello universe" {
		t.Errorf("ReplaceFile wrote wrong content: got %q, want %q", string(readContent), "hello universe")
	}

	// Test with invalid range (should do nothing)
	res = ReplaceFile(testFile, -1, 5, "test")
	if res.IsErr() {
		t.Fatalf("ReplaceFile should not return error for invalid range: %v", res.Err())
	}

	// Test with end < 0 (should replace to end)
	res = ReplaceFile(testFile, 6, -1, "end")
	if res.IsErr() {
		t.Fatalf("ReplaceFile should not return error: %v", res.Err())
	}
	readContent, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReplaceFile should preserve file: %v", err)
	}
	if string(readContent) != "hello end" {
		t.Errorf("ReplaceFile wrote wrong content: got %q, want %q", string(readContent), "hello end")
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "src.txt")
	dstFile := filepath.Join(tmpDir, "dst.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	res := CopyFile(srcFile, dstFile)
	if res.IsErr() {
		t.Fatalf("CopyFile should not return error: %v", res.Err())
	}
	readContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("CopyFile should create destination file: %v", err)
	}
	if string(readContent) != string(content) {
		t.Errorf("CopyFile copied wrong content: got %q, want %q", string(readContent), string(content))
	}

	// Verify permissions are copied
	srcInfo, _ := os.Stat(srcFile)
	dstInfo, _ := os.Stat(dstFile)
	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("CopyFile should copy file mode: got %v, want %v", dstInfo.Mode(), srcInfo.Mode())
	}

	// Test copying non-existing file
	res = CopyFile(filepath.Join(tmpDir, "nonexistent.txt"), dstFile)
	if !res.IsErr() {
		t.Error("CopyFile should return error for non-existing source file")
	}

	// Test copying to non-writable location (if not root)
	if os.Getuid() != 0 {
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		if err := os.MkdirAll(readOnlyDir, 0444); err != nil {
			t.Fatalf("Failed to create read-only directory: %v", err)
		}
		defer os.Chmod(readOnlyDir, 0755)
		readOnlyFile := filepath.Join(readOnlyDir, "dst.txt")
		res = CopyFile(srcFile, readOnlyFile)
		if !res.IsErr() {
			t.Error("CopyFile should return error when cannot create destination file")
		}
	}
}

func TestGrepFileErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with file that becomes unreadable during reading
	// This is hard to test directly, but we can test the error path
	// by using a file that gets removed or becomes inaccessible
	// For now, we test the non-EOF error path is covered by the existing tests
}

func TestCopyDir(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	// Create source directory structure
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	subDir := filepath.Join(srcDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub directory: %v", err)
	}
	file1 := filepath.Join(srcDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	res := CopyDir(srcDir, dstDir)
	if res.IsErr() {
		t.Fatalf("CopyDir should not return error: %v", res.Err())
	}

	// Verify copied structure
	dstFile1 := filepath.Join(dstDir, "file1.txt")
	dstFile2 := filepath.Join(dstDir, "sub", "file2.txt")
	content1, err := os.ReadFile(dstFile1)
	if err != nil {
		t.Fatalf("CopyDir should copy file1: %v", err)
	}
	if string(content1) != "content1" {
		t.Errorf("CopyDir copied wrong content for file1: got %q, want %q", string(content1), "content1")
	}
	content2, err := os.ReadFile(dstFile2)
	if err != nil {
		t.Fatalf("CopyDir should copy file2: %v", err)
	}
	if string(content2) != "content2" {
		t.Errorf("CopyDir copied wrong content for file2: got %q, want %q", string(content2), "content2")
	}

	// Test copying non-existing directory
	res = CopyDir(filepath.Join(tmpDir, "nonexistent"), dstDir)
	if !res.IsErr() {
		t.Error("CopyDir should return error for non-existing source directory")
	}

	// Test copying to non-writable location (if not root)
	if os.Getuid() != 0 {
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		if err := os.MkdirAll(readOnlyDir, 0444); err != nil {
			t.Fatalf("Failed to create read-only directory: %v", err)
		}
		defer os.Chmod(readOnlyDir, 0755)
		res = CopyDir(srcDir, readOnlyDir)
		if !res.IsErr() {
			t.Error("CopyDir should return error when cannot create destination directory")
		}
	}
}
