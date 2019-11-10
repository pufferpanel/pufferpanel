/*
 Copyright 2018 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package sftp

import (
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
)

type requestPrefix struct {
	prefix string
}

func CreateRequestPrefix(prefix string) sftp.Handlers {
	h := requestPrefix{prefix: prefix}

	return sftp.Handlers{FileCmd: h, FileGet: h, FileList: h, FilePut: h}
}

func (rp requestPrefix) Fileread(request *sftp.Request) (io.ReaderAt, error) {
	file, err := rp.getFile(request.Filepath, os.O_RDONLY, 0644)
	if err != nil {
		logging.Error().Printf("pp-sftp internal error: %s", err)
	}
	return file, err
}

func (rp requestPrefix) Filewrite(request *sftp.Request) (io.WriterAt, error) {
	file, err := rp.getFile(request.Filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	return file, err
}

func (rp requestPrefix) Filecmd(request *sftp.Request) error {
	sourceName, err := rp.validate(request.Filepath)
	if err != nil {
		logging.Error().Printf("pp-sftp internal error: %s", err)
		return rp.maskError(err)
	}
	var targetName string
	if request.Target != "" {
		targetName, err = rp.validate(request.Target)
		if err != nil {
			logging.Error().Printf("pp-sftp internal error: %s", err)
			return rp.maskError(err)
		}
	}
	switch request.Method {
	case "SetStat", "Setstat":
		{
			return nil
		}
	case "Rename":
		{
			return os.Rename(sourceName, targetName)
		}
	case "Rmdir":
		{
			return os.RemoveAll(sourceName)
		}
	case "Mkdir":
		{
			return os.Mkdir(sourceName, 0755)
		}
	case "Symlink":
		{
			return nil
		}
	case "Remove":
		{
			return os.Remove(sourceName)
		}
	default:
		return errors.New(fmt.Sprintf("Unknown request method: %v", request.Method))
	}
}

func (rp requestPrefix) Filelist(request *sftp.Request) (sftp.ListerAt, error) {
	sourceName, err := rp.validate(request.Filepath)
	if err != nil {
		logging.Error().Printf("pp-sftp internal error: %s", err)
		return nil, rp.maskError(err)
	}
	switch request.Method {
	case "List":
		{
			file, err := os.Open(sourceName)
			if err != nil {
				return nil, rp.maskError(err)
			}
			files, err := file.Readdir(0)
			if err != nil {
				return nil, err
			}
			err = file.Close()
			if err != nil {
				return nil, rp.maskError(err)
			}

			//validate any symlinks are valid
			files = pufferpanel.RemoveInvalidSymlinks(files, sourceName, rp.prefix)

			return listerat(files), nil
		}
	case "Stat":
		{
			file, err := os.Open(sourceName)
			if err != nil {
				return nil, rp.maskError(err)
			}
			fi, err := file.Stat()
			if err != nil {
				return nil, rp.maskError(err)
			}
			err = file.Close()
			if err != nil {
				return nil, rp.maskError(err)
			}
			return listerat([]os.FileInfo{fi}), nil
		}
	case "Readlink":
		{
			target, err := os.Readlink(sourceName)
			if err != nil {
				return nil, rp.maskError(err)
			}

			//determine if target is just a local link, or a full link
			//let's just assume linux really at this point
			if !strings.HasPrefix(target, string(os.PathSeparator)) {
				target = rp.prefix + string(os.PathSeparator) + target
			}

			if !pufferpanel.EnsureAccess(target, rp.prefix) {
				return nil, rp.maskError(errors.New("access denied"))
			}

			file, err := os.Open(target)
			if err != nil {
				return nil, rp.maskError(err)
			}
			fi, err := file.Stat()
			if err != nil {
				return nil, rp.maskError(err)
			}
			err = file.Close()
			if err != nil {
				return nil, rp.maskError(err)
			}
			return listerat([]os.FileInfo{fi}), nil
		}
	default:
		return nil, errors.New(fmt.Sprintf("Unknown request method: %s", request.Method))
	}
}

func (rp requestPrefix) getFile(path string, flags int, mode os.FileMode) (*os.File, error) {
	filePath, err := rp.validate(path)
	if err != nil {
		logging.Error().Printf("pp-sftp internal error: %s", err)
		return nil, rp.maskError(err)
	}

	folderPath := filepath.Dir(filePath)

	var file *os.File

	if flags&os.O_CREATE != 0 {
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			err = nil
			err = os.MkdirAll(folderPath, 0755)
			if err != nil {
				logging.Error().Printf("pp-sftp internal error: %s", err)
				return nil, rp.maskError(err)
			}
			file, err = os.Create(filePath)
		} else if err == nil {
			file, err = os.OpenFile(filePath, flags, mode)
		}
	} else {
		file, err = os.OpenFile(filePath, flags, mode)
	}
	if err != nil {
		logging.Error().Printf("pp-sftp internal error: %s", err)
		return nil, rp.maskError(err)
	}

	return file, err
}

func (rp requestPrefix) validate(path string) (string, error) {
	ok, path := rp.tryPrefix(path)
	if !ok {
		return "", errors.New("access denied")
	}
	return path, nil
}

func (rp requestPrefix) tryPrefix(path string) (bool, string) {
	newPath := filepath.Clean(filepath.Join(rp.prefix, path))
	if pufferpanel.EnsureAccess(newPath, rp.prefix) {
		return true, newPath
	} else {
		return false, ""
	}
}

func (rp requestPrefix) stripPrefix(path string) string {
	prefix, err := filepath.Abs(rp.prefix)
	if err != nil {
		prefix = rp.prefix
	}
	newStr := strings.TrimPrefix(path, prefix)
	if len(newStr) == 0 {
		newStr = "/"
	}
	return newStr
}

func (rp requestPrefix) maskError(err error) error {
	return errors.New(rp.stripPrefix(err.Error()))
}

type listerat []os.FileInfo

// Modeled after strings.Reader's ReadAt() implementation
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
