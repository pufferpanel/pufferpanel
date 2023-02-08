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

package database

import (
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"sync"
	"time"
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

	dialect := config.DatabaseDialect.Value()
	if dialect == "" {
		dialect = "sqlite3"
	}
	connString := config.DatabaseUrl.Value()
	if connString == "" {
		switch dialect {
		case "mysql":
			connString = "pufferpanel:pufferpanel@/pufferpanel"
		case "sqlite3":
			connString = "file:pufferpanel.db?cache=shared"
		}
	}

	if dialect == "mysql" {
		connString = addConnectionSetting(connString, "charset=utf8")
		connString = addConnectionSetting(connString, "parseTime=true")
	} else if dialect == "sqlite3" {
		connString = addConnectionSetting(connString, "_loc=auto")
		connString = addConnectionSetting(connString, "_foreign_keys=1")
	}

	var dialector gorm.Dialector
	switch dialect {
	case "sqlite3":
		dialector = sqlite.Open(connString)
	case "mysql":
		dialector = mysql.Open(connString)
	case "postgresql":
		dialector = postgres.Open(connString)
	case "sqlserver":
		dialector = sqlserver.Open(connString)
	default:
		return errors.New(fmt.Sprintf("unknown dialect %s", dialect))
	}

	gormConfig := gorm.Config{}
	gormConfig.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})

	if config.DatabaseLoggingEnabled.Value() {
		logging.Info.Printf("Database logging enabled")
		gormConfig.Logger = gormConfig.Logger.LogMode(logger.Info)
	}

	// Sqlite doesn't implement constraints see  https://github.com/go-gorm/gorm/wiki/GORM-V2-Release-Note-Draft#all-new-migratolease-Note-Draft#all-new-migrator
	gormConfig.DisableForeignKeyConstraintWhenMigrating = dialect == "sqlite3"

	dbConn, err = gorm.Open(dialector, &gormConfig)

	if err != nil {
		dbConn = nil
		logging.Error.Printf("Error connecting to database: %s", err)
		return pufferpanel.ErrDatabaseNotAvailable
	}
	if err := migrateModels(); err != nil {
		return err
	}
	return migrate(dbConn)
}

func GetConnection() (*gorm.DB, error) {
	var err error
	if dbConn == nil {
		err = openConnection()
	}

	return dbConn, err
}

func Close() {
	if dbConn != nil {
		sqlDB, _ := dbConn.DB()
		pufferpanel.Close(sqlDB)
	}
}

func migrateModels() error {
	dbObjects := []interface{}{
		&models.Node{},
		&models.Server{},
		&models.User{},
		&models.Template{},
		&models.Permissions{},
		&models.Client{},
		&models.UserSetting{},
		&models.Session{},
		&models.TemplateRepo{},
	}

	for _, v := range dbObjects {
		if err := dbConn.AutoMigrate(v); err != nil {
			return err
		}
	}
	return migrate(dbConn)
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
