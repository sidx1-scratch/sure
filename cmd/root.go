package cmd

import (
	"fmt"
	"os"

	"github.com/sidx1/sure/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sure",
	Short: "Your terminal's second opinion",
	Long: `sure is a second opinion layer for terminal commands.

It analyzes commands before execution and warns you about potential
issues like typos, destructive operations, or incorrect syntax.

sure NEVER modifies your commands. It only asks if you want to proceed.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version.Version
	rootCmd.SetVersionTemplate(fmt.Sprintf("sure %s (commit: %s, built: %s)\n",
		version.Version, version.Commit, version.Date))
}
