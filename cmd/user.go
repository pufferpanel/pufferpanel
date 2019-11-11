package main

import (
	"github.com/pufferpanel/pufferpanel/v2/cmd/user"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

func init() {
	userCmd.AddCommand(
		user.AddUserCmd,
		user.EditUserCmd)
}
