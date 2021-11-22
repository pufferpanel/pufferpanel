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

package config

import (
	"encoding/hex"
	"errors"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func init() {
	viper.SetEnvPrefix("PUFFER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/pufferpanel/")
	viper.AddConfigPath(".")

	for k, v := range defaultSettings {
		viper.SetDefault(k, v)
	}
}

type Setting struct {
	Key   string `gorm:"type:varchar(100);primaryKey"`
	Value string `gorm:"type:varchar(100)"`
}

type db interface {
	GetConnection() (*gorm.DB, error)
}

var database db
var defaultSettings = map[string]interface{}{
	//global settings
	"logs":          "logs",
	"web.host":      "0.0.0.0:8080",
	"token.private": "private.pem",
	"token.public":  "public.pem",

	//panel specific settings
	"panel.enable":                true,
	"panel.database.session":      60,
	"panel.database.dialect":      "sqlite3",
	"panel.database.log":          false,
	"panel.web.files":             "www",
	"panel.email.templates":       "email/emails.json",
	"panel.email.provider":        "",
	"panel.email.from":            "",
	"panel.email.domain":          "",
	"panel.email.key":             "",
	"panel.settings.companyName":  "PufferPanel",
	"panel.settings.defaultTheme": "PufferPanel",
	"panel.settings.masterUrl":    "http://localhost:8080",
	"panel.sessionKey":            []uint8{},
	"panel.registrationEnabled":   true,

	//daemon specific settings
	"daemon.enable":                 true,
	"daemon.console.buffer":         50,
	"daemon.console.forward":        false,
	"daemon.sftp.host":              "0.0.0.0:5657",
	"daemon.sftp.key":               "sftp.key",
	"daemon.auth.url":               "http://localhost:8080",
	"daemon.auth.clientId":          "",
	"daemon.auth.clientSecret":      "",
	"daemon.data.cache":             "cache",
	"daemon.data.servers":           "servers",
	"daemon.data.modules":           "modules",
	"daemon.data.crashLimit":        3,
	"daemon.data.maxWSDownloadSize": int64(1024 * 1024 * 20),
}

func LoadConfigFile(path string) error {
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

func fromDatabase(key string) (value *string, err error) {
	db, err := database.GetConnection()
	if err != nil {
		return
	}

	configEntry := &Setting{
		Key: key,
	}

	err = db.First(configEntry).Error
	if err != nil && gorm.ErrRecordNotFound != err {
		return nil, err
	}

	if err == nil {
		value = &configEntry.Value
	}

	return value, nil
}

func LoadConfigDatabase(db db) error {
	database = db
	for k, v := range defaultSettings {
		fromDB, err := fromDatabase(k)
		if err != nil {
			return err
		}

		if fromDB == nil {
			continue
		}

		switch v.(type) {
		case string:
			viper.Set(k, fromDB)
		case bool:
			viper.Set(k, cast.ToBool(fromDB))
		case int64:
			viper.Set(k, cast.ToInt64(fromDB))
		case int:
			viper.Set(k, cast.ToInt(fromDB))
		default:
			viper.Set(k, fromDB)
		}
	}
	return nil
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func toDatabase(key string, value string) error {
	setting := &Setting{
		Key:   key,
		Value: value,
	}

	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	fromDB, err := fromDatabase(key)
	if err != nil {
		return err
	}

	if fromDB == nil {
		return db.Create(setting).Error
	} else {
		return db.Save(setting).Error
	}
}

func Set(key string, value interface{}) (err error) {
	viper.Set(key, value)

	if database != nil {
		if v, ok := value.(string); ok {
			err = toDatabase(key, v)
		} else if v, ok := value.(bool); ok {
			err = toDatabase(key, cast.ToString(v))
		} else if v, ok := value.(int); ok {
			err = toDatabase(key, cast.ToString(v))
		} else if v, ok := value.(int64); ok {
			err = toDatabase(key, cast.ToString(v))
		} else if v, ok := value.([]uint8); ok {
			err = toDatabase(key, hex.EncodeToString(v))
			//we have to override it here
			viper.Set(key, hex.EncodeToString(v))
		} else {
			err = errors.New("Unsupported type for " + key + ": " + reflect.TypeOf(value).String())
		}

		if err != nil {
			return
		}
	}

	return
}
