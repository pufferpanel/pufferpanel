package pufferpanel

import (
	"archive/tar"
	"errors"
	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const PathSeparator = string(os.PathSeparator)

type ExtractOptions struct {
	FileServer   FileServer
	SourceFile   string
	TargetPath   string
	Filter       string
	SkipRoot     bool
	ForcedWalker Walker
}

func DetermineIfSingleRoot(sourceFile string) (bool, error) {
	isSingleRoot := false

	var rootName string

	err := archiver.Walk(sourceFile, func(file archiver.File) (err error) {
		name := getCompressedItemName(file)

		if name == "" || name == PathSeparator {
			return
		}
		root := strings.Split(name, PathSeparator)[0]
		if rootName == "" {
			rootName = root
			return nil
		}
		if root != rootName {
			return archiver.ErrStopWalk
		}
		return nil
	})

	if errors.Is(err, archiver.ErrStopWalk) {
		isSingleRoot = false
	}

	return isSingleRoot, err
}

func Extract(fs FileServer, sourceFile, targetPath, filter string, skipRoot bool, forcedType Walker) error {
	if fs != nil {
		sourceFile = filepath.Join(fs.Prefix(), sourceFile)
	}

	if skipRoot {
		var err error
		skipRoot, err = DetermineIfSingleRoot(sourceFile)
		if err != nil {
			return err
		}
	}

	if forcedType != nil {
		return forcedType.Walk(sourceFile, walker(fs, targetPath, filter, skipRoot))
	}

	return archiver.Walk(sourceFile, walker(fs, targetPath, filter, skipRoot))
}

func Compress(fs FileServer, targetFile string, files []string) error {
	if fs != nil {
		p := fs.Prefix()

		targetFile = filepath.Join(p, targetFile)

		for k, v := range files {
			files[k] = filepath.Join(p, v)
		}
	}

	return archiver.Archive(files, targetFile)
}

func walker(fs FileServer, targetPath, filter string, skipRoot bool) archiver.WalkFunc {
	return func(file archiver.File) (err error) {
		path := getCompressedItemName(file)

		if !CompareWildcard(file.Name(), filter) {
			return
		}

		if skipRoot {
			path = strings.Join(strings.Split(path, PathSeparator)[1:], PathSeparator)
		}

		parent := filepath.Join(targetPath, filepath.Dir(path))
		path = filepath.Join(targetPath, path)

		if file.Mode().IsDir() {
			if fs != nil {
				if err = fs.MkdirAll(path, 0755); err != nil {
					return err
				}
			} else {
				if err = os.MkdirAll(path, 0755); err != nil {
					return err
				}
			}
		} else if file.Mode().IsRegular() {
			if fs != nil {
				if err = fs.MkdirAll(parent, 0755); err != nil {
					return err
				}
			} else {
				if err = os.MkdirAll(parent, 0755); err != nil {
					return err
				}
			}
			var outFile *os.File
			if fs != nil {
				outFile, err = fs.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Mode())
			} else {
				outFile, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Mode())
			}

			if err != nil {
				return err
			}
			defer Close(outFile)
			_, err = io.Copy(outFile, file.ReadCloser)
		}

		return
	}
}

// getCompressedItemName Resolves headers in the event the wrapped interface fails
func getCompressedItemName(file archiver.File) string {
	//For certain headers, the actual File interface uses the wrong value
	//Example, ZIP gives the filename, not the full path

	switch v := file.Header.(type) {
	case zip.FileHeader:
		return v.Name
	case *tar.Header:
		return v.Name
	default:
		return file.Name()
	}
}

type Walker interface {
	Walk(archive string, walkFn archiver.WalkFunc) error
}
