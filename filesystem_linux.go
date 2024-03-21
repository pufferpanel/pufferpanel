package pufferpanel

import (
	"errors"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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
	fd, err := unix.Openat2(getFd(sfp.root), path, &unix.OpenHow{
		Flags:   uint64(flags),
		Mode:    uint64(syscallMode(mode)),
		Resolve: unix.RESOLVE_BENEATH,
	})
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), filepath.Base(path))
	return file, nil
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

	err = unix.Renameat2(getFd(sourceFolder), sourceName, getFd(targetFolder), targetName, 0)
	return err
}

func (sfp *fileServer) Mkdir(path string, mode os.FileMode) error {
	path = prepPath(path)
	parent := filepath.Dir(path)
	f := filepath.Base(path)

	folder, err := sfp.OpenFile(parent, os.O_RDONLY, mode)
	if err != nil {
		return err
	}
	defer Close(folder)
	return unix.Mkdirat(getFd(folder), f, syscallMode(mode))
}

func (sfp *fileServer) Remove(path string) error {
	path = prepPath(path)
	parent := filepath.Dir(path)
	f := filepath.Base(path)

	folder, err := sfp.OpenFile(parent, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer Close(folder)

	expected, err := sfp.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	stat, err := expected.Stat()
	Close(expected)

	if stat.IsDir() {
		return unix.Unlinkat(getFd(folder), f, unix.AT_REMOVEDIR)
	} else {
		return unix.Unlinkat(getFd(folder), f, 0)
	}
}

func (sfp *fileServer) RemoveAll(path string) error {
	path = prepPath(path)

	folder, err := sfp.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer Close(folder)

	files, err := folder.ReadDir(0)
	if err != nil {
		return err
	}

	//go through all the files, and call our stuff to delete it
	for _, v := range files {
		if v.Type()&os.ModeSymlink == 0 && v.IsDir() {
			//recursive call, so we need to go into this one and delete things
			err = sfp.RemoveAll(filepath.Join(path, v.Name()))
			if err != nil {
				return err
			}
			err = unix.Unlinkat(getFd(folder), v.Name(), unix.AT_REMOVEDIR)
			if err != nil {
				return err
			}
		} else {
			err = unix.Unlinkat(getFd(folder), v.Name(), 0)
			if err != nil {
				return err
			}
		}
	}

	err = sfp.Remove(path)
	return err
}

func syscallMode(i os.FileMode) (o uint32) {
	o |= uint32(i.Perm())
	if i&os.ModeSetuid != 0 {
		o |= syscall.S_ISUID
	}
	if i&os.ModeSetgid != 0 {
		o |= syscall.S_ISGID
	}
	if i&os.ModeSticky != 0 {
		o |= syscall.S_ISVTX
	}
	return
}

func getFd(f *os.File) int {
	return int(f.Fd())
}

func prepPath(path string) string {
	path = filepath.Clean(path)
	path = strings.TrimPrefix(path, "/")
	return path
}
