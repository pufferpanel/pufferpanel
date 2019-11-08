package pufferpanel

import (
	"github.com/pufferpanel/pufferd/v2"
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
	viper.SetDefault("database.url", "file:pufferpanel.db?cache=shared")

	viper.SetDefault("token.private", "private.pem")
	viper.SetDefault("token.public", "public.pem")

	viper.SetDefault("web.host", "0.0.0.0:8080")
	//viper.SetDefault("web.socket", "/var/run/pufferpanel.sock")
	viper.SetDefault("web.files", "www")

	viper.SetDefault("email.templates", "email/emails.json")
	viper.SetDefault("email.provider", "")

	viper.SetDefault("settings.companyName", "PufferPanel")
	viper.SetDefault("settings.masterUrl", "http://localhost:8080")

	viper.SetDefault("logs", "logs")
	viper.SetDefault("localNode", true)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	pufferd.SetDefaults()
}

func LoadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//this is just a missing config, since ENV is supported, ignore
		} else {
			return err
		}
	}
	return nil
}
