package config

type Config interface {
	GetString(key string) (val string, exist bool)
}

func Get() Config {
	return &dbConfig{}
}