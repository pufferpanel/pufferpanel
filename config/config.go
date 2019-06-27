package config

import (
	"github.com/spf13/viper"
	"strings"
)

func init() {
	viper.SetEnvPrefix("PUFFERPANEL")
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/pufferpanel/")
	viper.AddConfigPath(".")

	viper.SetDefault("database.session", 60)
	viper.SetDefault("database.dialect", "sqlite3")
	viper.SetDefault("database.url", "file:pufferpanel.db?cache=shared&mode=memory")

	viper.SetDefault("web.host", "0.0.0.0")
	viper.SetDefault("web.port", "8080")
	viper.SetDefault("web.socket", "/var/run/pufferpanel.sock")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func Load() error {
	return viper.ReadInConfig()
}