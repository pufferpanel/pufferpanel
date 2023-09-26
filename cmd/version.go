package main

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of PufferPanel",
	Run:   executeVersion,
}

func executeVersion(cmd *cobra.Command, args []string) {
	fmt.Println(pufferpanel.Display)
}
