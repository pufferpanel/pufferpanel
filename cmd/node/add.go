package node

import (
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use: "add",
	Short: "Adda a node",
	Run: runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
}
