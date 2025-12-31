package fileutil

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/andeya/gust/result"
	"github.com/andeya/gust/void"
)

// FileExist reports whether the named file or directory exists.
// It returns two values: whether the file exists, and whether it's a directory.
func FileExist(name string) (existed bool, isDir bool) {
	info, err := os.Stat(name)
	if err != nil {
		return !os.IsNotExist(err), false
	}
	return true, info.IsDir()
}

// FileExists reports whether the named file or directory exists.
// It returns true if the file exists, false otherwise.
func FileExists(name string) bool {
	existed, _ := FileExist(name)
	return existed
}

// SearchFile searches for a file in the given paths and returns a Result.
// This is often used to search for config files in /etc, ~/, etc.
func SearchFile(filename string, paths ...string) result.Result[string] {
	for _, path := range paths {
		fullpath := filepath.Join(path, filename)
		if FileExists(fullpath) {
			return result.Ok(fullpath)
		}
	}
	return result.TryErr[string](filepath.Join(paths[len(paths)-1], filename) + " not found in paths")
}

// GrepFile searches for lines matching a pattern in a file and returns a Result.
// Newlines are stripped while reading.
func GrepFile(pattern string, filename string) result.Result[[]string] {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return result.TryErr[[]string](err)
	}

	fd, err := os.Open(filename)
	if err != nil {
		return result.TryErr[[]string](err)
	}
	defer fd.Close()

	lines := make([]string, 0)
	reader := bufio.NewReader(fd)
	var prefix string

	for {
		byteLine, isPrefix, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return result.TryErr[[]string](err)
		}
		if err == io.EOF {
			break
		}

		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		}

		line = prefix + line
		prefix = ""

		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return result.Ok(lines)
}

// FilepathSplitExt splits the filename into a pair (root, ext) such that root + ext == filename,
// and ext is empty or begins with a period and contains at most one period.
// Leading periods on the basename are ignored; splitext('.cshrc') returns ("", '.cshrc').
// If slashInsensitive is true, it ignores the difference between slash and backslash.
func FilepathSplitExt(filename string, slashInsensitive ...bool) (root, ext string) {
	insensitive := false
	if len(slashInsensitive) > 0 {
		insensitive = slashInsensitive[0]
	}
	if insensitive {
		filename = FilepathSlashInsensitive(filename)
	}
	for i := len(filename) - 1; i >= 0 && !os.IsPathSeparator(filename[i]); i-- {
		if filename[i] == '.' {
			return filename[:i], filename[i:]
		}
	}
	return filename, ""
}

// FilepathStem returns the stem of filename.
// Example:
//
//	FilepathStem("/root/dir/sub/file.ext") // output "file"
//
// If slashInsensitive is true, it ignores the difference between slash and backslash.
func FilepathStem(filename string, slashInsensitive ...bool) string {
	insensitive := false
	if len(slashInsensitive) > 0 {
		insensitive = slashInsensitive[0]
	}
	if insensitive {
		filename = FilepathSlashInsensitive(filename)
	}
	base := filepath.Base(filename)
	for i := len(base) - 1; i >= 0; i-- {
		if base[i] == '.' {
			return base[:i]
		}
	}
	return base
}

// FilepathSlashInsensitive ignores the difference between the slash and the backslash,
// and converts to the same as the current system.
func FilepathSlashInsensitive(path string) string {
	if filepath.Separator == '/' {
		return strings.ReplaceAll(path, "\\", "/")
	}
	return strings.ReplaceAll(path, "/", "\\")
}

// FilepathContains checks if the basepath contains all the subpaths and returns a VoidResult.
func FilepathContains(basepath string, subpaths []string) (r result.VoidResult) {
	defer r.Catch()
	baseAbs := result.Ret(filepath.Abs(basepath)).Unwrap()
	iterator.FromSlice(subpaths).ForEach(func(p string) {
		pAbs := result.Ret(filepath.Abs(p)).Unwrap()
		rel := result.Ret(filepath.Rel(baseAbs, pAbs)).Unwrap()
		if strings.HasPrefix(rel, "..") {
			panic(result.FmtErrVoid("%s is not include %s", baseAbs, pAbs))
		}
	})
	return result.OkVoid()
}

