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
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

var dbConn *gorm.DB
var lock sync.Mutex

func openConnection() (err error) {
	//lock system so we can only connect one at a time
	lock.Lock()
	defer lock.Unlock()

	//if we had 2 calls to this before it was established, quick out since it's already created
	if dbConn != nil {
		return
	}

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
	} else if dialect == "sqlite3" {
		connString = addConnectionSetting(connString, "_loc=auto")
	}

	//attempt to open database connection to validate
	dbConn, err = gorm.Open(dialect, connString)

	if err != nil {
		dbConn = nil
		logging.Error().Printf("Error connecting to database: %s", err)
		return pufferpanel.ErrDatabaseNotAvailable
	}

	if viper.GetBool("database.log") {
		logging.Info().Printf("Database logging enabled")
		dbConn.LogMode(true)
	}

	err = migrateModels()
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
	pufferpanel.Close(dbConn)
}

func migrateModels() (err error) {
	dbObjects := []interface{}{
		&models.Node{},
		&models.Server{},
		&models.User{},
		&models.Template{},
		&models.Permissions{},
		&models.Client{},
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

	err = dbConn.Model(&models.Permissions{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return
	}

	err = dbConn.Model(&models.Permissions{}).AddForeignKey("server_identifier", "servers(identifier)", "CASCADE", "CASCADE").Error
	if err != nil {
		return
	}

	err = dbConn.Model(&models.Permissions{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE").Error
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
	connString += setting

	return connString
}
