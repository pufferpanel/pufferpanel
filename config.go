/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package pufferpanel

import (
	"github.com/spf13/viper"
	"strings"
)

func init() {
	viper.SetEnvPrefix("PUFFER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/pufferpanel/")
	viper.AddConfigPath(".")

	//global settings
	viper.SetDefault("logs", "logs")
	viper.SetDefault("web.host", "0.0.0.0:8080")

	//panel specific settings
	viper.SetDefault("panel.enable", true)
	viper.SetDefault("panel.database.session", 60)
	viper.SetDefault("panel.database.dialect", "sqlite3")
	//viper.SetDefault("panel.database.url", "file:pufferpanel.db?cache=shared")
	viper.SetDefault("panel.database.log", false)
	viper.SetDefault("panel.token.private", "private.pem")
	viper.SetDefault("panel.token.public", "public.pem")

	viper.SetDefault("panel.web.files", "www")
	viper.SetDefault("panel.email.templates", "email/emails.json")
	viper.SetDefault("panel.email.provider", "")
	viper.SetDefault("panel.email.from", "")
	viper.SetDefault("panel.email.domain", "")
	viper.SetDefault("panel.email.key", "")
	viper.SetDefault("panel.settings.companyName", "PufferPanel")
	viper.SetDefault("panel.settings.masterUrl", "http://localhost:8080")

	//daemon specific settings
	viper.SetDefault("daemon.enable", true)
	viper.SetDefault("daemon.console.buffer", 50)
	viper.SetDefault("daemon.console.forward", false)
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
	viper.SetDefault("daemon.data.maxWSDownloadSize", int64(1024 * 1024 * 20))
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