// FilepathAbsolute converts all paths to absolute paths and returns a Result.
func FilepathAbsolute(paths []string) (r result.Result[[]string]) {
	defer r.Catch()
	res := iterator.FromSlice(paths).Map(func(p string) string {
		return result.Ret(filepath.Abs(p)).Unwrap()
	}).Collect()
	return result.Ok(res)
}

// FilepathAbsoluteMap converts all paths to absolute paths and returns a Result.
func FilepathAbsoluteMap(paths []string) (r result.Result[map[string]string]) {
	defer r.Catch()
	res := iterator.Fold(
		iterator.Zip(
			iterator.FromSlice(paths),
			iterator.FromSlice(paths).Map(func(p string) string {
				return result.Ret(filepath.Abs(p)).Unwrap()
			}),
		),
		make(map[string]string, len(paths)),
		func(acc map[string]string, p pair.Pair[string, string]) map[string]string {
			acc[p.A] = p.B
			return acc
		},
	)
	return result.Ok(res)
}

// FilepathRelative converts all target paths to relative paths from the base path and returns a Result.
func FilepathRelative(basepath string, targpaths []string) (r result.Result[[]string]) {
	defer r.Catch()
	baseAbs := result.Ret(filepath.Abs(basepath)).Unwrap()
	res := iterator.FromSlice(targpaths).Map(func(p string) string {
		return filepathRelative(baseAbs, p).Unwrap()
	}).Collect()
	return result.Ok(res)
}

// FilepathRelativeMap converts all target paths to relative paths from the base path and returns a Result.
func FilepathRelativeMap(basepath string, targpaths []string) (r result.Result[map[string]string]) {
	defer r.Catch()
	baseAbs := result.Ret(filepath.Abs(basepath)).Unwrap()
	res := iterator.Fold(
		iterator.Zip(
			iterator.FromSlice(targpaths),
			iterator.FromSlice(targpaths).Map(func(p string) string {
				return filepathRelative(baseAbs, p).Unwrap()
			}),
		),
		make(map[string]string, len(targpaths)),
		func(acc map[string]string, p pair.Pair[string, string]) map[string]string {
			acc[p.A] = p.B
			return acc
		},
	)
	return result.Ok(res)
}

func filepathRelative(basepath, targpath string) (r result.Result[string]) {
	defer r.Catch()
	abs := result.Ret(filepath.Abs(targpath)).Unwrap()
	rel := result.Ret(filepath.Rel(basepath, abs)).Unwrap()
	if strings.HasPrefix(rel, "..") {
		return result.FmtErr[string]("%s is not include %s", basepath, abs)
	}
	return result.Ok(rel)
}

// FilepathDistinct removes duplicate paths and returns a Result.
// If toAbs is true, returns absolute paths; otherwise returns original paths.
func FilepathDistinct(paths []string, toAbs bool) (r result.Result[[]string]) {
	defer r.Catch()
	seen := make(map[string]bool, len(paths))
	res := iterator.FromSlice(paths).FilterMap(func(p string) option.Option[string] {
		abs := result.Ret(filepath.Abs(p)).Unwrap()
		if seen[abs] {
			return option.None[string]()
		}
		seen[abs] = true
		if toAbs {
			return option.Some(abs)
		}
		return option.Some(p)
	}).Collect()
	return result.Ok(res)
}

