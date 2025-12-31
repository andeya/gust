package fileutil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/andeya/gust/result"
	"github.com/andeya/gust/void"
)

// TarGz compresses and archives files into a tar.gz file and returns a VoidResult.
// src is the source file or directory to archive.
// dst is the destination tar.gz file path.
// includePrefix determines whether to include the parent directory in the archive.
// logOutput is an optional logging function that will be called for each file processed.
// ignoreElem is a list of file or directory names to ignore (e.g., ".git", ".DS_Store").
func TarGz(src, dst string, includePrefix bool, logOutput func(string, ...interface{}), ignoreElem ...string) result.VoidResult {
	return result.AndThen(
		result.Ret(os.Create(dst)),
		func(fw *os.File) result.VoidResult {
			defer fw.Close()
			return TarGzTo(src, fw, includePrefix, logOutput, ignoreElem...).
				OrElse(func(err error) result.VoidResult {
					os.Remove(dst)
					return result.TryErr[void.Void](err)
				})
		},
	)
}

// TarGzTo compresses and archives files to the given writer and returns a VoidResult.
// src is the source file or directory to archive.
// dstWriter is the destination writer.
// includePrefix determines whether to include the parent directory in the archive.
// logOutput is an optional logging function that will be called for each file processed.
// ignoreElem is a list of file or directory names to ignore (e.g., ".git", ".DS_Store").
func TarGzTo(src string, dstWriter io.Writer, includePrefix bool, logOutput func(string, ...interface{}), ignoreElem ...string) result.VoidResult {
	return result.AndThen(
		result.Ret(filepath.Abs(src)),
		func(srcAbs string) result.VoidResult {
			return result.AndThen(
				result.Ret(os.Stat(srcAbs)),
				func(srcFi os.FileInfo) result.VoidResult {
					return tarGzToImpl(srcAbs, srcFi, dstWriter, includePrefix, logOutput, ignoreElem...)
				},
			)
		},
	)
}

func tarGzToImpl(src string, srcFi os.FileInfo, dstWriter io.Writer, includePrefix bool, logOutput func(string, ...interface{}), ignoreElem ...string) result.VoidResult {

	gw := gzip.NewWriter(dstWriter)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	separator := string(filepath.Separator)

	cleanedIgnore := make([]string, 0, len(ignoreElem)+1)
	for _, v := range ignoreElem {
		v = strings.Trim(v, separator)
		if v == "" {
			continue
		}
		cleanedIgnore = append(cleanedIgnore, v)
	}
	ignoreElem = append(cleanedIgnore, ".DS_Store")

	var prefix string
	if !srcFi.IsDir() || includePrefix {
		prefix, _ = filepath.Split(src)
	} else {
		prefix = src + separator
	}

	walkErr := filepath.Walk(src, func(fileName string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}

		// Because hdr.Name is base name,
		// once packaged, all files will pile up and destroy the original directory structure.
		hdr.Name = strings.TrimPrefix(fileName, prefix)

		// ignore files
		for _, v := range ignoreElem {
			if hdr.Name == v ||
				strings.HasPrefix(hdr.Name, v+separator) ||
				strings.HasSuffix(hdr.Name, separator+v) ||
				strings.Contains(hdr.Name, separator+v+separator) {
				return nil
			}
		}

		// If it is not a standard file, it will not be processed, such as a directory.
		if !fi.Mode().IsRegular() {
			return nil
		}

		// write file information
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		fr, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer fr.Close()

		n, err := io.Copy(tw, fr)
		if err != nil {
			return err
		}
		if logOutput != nil {
			logOutput("tar.gz: packaged %s, written %d bytes\n", hdr.Name, n)
		}
		return nil
	})
	if walkErr != nil {
		return result.TryErr[void.Void](walkErr)
	}
	return result.OkVoid()
}
