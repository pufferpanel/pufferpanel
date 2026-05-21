package files

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pufferpanel/pufferpanel/v3/utils"
)

func CopyFile(sourceFS FileServer, src string, destFS FileServer, dest string) error {
	source, err := sourceFS.OpenFile(src, os.O_RDONLY, DefaultFilePermissions)
	if err != nil {
		return err
	}
	defer utils.Close(source)

	err = destFS.MkdirAll(filepath.Dir(dest), DefaultFolderPermissions)
	if err != nil {
		return err
	}
	destination, err := destFS.Create(dest)
	if err != nil {
		return err
	}
	defer utils.Close(destination)
	_, err = io.Copy(destination, source)
	return err
}

func WriteFile(src io.Reader, destFS FileServer, dest string) error {
	destination, err := destFS.Create(dest)
	if err != nil {
		return err
	}
	defer utils.Close(destination)
	_, err = io.Copy(destination, src)
	return err
}
