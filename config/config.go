package config

import (
	"encoding/json"
	"github.com/pufferpanel/apufferi/logging"
	"io"
	"os"
	"sync"
)

type Config struct {
	Database Database `json:"database"`
}

type Database struct {
	Url     string `json:"url"`
	Dialect string `json:"dialect"`
}

var config Config
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
			if err != nil {
				logging.Error("Error creating config", err)
				return err
			}
			data := Config{Database: Database{}}
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(true)
			err = encoder.Encode(&data)
			if err != nil {
				logging.Error("Error creating config", err)
				return err
			}
		} else {
			if err != nil {
				logging.Error("Error reading config", err)
				return err
			}
			return err
		}
	}

	err = Load(cfg)
	if err != nil {
		logging.Error("Error reading config", err)
		return err
	}
	loaded = true
	return cfg.Close()
}

func Load(reader io.Reader) error {
	config = Config{}
	return json.NewDecoder(reader).Decode(&config)
}

func Get() Config {
	if !loaded {
		LoadDefault()
	}
	return config
}
