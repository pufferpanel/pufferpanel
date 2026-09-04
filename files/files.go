package files

import (
	"io"
	"path/filepath"

	"github.com/pufferpanel/pufferpanel/v3/utils"
)

func CopyFile(fs FileServer, src, dest string) error {
	source, err := fs.Open(src)
	if err != nil {
		return err
	}
	defer utils.Close(source)

	err = fs.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}
	destination, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer utils.Close(destination)
	_, err = io.Copy(destination, source)
	return err
}

func WriteFile(fs FileServer, src io.Reader, dest string) error {
	destination, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer utils.Close(destination)
	_, err = io.Copy(destination, src)
	return err
}
