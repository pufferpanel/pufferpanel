package files

import (
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer utils.Close(source)

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}
	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer utils.Close(destination)
	_, err = io.Copy(destination, source)
	return err
}

func WriteFile(src io.ReadCloser, dest string) error {
	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer utils.Close(destination)
	_, err = io.Copy(destination, src)
	return err
}
