package main

import (
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

var migrateCmd = &cobra.Command{
	Use:    "dbmigrate",
	Short:  "Runs the database migrations",
	Run:    executeDbMigrations,
	Hidden: true,
}

func executeDbMigrations(cmd *cobra.Command, args []string) {
	var currentFile string
	var backupFile string

	if database.GetDialect() == "sqlite3" {
		//we could get the filename... let's get it
		drv := sqlite3.SQLiteDriver{}
		conn, err := drv.Open(database.GetConnectionString())
		if err != nil {
			logging.Error.Printf("error connecting to database: %s", err.Error())
			os.Exit(1)
			return
		}
		s3 := conn.(*sqlite3.SQLiteConn)
		currentFile = s3.GetFilename("")
		_ = conn.Close()

		//look for a new name we can give this....
		suffix := "backup"
		num := 0
		for {
			backupFile = fmt.Sprintf("%s.%d.%s", currentFile, num, suffix)
			fi, err := os.Lstat(backupFile)
			if os.IsNotExist(err) && fi == nil {
				break
			}
			num++
		}

		err = copyDatabaseFile(currentFile, backupFile)
		if err != nil {
			logging.Error.Printf("error backing up database: %s", err.Error())
			os.Exit(1)
			return
		}
	}

	db, err := database.GetConnection()
	if err != nil {
		logging.Error.Printf("error connecting to database: %s", err.Error())
		rollback(backupFile, currentFile)
		os.Exit(1)
		return
	}

	logging.Info.Printf("Starting database migration")
	err = database.Migrate(db)
	if err != nil {
		logging.Error.Printf("error upgrading database: %s", err.Error())
		rollback(backupFile, currentFile)
		os.Exit(1)
		return
	}
}

func rollback(backup, overrideTo string) {
	if backup == "" || overrideTo == "" {
		return
	}
	err := copyDatabaseFile(backup, overrideTo)
	if err != nil {
		logging.Error.Printf("error restoring database: %s", err.Error())
		return
	}
}

func copyDatabaseFile(src, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer pufferpanel.Close(source)

	err = os.MkdirAll(filepath.Dir(dest), 0750)
	if err != nil {
		return err
	}
	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer pufferpanel.Close(destination)
	_, err = io.Copy(destination, source)
	return err
}
