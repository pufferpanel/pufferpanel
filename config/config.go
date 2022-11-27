/*
 Copyright 2022 Padduck, LLC
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
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	viper.SetEnvPrefix("PUFFER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func LoadConfigFile(workDir string) error {
	if workDir != "" {
		viper.SetConfigFile(filepath.Join(workDir, "config.json"))
	} else if runtime.GOOS != "windows" {
		viper.SetConfigFile(filepath.Join("/", "etc", "pufferpanel", "config.json"))
	} else {
		viper.SetConfigFile("config.json")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			return err
		}
	}

	return nil
}
