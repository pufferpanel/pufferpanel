package main

import (
	"github.com/pufferpanel/pufferpanel/v2/panel/cmd/template"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage templates",
}

func init() {
	templatesCmd.AddCommand(
		template.AddCmd,
		template.ImportCmd)
}
