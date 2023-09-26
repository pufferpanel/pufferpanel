package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
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
			configFile = "config.json"
		}
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	return nil
}

var configFileNotFoundError viper.ConfigFileNotFoundError
