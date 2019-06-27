/*
 Copyright 2018 Padduck, LLC
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

package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pufferpanel/apufferi"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/spf13/viper"
	"strings"
)

var dbConn *gorm.DB

func Load() error {
	err := openConnection()
	if err != nil {
		return err
	}

	err = migrateModels()

	return err
}

func openConnection() (err error) {
	dialect := viper.GetString("database.dialect")
	if dialect == "" {
		dialect = "sqlite3"
	}
	connString := viper.GetString("database.url")
	if connString == "" {
		switch dialect {
		case "mysql":
			connString = "pufferpanel:pufferpanel@/pufferpanel"
		case "sqlite3":
			connString = "file:pufferpanel.db?cache=shared&mode=memory"
		}
	}

	if dialect == "mysql" {
		connString = addConnectionSetting(connString, "charset=utf8")
		connString = addConnectionSetting(connString, "parseTime=true")
	} else if dialect == "sqlite3"{
		connString = addConnectionSetting(connString, "_loc=auto")
	}

	//attempt to open database connection to validate
	dbConn, err = gorm.Open(dialect, connString)

	if err != nil {
		return
	}

	if viper.GetBool("database.log") {
		logging.Info("Database logging enabled")
		dbConn.LogMode(true)
	}

	return
}

func GetConnection() (*gorm.DB, error) {
	var err error
	if dbConn == nil {
		err = openConnection()
	}

	return dbConn, err
}

func Close() {
	apufferi.Close(dbConn)
}

func migrateModels() (err error) {
	dbObjects := []interface{}{
		&models.Node{},
		&models.Server{},
		&models.User{},
		&models.ClientInfo{},
		&models.ClientServerScopes{},
		&models.TokenInfo{},
	}

	for _, v := range dbObjects {
		dbConn.AutoMigrate(v)
	}

	dialect := viper.GetString("database.dialect")
	if dialect == "" || dialect == "sqlite3" {
		return
	}

	err = dbConn.Model(&models.Server{}).AddForeignKey("node_id", "nodes(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return
	}

	err = dbConn.Model(&models.ClientInfo{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return
	}

	err = dbConn.Model(&models.ClientServerScopes{}).AddForeignKey("server_id", "servers(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return
	}

	err = dbConn.Model(&models.ClientServerScopes{}).AddForeignKey("client_info_id", "client_infos(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return
	}

	err = dbConn.Model(&models.TokenInfo{}).AddForeignKey("client_info_id", "client_infos(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return
	}

	return
}

func addConnectionSetting(connString, setting string) string {
	if strings.Contains(connString, setting) {
		return connString
	}

	if !strings.Contains(connString, "?") {
		connString += "?"
	} else {
		connString += "&"
	}
	connString += "charset=utf8"

	return connString
}