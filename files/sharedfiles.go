package files

import "github.com/pufferpanel/pufferpanel/v3/config"

var BinaryFS FileServer
var ServerFS FileServer
var CacheFS FileServer
var BackupFS FileServer

func InitSharedFileSystems() error {
	var err error

	BinaryFS, err = NewFileServer(config.BinariesFolder.Value(), 0, 0, true)
	if err != nil {
		return err
	}

	ServerFS, err = NewFileServer(config.ServersFolder.Value(), 0, 0, false)
	if err != nil {
		return err
	}

	CacheFS, err = NewFileServer(config.CacheFolder.Value(), 0, 0, false)
	if err != nil {
		return err
	}

	BackupFS, err = NewFileServer(config.BackupsFolder.Value(), 0, 0, false)
	if err != nil {
		return err
	}

	return err
}
