package main

import (
	"github.com/spf13/cobra"
)

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates the database to a new backend",
	//Run:   executeDbMigrateToNewSystem,
}

/*
func executeDbMigrateToNewSystem(cmd *cobra.Command, args []string) {
	survey.AskOne(&survey.OptionAnswer{
		Message: ""
	})
}
*/
