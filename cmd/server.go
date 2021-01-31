package main

import (
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers",
	PreRunE: prepare,
}

var startServerCmd = &cobra.Command{
	Use: "start",
	Args: cobra.MinimumNArgs(1),
	Run: runServerStart,
}
var stopServerCmd = &cobra.Command{
	Use: "stop",
	Args: cobra.MinimumNArgs(1),
	Run: runServerStop,
}

var waitForStop bool

func init() {
	serverCmd.AddCommand(startServerCmd, stopServerCmd)

	stopServerCmd.Flags().BoolVarP(&waitForStop, "wait", "w", false, "Wait for servers to stop before exiting")
}

func prepare(cmd *cobra.Command, args []string) error {
	logging.DisableFileLogger()

	if err := pufferpanel.LoadConfig(""); err != nil {
		return err
	}

	return nil
}

func runServerStart(cmd *cobra.Command, args []string) {
}

func runServerStop(cmd *cobra.Command, args []string) {

}