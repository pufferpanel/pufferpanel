package pufferpanel

import (
	"errors"
	"github.com/mholt/archiver/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const PathSeparator = string(os.PathSeparator)

func DetermineIfSingleRoot(sourceFile string) (bool, error) {
	isSingleRoot := false

	var rootName string

	err := archiver.Walk(sourceFile, func(file archiver.File) (err error) {
		if file.Name() == "" || file.Name() == PathSeparator {
			return
		}
		root := strings.Split(file.Name(), PathSeparator)[0]
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

func Extract(fs FileServer, sourceFile, targetPath, filter string, skipRoot bool) error {
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

	return archiver.Walk(sourceFile, func(file archiver.File) (err error) {
		path := file.Name()

		if !CompareWildcard(file.Name(), filter) {
			return
		}

		if skipRoot {
			path = strings.Join(strings.Split(path, PathSeparator)[1:], PathSeparator)
		}

		parent := filepath.Join(targetPath, filepath.Dir(path))
		path = filepath.Join(targetPath, file.Name())

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
	})
}
