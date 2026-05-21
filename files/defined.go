package files

import (
	"os"

	"github.com/pufferpanel/pufferpanel/v3/config"
)

var CacheFS FileServer
var BinariesFS FileServer
var RootServerFS FileServer

func CreateSharedFS() error {
	var err error

	CacheFS, err = NewFileServer(config.CacheFolder.Value(), os.Getuid(), os.Getgid())
	if err != nil {
		return err
	}

	BinariesFS, err = NewFileServer(config.BinariesFolder.Value(), os.Getuid(), os.Getgid())
	if err != nil {
		return err
	}

	RootServerFS, err = NewFileServer(config.ServersFolder.Value(), os.Getuid(), os.Getgid())
	if err != nil {
		return err
	}

	return nil
}
