package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"io"
	"net"
	"os"
	"path/filepath"
)

type requestPrefix struct {
	fs         files.FileServer
	remoteAddr net.Addr
	serverId   string
}

func CreateRequestPrefix(remoteAddr net.Addr, serverId string, fs files.FileServer) sftp.Handlers {
	h := requestPrefix{fs: fs, serverId: serverId, remoteAddr: remoteAddr}

	return sftp.Handlers{FileCmd: h, FileGet: h, FileList: h, FilePut: h}
}

func (rp requestPrefix) Fileread(request *sftp.Request) (io.ReaderAt, error) {
	rp.log(request)

	file, err := rp.getFile(request.Filepath, os.O_RDONLY, 0644)
	return file, err
}

func (rp requestPrefix) Filewrite(request *sftp.Request) (io.WriterAt, error) {
	rp.log(request)

	file, err := rp.getFile(request.Filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	return file, err
}

func (rp requestPrefix) Filecmd(request *sftp.Request) error {
	rp.log(request)

	switch request.Method {
	case "SetStat", "Setstat":
		{
			return nil
		}
	case "Rename":
		{
			return rp.fs.Rename(request.Filepath, request.Target)
		}
	case "Rmdir":
		{
			return rp.fs.RemoveAll(request.Filepath)
		}
	case "Mkdir":
		{
			return rp.fs.Mkdir(request.Filepath, 0755)
		}
	case "Symlink":
		{
			return nil
		}
	case "Remove":
		{
			return rp.fs.Remove(request.Filepath)
		}
	default:
		return fmt.Errorf("unknown request method: %v", request.Method)
	}
}

func (rp requestPrefix) Filelist(request *sftp.Request) (sftp.ListerAt, error) {
	rp.log(request)

	switch request.Method {
	case "List":
		{
			f, err := rp.fs.ReadDir(request.Filepath)
			if err != nil {
				return nil, err
			}

			return toListerAt(rp.fs, request.Filepath, f), nil
		}
	case "Stat":
		{
			file, err := rp.getFile(request.Filepath, os.O_RDONLY, 0644)
			if err != nil {
				return nil, err
			}
			fi, err := file.Stat()
			if err != nil {
				return nil, err
			}
			err = file.Close()
			if err != nil {
				return nil, err
			}
			return listerat([]os.FileInfo{fi}), nil
		}
	case "Readlink":
		{
			file, err := rp.fs.Open(request.Filepath)
			if err != nil {
				return nil, err
			}
			fi, err := file.Stat()
			if err != nil {
				return nil, err
			}
			err = file.Close()
			if err != nil {
				return nil, err
			}
			return listerat([]os.FileInfo{fi}), nil
		}
	default:
		return nil, fmt.Errorf("unknown request method: %s", request.Method)
	}
}

func (rp requestPrefix) log(request *sftp.Request) {
	if config.SftpDebugLog.Value() {
		logging.Debug.Printf("[SFTP] [%s] [%s] [%s] [%s] ", rp.remoteAddr.String(), request.Method, rp.serverId, request.Filepath)
	}
}

func (rp requestPrefix) getFile(path string, flags int, mode os.FileMode) (*os.File, error) {
	//if this is a file create, then ensure the folder path exists
	if flags&os.O_CREATE != 0 {
		_, err := rp.fs.Stat(path)
		if os.IsNotExist(err) {
			err = nil
			err = rp.fs.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				return nil, err
			}
		}
	}

	file, err := rp.fs.OpenFile(path, flags, mode)

	if err != nil {
		return nil, err
	}

	return file, err
}

type listerat []os.FileInfo

func toListerAt(fs files.FileServer, root string, entries []os.DirEntry) listerat {
	result := listerat{}

	for _, v := range entries {
		file, err := fs.Stat(filepath.Join(root, v.Name()))
		if err == nil {
			result = append(result, file)
		}
	}

	return result
}

// ListAt Modeled after strings.Reader's ReadAt() implementation
func (f listerat) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	var n int
	if offset >= int64(len(f)) {
		return 0, io.EOF
	}
	n = copy(ls, f[offset:])
	if n < len(ls) {
		return n, io.EOF
	}
	return n, nil
}
