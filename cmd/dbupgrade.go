package main

import (
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/pterm/pterm"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

var dbUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Runs the database upgrades",
	Run:   executeDbUpgrade,
}

func executeDbUpgrade(cmd *cobra.Command, args []string) {
	var currentFile string
	var backupFile string

	if !config.PanelEnabled.Value() {
		pterm.Info.Printfln("Panel not enabled, skipping upgrade")
		os.Exit(0)
	}

	if database.GetDialect() == "sqlite3" {
		//we could get the filename... let's get it
		drv := sqlite3.SQLiteDriver{}
		conn, err := drv.Open(database.GetConnectionString())
		if err != nil {
			pterm.Error.Printfln("error connecting to database: %s", err.Error())
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
			pterm.Error.Printfln("error backing up database: %s", err.Error())
			os.Exit(1)
			return
		}
	}

	db, err := database.GetConnection()
	if err != nil {
		pterm.Error.Printfln("error connecting to database: %s", err.Error())
		rollback(backupFile, currentFile)
		os.Exit(1)
		return
	}

	pterm.Info.Printfln("Starting database upgrade")
	err = database.Upgrade(db, true)
	if err != nil {
		pterm.Error.Printfln("error upgrading database: %s", err.Error())
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
		pterm.Error.Printfln("error restoring database: %s", err.Error())
		return
	}
}

func copyDatabaseFile(src, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer utils.Close(source)

	err = os.MkdirAll(filepath.Dir(dest), 0750)
	if err != nil {
		return err
	}
	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer utils.Close(destination)
	_, err = io.Copy(destination, source)
	return err
}
