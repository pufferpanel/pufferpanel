package files

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pufferpanel/pufferpanel/v3/sys"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"golang.org/x/sys/unix"
)

const O_CREATEFILE = os.O_CREATE | os.O_TRUNC | os.O_RDWR
const DefaultFolderPermissions = 0755
const DefaultFilePermissions = 0644

type FileServer interface {
	fs.FS
	fs.ReadDirFS
	fs.StatFS
	fs.ReadFileFS

	Prefix() string

	Stat(name string) (fs.FileInfo, error)
	Mkdir(path string, mode os.FileMode) error
	MkdirAll(path string, mode os.FileMode) error
	OpenFile(path string, flags int, mode os.FileMode) (*os.File, error)
	OpenRead(name string) (*os.File, error)
	Remove(path string) error
	Rename(source, target string) error
	RemoveAll(path string) error
	Glob(pattern string) ([]string, error)
	Symlink(oldname, newname string) error
	Create(file string) (*os.File, error)
	ReadFile(file string) ([]byte, error)
	WriteFile(file string, data []byte) error
	Chmod(file string, mode os.FileMode) error

	Close() error
}

type fileServer struct {
	dir  string
	root *os.File

	wrapperRoot *os.Root

	uid int
	gid int
}

func NewFileServer(prefix string, uid, gid int) (FileServer, error) {
	f := &fileServer{dir: prefix, uid: uid, gid: gid}
	var err error
	f.root, err = f.resolveRootFd()
	if err != nil {
		return nil, err
	}

	f.wrapperRoot, err = os.OpenRoot(prefix)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (sfp *fileServer) Prefix() string {
	return sfp.dir
}

func (sfp *fileServer) OpenRead(name string) (*os.File, error) {
	return sfp.OpenFile(name, os.O_RDONLY, 0644)
}

func (sfp *fileServer) Open(name string) (fs.File, error) {
	return sfp.OpenRead(name)
}

func (sfp *fileServer) Stat(name string) (fs.FileInfo, error) {
	f, err := sfp.Open(name)
	if err != nil {
		return nil, err
	}
	defer utils.Close(f)
	return f.Stat()
}

func (sfp *fileServer) Symlink(oldpath, newpath string) error {
	return unix.Symlinkat(oldpath, getFd(sfp.root), newpath)
}

func (sfp *fileServer) ReadDir(name string) ([]fs.DirEntry, error) {
	folder, err := sfp.OpenFile(name, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer utils.Close(folder)

	return folder.ReadDir(0)
}

func (sfp *fileServer) Glob(pattern string) ([]string, error) {
	parent := filepath.Base(pattern)
	if parent == pattern {
		parent = "."
	}

	files, err := sfp.ReadDir(parent)

	if err != nil {
		return nil, err
	}

	results := make([]string, 0)
	for _, v := range files {
		if matches, _ := filepath.Match(pattern, v.Name()); matches {
			results = append(results, v.Name())
		}
	}
	return results, nil
}

func (sfp *fileServer) Close() error {
	return sfp.root.Close()
}

func (sfp *fileServer) OpenFile(path string, flags int, mode os.FileMode) (*os.File, error) {
	path = prepPath(path)

	if path == "" {
		return os.Open(sfp.dir)
	}

	//if this is not a create request, nuke mode
	if flags&os.O_CREATE == 0 {
		mode = 0
	}

	var fd int
	var err error
	if utils.UseOpenat2() {
		//at this point, we are going to work on openat2
		fd, err = unix.Openat2(getFd(sfp.root), path, &unix.OpenHow{
			Flags:   uint64(flags),
			Mode:    uint64(sys.SyscallMode(mode)),
			Resolve: unix.RESOLVE_BENEATH,
		})
		if err != nil {
			return nil, err
		}
	} else {
		//because openat is not permitted, we will have to play a game...
		parts := strings.Split(path, string(filepath.Separator))

		//follow the chain, this is just directories we're going through
		var rootFd = getFd(sfp.root)
		var previousFd = rootFd
		for _, v := range parts[:len(parts)-1] {
			fd, err = unix.Openat(previousFd, v, unix.O_NOFOLLOW|unix.O_PATH, sys.SyscallMode(0))
			if previousFd != rootFd {
				_ = unix.Close(previousFd)
			}
			if err != nil {
				return nil, err
			}
			previousFd = fd
		}
		//now.... we can open the file
		fd, err = unix.Openat(previousFd, parts[len(parts)-1], unix.O_NOFOLLOW|flags, sys.SyscallMode(mode))
		if previousFd != rootFd {
			_ = unix.Close(previousFd)
		}
		if err != nil {
			return nil, err
		}
	}

	file := os.NewFile(uintptr(fd), filepath.Base(path))
	if flags&os.O_CREATE == 1 && sfp.uid != -1 {
		err = file.Chown(sfp.uid, sfp.gid)
	}
	return file, err
}

func (sfp *fileServer) Create(path string) (*os.File, error) {
	return sfp.OpenFile(path, O_CREATEFILE, 0644)
}

func (sfp *fileServer) MkdirAll(path string, mode os.FileMode) error {
	//this is going to be recursive...
	path = prepPath(path)

	//now for each one, we just need to make each path, and hope this works
	//in theory, the mkdir will be safe enough
	parts := strings.Split(path, string(filepath.Separator))
	//if it was just mkdir root... we don't do anything
	if len(parts) == 0 {
		return nil
	}

	var err error
	for i := range parts {
		err = sfp.Mkdir(filepath.Join(parts[:i+1]...), mode)
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
	defer utils.Close(sourceFolder)

	targetFolder, err := sfp.OpenFile(targetParent, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer utils.Close(targetFolder)

	err = unix.Renameat2(getFd(sourceFolder), sourceName, getFd(targetFolder), targetName, 0)
	return err
}

func (sfp *fileServer) Mkdir(path string, mode os.FileMode) error {
	path = prepPath(path)
	parent := filepath.Dir(path)
	f := filepath.Base(path)

	if parent == "" {
		err := unix.Mkdirat(getFd(sfp.root), f, sys.SyscallMode(mode))
		if err != nil {
			return err
		}
		if sfp.uid != -1 {
			err = unix.Fchown(getFd(sfp.root), sfp.uid, sfp.gid)
		}

		return err
	} else {
		folder, err := sfp.OpenFile(parent, os.O_RDONLY, mode)
		if err != nil {
			return err
		}
		defer utils.Close(folder)
		err = unix.Mkdirat(getFd(folder), f, sys.SyscallMode(mode))
		if err != nil {
			return err
		}
		if sfp.uid != -1 {
			err = unix.Fchown(getFd(folder), sfp.uid, sfp.gid)
		}
		return err
	}
}

func (sfp *fileServer) Remove(path string) error {
	path = prepPath(path)
	parent := filepath.Dir(path)
	f := filepath.Base(path)

	folder, err := sfp.OpenFile(parent, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer utils.Close(folder)

	expected, err := sfp.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	stat, err := expected.Stat()
	utils.Close(expected)
	if err != nil {
		return err
	}

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
	defer utils.Close(folder)

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

func (sfp *fileServer) WriteFile(path string, data []byte) error {
	file, err := sfp.Create(path)
	defer utils.Close(file)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func (sfp *fileServer) ReadFile(path string) ([]byte, error) {
	file, err := sfp.Open(path)
	defer utils.Close(file)
	if err != nil {
		return nil, err
	}
	reader := io.LimitReader(file, 1024)
	return io.ReadAll(reader)
}

func (sfp *fileServer) Chmod(path string, mode os.FileMode) error {
	return nil
}

func getFd(f *os.File) int {
	return int(f.Fd())
}

func prepPath(path string) string {
	path = filepath.Clean(path)
	path = strings.TrimPrefix(path, "/")

	if path == "." || path == "/" {
		return ""
	}

	return path
}

// shorten maps name, which should start with f.dir, back to the suffix after f.dir.
func (sfp *fileServer) shorten(name string) (rel string, ok bool) {
	if name == sfp.dir {
		return ".", true
	}
	if len(name) >= len(sfp.dir)+2 && name[len(sfp.dir)] == '/' && name[:len(sfp.dir)] == sfp.dir {
		return name[len(sfp.dir)+1:], true
	}
	return "", false
}

// fixErr shortens any reported names in PathErrors by stripping f.dir.
func (sfp *fileServer) fixErr(err error) error {
	var e *fs.PathError
	if errors.As(err, &e) {
		if short, ok := sfp.shorten(e.Path); ok {
			e.Path = short
		}
	}
	return err
}

func (sfp *fileServer) resolveRootFd() (*os.File, error) {
	return os.Open(sfp.dir)
}
