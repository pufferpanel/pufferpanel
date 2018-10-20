package config

import (
	"encoding/json"
	"io"
)

type Config struct {
	Database Database `json:"database"`
}

type Database struct {
	Url     string `json:"url"`
	Dialect string `json:"dialect"`
}

var config Config

func Load(reader io.Reader) error {
	config = Config{}
	return json.NewDecoder(reader).Decode(&config)
}

func Get() Config {
	return config
}
