package cmd

import (
	"fmt"

	"github.com/sidx1/sure/internal/discovery"
	"github.com/sidx1/sure/internal/tui"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for CLI utilities",
	Run: func(cmd *cobra.Command, args []string) {
		tui.PrintInfo("Scanning for commands in PATH...")
		cmds, err := discovery.DiscoverCommands()
		if err != nil {
			tui.PrintError(fmt.Sprintf("Scan failed: %v", err))
			return
		}
		tui.PrintSuccess(fmt.Sprintf("Found %d commands.", len(cmds)))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
