package main

import (
	"github.com/pufferpanel/pufferpanel/v2/cmd/templates"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage templates",
	Run:   executeRun,
}

func init() {
	templatesCmd.AddCommand(
		templates.AddCmd,
		templates.ImportCmd)
}
