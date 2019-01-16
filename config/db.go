package config

type dbConfig struct {
	Config
}

func (cfg *dbConfig) Get(key string) (string, bool) {
	return "", false
}

