package pufferpanel

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func (sfp *fileServer) OpenFile(path string, flags int, mode os.FileMode) (*os.File, error) {
	path = prepPath(path)

	if path == "" {
		return os.Open(sfp.dir)
	}

	//if this is not a create request, nuke mode
	if flags&os.O_CREATE == 0 {
		mode = 0
	}

	//at this point, we are going to work on openat2
	fd, err := os.OpenFile(filepath.Join(sfp.dir, path), flags, mode)
	if err != nil {
		return nil, err
	}
	fi, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	//for windows, hard deny symlinks
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return nil, errors.New("access denied")
	}
	return fd, nil
}

func (sfp *fileServer) MkdirAll(path string, mode os.FileMode) error {
	//this is going to be recursive...
	path = prepPath(path)

	//now for each one, we just need to make each path, and hope this works
	//in theory, the mkdir will be safe enough
	parts := strings.Split(path, string(filepath.Separator))
	//if it was just mkdir root... we don't do anything
	if len(parts) <= 1 {
		return nil
	}

	var err error
	for i := range parts {
		err = sfp.Mkdir(filepath.Join(parts[:i]...), mode)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	return nil
}

func (sfp *fileServer) Rename(source, target string) error {
	source = prepPath(source)
	target = prepPath(target)

	sourceParent := filepath.Dir(source)
	targetParent := filepath.Dir(target)

	sourceName := filepath.Base(source)
	targetName := filepath.Base(target)

	sourceFolder, err := sfp.OpenFile(sourceParent, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer Close(sourceFolder)

	targetFolder, err := sfp.OpenFile(targetParent, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer Close(targetFolder)

	err = os.Rename(filepath.Join(sfp.dir, sourceParent, sourceName), filepath.Join(sfp.dir, targetParent, targetName))
	return err
}

func (sfp *fileServer) Mkdir(path string, mode os.FileMode) error {
	path = prepPath(path)
	parent := filepath.Dir(path)

	folder, err := sfp.OpenFile(parent, os.O_RDONLY, mode)
	if err != nil {
		return err
	}
	defer Close(folder)
	return os.Mkdir(filepath.Join(sfp.dir, path), mode)
}

func (sfp *fileServer) Remove(path string) error {
	path = prepPath(path)
	parent := filepath.Dir(path)

	folder, err := sfp.OpenFile(parent, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer Close(folder)

	return os.Remove(filepath.Join(sfp.dir, path))
}

func (sfp *fileServer) RemoveAll(path string) error {
	path = prepPath(path)
	parent := filepath.Dir(path)

	folder, err := sfp.OpenFile(parent, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer Close(folder)

	return os.Remove(filepath.Join(sfp.dir, path))
}

func prepPath(path string) string {
	path = strings.Replace(path, "/", "\\", -1)
	path = filepath.Clean(path)
	path = strings.TrimPrefix(path, "\\")
	return path
}
