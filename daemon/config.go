package daemon

import (
	"github.com/spf13/viper"
	"strings"
)

func init() {
	//env configuration
	viper.SetEnvPrefix("PUFFERPANEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func SetDefaults() {
	//defaults we can set at this point in time
	viper.SetDefault("console.buffer", 50)
	viper.SetDefault("console.forward", false)

	viper.SetDefault("listen.web", "0.0.0.0:5656")
	//viper.SetDefault("listen.socket", "unix:/var/run/daemon.sock")
	viper.SetDefault("listen.webCert", "https.pem")
	viper.SetDefault("listen.webKey", "https.key")
	viper.SetDefault("listen.sftp", "0.0.0.0:5657")
	viper.SetDefault("listen.sftpKey", "sftp.key")

	viper.SetDefault("auth.publicKey", "panel.pem")

	viper.SetDefault("auth.url", "http://localhost:8080")

	viper.SetDefault("auth.clientId", "")
	viper.SetDefault("auth.clientSecret", "")

	viper.SetDefault("data.cache", "cache")
	viper.SetDefault("data.servers", "servers")
	viper.SetDefault("data.modules", "modules")
	viper.SetDefault("data.logs", "logs")
	viper.SetDefault("data.crashLimit", 3)
	viper.SetDefault("data.maxWSDownloadSize", int64(1024*1024*20)) //1024 bytes (1KB) * 1024 (1MB) * 50 (50MB))
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
