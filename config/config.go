package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	viper.SetEnvPrefix("PUFFER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func LoadConfigFile(configFile string) error {
	if configFile == "" {
		var exists bool
		configFile, exists = os.LookupEnv("PUFFER_CONFIG")
		if !exists || configFile == "" {
			if runtime.GOOS == "windows" {
				//well, check for the program files path....
				if _, err := os.Lstat("C:\\Program Files\\PufferPanel\\config.json"); err == nil {
					configFile = "C:\\Program Files\\PufferPanel\\config.json"
				}
			} else {
				//well, check for the /etc path
				if _, err := os.Lstat("/etc/pufferpanel/config.json"); err == nil {
					configFile = "/etc/pufferpanel/config.json"
				}
			}
			//we got nothing
			if configFile == "" {
				configFile = "config.json"
			}
		}
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	if DataRootFolder.Value() == "" {
		//we need to set root and save it
		serversDir := ServersFolder.Value()
		err := DataRootFolder.Set(filepath.Dir(serversDir), true)
		if err != nil {
			return err
		}
	}

	return nil
}

var configFileNotFoundError viper.ConfigFileNotFoundError
