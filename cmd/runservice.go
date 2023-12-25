package main

import "github.com/spf13/cobra"

var runServiceCmd = &cobra.Command{
	Use:   "runService",
	Short: "Runs the panel as a service",
	Run:   executeRunService,
}