// FilepathSame checks if the two paths refer to the same file or directory and returns a Result.
func FilepathSame(path1, path2 string) (r result.Result[bool]) {
	defer r.Catch()
	if path1 == path2 {
		return result.Ok(true)
	}
	p1 := result.Ret(filepath.Abs(path1)).Unwrap()
	p2 := result.Ret(filepath.Abs(path2)).Unwrap()
	return result.Ok(p1 == p2)
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns a VoidResult.
// The permission bits perm (before umask) are used for all directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing and returns Ok.
// If perm is empty, default use 0755.
func MkdirAll(path string, perm ...os.FileMode) result.VoidResult {
	var fm os.FileMode = 0755
	if len(perm) > 0 {
		fm = perm[0]
	}
	return result.RetVoid(os.MkdirAll(path, fm))
}

// WriteFile writes data to a file, and automatically creates the directory if necessary.
// If perm is empty, automatically determines the file permissions based on extension.
func WriteFile(filename string, data []byte, perm ...os.FileMode) result.VoidResult {
	filename = filepath.FromSlash(filename)
	if errRes := MkdirAll(filepath.Dir(filename)); errRes.IsErr() {
		return errRes
	}
	if len(perm) > 0 {
		return result.RetVoid(os.WriteFile(filename, data, perm[0]))
	}
	var ext string
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		ext = filename[idx:]
	}
	switch ext {
	case ".sh", ".py", ".rb", ".bat", ".com", ".vbs", ".htm", ".run", ".App", ".exe", ".reg":
		return result.RetVoid(os.WriteFile(filename, data, 0755))
	default:
		return result.RetVoid(os.WriteFile(filename, data, 0644))
	}
}

// RewriteFile rewrites the file content using the provided function and returns a VoidResult.
func RewriteFile(filename string, fn func(content []byte) result.Result[[]byte]) result.VoidResult {
	f, err := os.OpenFile(filename, os.O_RDWR, 0777)
	if err != nil {
		return result.TryErr[void.Void](err)
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return result.TryErr[void.Void](err)
	}
	newContentRes := fn(content)
	if newContentRes.IsErr() {
		return result.TryErr[void.Void](newContentRes.Err())
	}
	newContent := newContentRes.Unwrap()
	if bytes.Equal(content, newContent) {
		return result.OkVoid()
	}
	f.Seek(0, 0)
	f.Truncate(0)
	_, err = f.Write(newContent)
	return result.RetVoid(err)
}

// RewriteToFile rewrites the file to a new filename using the provided function and returns a VoidResult.
// If newFilename already exists and is not a directory, replaces it.
func RewriteToFile(filename, newFilename string, fn func(content []byte) result.Result[[]byte]) result.VoidResult {
	f, err := os.Open(filename)
	if err != nil {
		return result.TryErr[void.Void](err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return result.TryErr[void.Void](err)
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return result.TryErr[void.Void](err)
	}
	newContentRes := fn(content)
	if newContentRes.IsErr() {
		return result.TryErr[void.Void](newContentRes.Err())
	}
	return WriteFile(newFilename, newContentRes.Unwrap(), info.Mode())
}

// ReplaceFile replaces the bytes selected by [start, end] with the new content and returns a VoidResult.
func ReplaceFile(filename string, start, end int, newContent string) result.VoidResult {
	if start < 0 || (end >= 0 && start > end) {
		return result.OkVoid()
	}
	return RewriteFile(filename, func(content []byte) result.Result[[]byte] {
		if end < 0 || end > len(content) {
			end = len(content)
		}
		if start > end {
			start = end
		}
		return result.Ok(bytes.Replace(content, content[start:end], []byte(newContent), 1))
	})
}

// CopyFile copies a single file from src to dst and returns a VoidResult.
func CopyFile(src, dst string) result.VoidResult {
	srcfd, err := os.Open(src)
	if err != nil {
		return result.TryErr[void.Void](err)
	}
	defer srcfd.Close()

	dstfd, err := os.Create(dst)
	if err != nil {
		return result.TryErr[void.Void](err)
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return result.TryErr[void.Void](err)
	}
	srcinfo, err := os.Stat(src)
	if err != nil {
		return result.TryErr[void.Void](err)
	}
	return result.RetVoid(os.Chmod(dst, srcinfo.Mode()))
}

// CopyDir copies a whole directory recursively and returns a VoidResult.
func CopyDir(src string, dst string) (r result.VoidResult) {
	defer r.Catch()
	srcinfo := result.Ret(os.Stat(src)).Unwrap()
	result.RetVoid(os.MkdirAll(dst, srcinfo.Mode())).Unwrap()
	fds := result.Ret(os.ReadDir(src)).Unwrap()
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			CopyDir(srcfp, dstfp).Unwrap()
		} else {
			CopyFile(srcfp, dstfp).Unwrap()
		}
	}
	return result.OkVoid()
}
