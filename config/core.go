package config

import (
	"encoding/json"
	"github.com/pufferpanel/apufferi/common"
	"github.com/pufferpanel/apufferi/logging"
	"io"
	"os"
	"sync"
)

type CoreConfig struct {
	Database Database `json:"database"`
	Session  Session  `json:"session"`
}

type Database struct {
	Url     string `json:"url"`
	Dialect string `json:"dialect"`
}

type Session struct {
	Timeout int `json:"timeout"`
}

var config *CoreConfig
var loaded bool
var locker sync.Locker = &sync.Mutex{}

func LoadDefault() error {
	locker.Lock()
	defer locker.Unlock()
	configPath, exists := os.LookupEnv("PUFFERPANEL_CONFIG")

	if !exists {
		configPath = "config.json"
	}

	cfg, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(configPath)
			defer common.Close(file)
			if err != nil {
				logging.Error("Error creating config: %s", err.Error())
				return err
			}
			data := getDefault()
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(true)
			err = encoder.Encode(&data)
			if err != nil {
				logging.Error("Error creating config: %s", err.Error())
				return err
			}

		} else {
			if err != nil {
				logging.Error("Error reading config: %s", err.Error())
				return err
			}
			return err
		}
	}

	err = Load(cfg)
	if err != nil {
		logging.Error("Error reading config: %s", err.Error())
		return err
	}
	loaded = true
	return cfg.Close()
}

func Load(reader io.Reader) error {
	config = &CoreConfig{}
	return json.NewDecoder(reader).Decode(config)
}

func GetCore() (*CoreConfig, error) {
	var err error = nil
	if !loaded {
		err = LoadDefault()
	}
	return config, err
}

func getDefault() *CoreConfig {
	return &CoreConfig{
		Database: Database{},
		Session: Session{
			Timeout: 60,
		},
	}
}
