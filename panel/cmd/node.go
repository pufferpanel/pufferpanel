package main

import (
	"github.com/pufferpanel/pufferpanel/v2/panel/cmd/node"
	"github.com/spf13/cobra"
)

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage nodes",
}

func init() {
	nodesCmd.AddCommand(
		node.AddCmd,
	)
}
