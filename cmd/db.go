package main

import "github.com/spf13/cobra"

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Runs DB operations",
}

func init() {
	dbCmd.AddCommand(dbUpgradeCmd, dbMigrateCmd)
}
