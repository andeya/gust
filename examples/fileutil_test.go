package examples_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andeya/gust/fileutil"
	"github.com/andeya/gust/result"
)

// ExampleFileutil_copyFile demonstrates file copying with Catch pattern.
func Example_fileutil_copyFile() {
	// Before: Traditional Go (multiple error checks)
	// func copyFile(src, dst string) error {
	//     srcFile, err := os.Open(src)
	//     if err != nil {
	//         return err
	//     }
	//     defer srcFile.Close()
	//     dstFile, err := os.Create(dst)
	//     if err != nil {
	//         return err
	//     }
	//     defer dstFile.Close()
	//     _, err = io.Copy(dstFile, srcFile)
	//     return err
	// }

	// After: gust fileutil with Catch pattern (linear flow, automatic error propagation)
	copyFileExample := func(src, dst string) (r result.VoidResult) {
		defer r.Catch()
		fileutil.CopyFile(src, dst).Unwrap()
		return result.OkVoid()
	}

	// Create a temporary file for demonstration
	tmpfile, _ := os.CreateTemp("", "example")
	tmpfile.WriteString("test content")
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	dstFile := tmpfile.Name() + ".copy"
	defer os.Remove(dstFile)

	res := copyFileExample(tmpfile.Name(), dstFile)
	if res.IsOk() {
		fmt.Println("File copied successfully")
	} else {
		fmt.Println("Error:", res.Err())
	}
	// Output: File copied successfully
}

// Example_fileutil_writeFile demonstrates writing files with automatic directory creation.
func Example_fileutil_writeFile() {
	// Before: Traditional Go (manual directory creation, error checks)
	// func writeConfig(filename string, data []byte) error {
	//     dir := filepath.Dir(filename)
	//     if err := os.MkdirAll(dir, 0755); err != nil {
	//         return err
	//     }
	//     return os.WriteFile(filename, data, 0644)
	// }

	// After: gust fileutil (automatic directory creation, Catch pattern)
	writeConfig := func(filename string, data []byte) (r result.VoidResult) {
		defer r.Catch()
		fileutil.WriteFile(filename, data).Unwrap()
		return result.OkVoid()
	}

	// Create a temporary directory for demonstration
	tmpDir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "subdir", "config.json")
	res := writeConfig(configPath, []byte(`{"key": "value"}`))
	if res.IsOk() {
		fmt.Println("Config written successfully")
	} else {
		fmt.Println("Error:", res.Err())
	}
	// Output: Config written successfully
}

// Example_fileutil_copyDirectory demonstrates recursive directory copying with Catch pattern.
func Example_fileutil_copyDirectory() {
	// Before: Traditional Go (multiple error checks, nested conditions)
	// func copyDirectory(src, dst string) error {
	//     info, err := os.Stat(src)
	//     if err != nil {
	//         return err
	//     }
	//     if err = os.MkdirAll(dst, info.Mode()); err != nil {
	//         return err
	//     }
	//     entries, err := os.ReadDir(src)
	//     if err != nil {
	//         return err
	//     }
	//     for _, entry := range entries {
	//         // ... more error checks
	//     }
	//     return nil
	// }

	// After: gust fileutil with Catch pattern (linear flow, single error handler)
	copyDirectoryExample := func(src, dst string) (r result.VoidResult) {
		defer r.Catch()
		fileutil.CopyDir(src, dst).Unwrap()
		return result.OkVoid()
	}

	// Create temporary directories for demonstration
	tmpDir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "test.txt"), []byte("test"), 0644)

	res := copyDirectoryExample(srcDir, dstDir)
	if res.IsOk() {
		fmt.Println("Directory copied successfully")
	} else {
		fmt.Println("Error:", res.Err())
	}
	// Output: Directory copied successfully
}

// Example_fileutil_searchFile demonstrates searching for files in multiple paths.
func Example_fileutil_searchFile() {
	// Before: Traditional Go (manual path iteration, error handling)
	// func findConfig(filename string, paths []string) (string, error) {
	//     for _, path := range paths {
	//         fullpath := filepath.Join(path, filename)
	//         if _, err := os.Stat(fullpath); err == nil {
	//             return fullpath, nil
	//         }
	//     }
	//     return "", fmt.Errorf("file not found")
	// }

	// After: gust fileutil (declarative, Result-based)
	searchConfig := func(filename string, paths ...string) result.Result[string] {
		return fileutil.SearchFile(filename, paths...)
	}

	// Create a temporary file for demonstration
	tmpDir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(tmpDir)

	configFile := "config.json"
	configPath := filepath.Join(tmpDir, configFile)
	os.WriteFile(configPath, []byte(`{}`), 0644)

	// Search in multiple paths
	res := searchConfig(configFile, "/etc", tmpDir, "/usr/local/etc")
	if res.IsOk() {
		foundPath := res.Unwrap()
		// Verify it's in the temp directory
		if filepath.Dir(foundPath) == tmpDir {
			fmt.Println("Found config in temp directory")
		} else {
			fmt.Println("Found config at:", foundPath)
		}
	} else {
		fmt.Println("Config not found")
	}
	// Output: Found config in temp directory
}

// Example_fileutil_fileOperations demonstrates various file operations with Catch pattern.
func Example_fileutil_fileOperations() {
	// Demonstrate multiple file operations in a single function with Catch pattern
	performFileOps := func() (r result.VoidResult) {
		defer r.Catch()

		// Create temporary directory
		tmpDir, _ := os.MkdirTemp("", "example")
		defer os.RemoveAll(tmpDir)

		// Write file (automatically creates directory)
		filePath := filepath.Join(tmpDir, "data", "file.txt")
		fileutil.WriteFile(filePath, []byte("Hello, World!")).Unwrap()

		// Copy file
		copyPath := filepath.Join(tmpDir, "data", "file.copy.txt")
		fileutil.CopyFile(filePath, copyPath).Unwrap()

		// Check if file exists
		exists, _ := fileutil.FileExist(filePath)
		if exists {
			fmt.Println("File operations completed successfully")
		}

		return result.OkVoid()
	}

	res := performFileOps()
	if res.IsOk() {
		fmt.Println("All operations succeeded")
	} else {
		fmt.Println("Error:", res.Err())
	}
	// Output:
	// File operations completed successfully
	// All operations succeeded
}
