package pufferpanel

import (
	"github.com/spf13/viper"
	"strings"
)

func init() {
	viper.SetEnvPrefix("PUFFERPANEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/pufferpanel/")
	viper.AddConfigPath(".")

	//global settings
	viper.SetDefault("logs", "logs")

	//panel specific settings
	viper.SetDefault("panel.database.session", 60)
	viper.SetDefault("panel.database.dialect", "sqlite3")
	viper.SetDefault("panel.database.url", "file:pufferpanel.db?cache=shared")
	viper.SetDefault("panel.token.private", "private.pem")
	viper.SetDefault("panel.token.public", "public.pem")
	viper.SetDefault("panel.web.host", "0.0.0.0:8080")
	viper.SetDefault("panel.web.files", "www")
	viper.SetDefault("panel.email.templates", "email/emails.json")
	viper.SetDefault("panel.email.provider", "")
	viper.SetDefault("panel.settings.companyName", "PufferPanel")
	viper.SetDefault("panel.settings.masterUrl", "http://localhost:8080")
	viper.SetDefault("panel.localNode", true)

	//daemon specific settings
	viper.SetDefault("daemon.console.buffer", 50)
	viper.SetDefault("daemon.console.forward", false)
	viper.SetDefault("daemon.web.host", "0.0.0.0:5656")
	viper.SetDefault("daemon.sftp.host", "0.0.0.0:5657")
	viper.SetDefault("daemon.sftp.key", "sftp.key")
	viper.SetDefault("daemon.auth.publicKey", "panel.pem")
	viper.SetDefault("daemon.auth.url", "http://localhost:8080")
	viper.SetDefault("daemon.auth.clientId", "")
	viper.SetDefault("daemon.auth.clientSecret", "")
	viper.SetDefault("daemon.data.cache", "cache")
	viper.SetDefault("daemon.data.servers", "servers")
	viper.SetDefault("daemon.data.modules", "modules")
	viper.SetDefault("daemon.data.crashLimit", 3)
}

func LoadConfig(path string) error {
	if path != "" {
		viper.SetConfigFile(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

		} else {
			return err
		}
	}
	return nil
}
