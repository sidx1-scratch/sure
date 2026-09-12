package cmd

import (
	"fmt"

	"github.com/sidx1/sure/internal/cache"
	"github.com/sidx1/sure/internal/config"
	"github.com/sidx1/sure/internal/keyring"
	"github.com/sidx1/sure/internal/tui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose issues",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running sure doctor...")
		
		cfg, err := config.Load()
		if err != nil {
			tui.PrintError("Config: Failed to load")
		} else {
			tui.PrintSuccess("Config: Loaded ok")
		}
		
		_, err = keyring.GetAPIKey(cfg.Provider)
		if err != nil {
			tui.PrintError(fmt.Sprintf("API Key: Missing for %s", cfg.Provider))
		} else {
			tui.PrintSuccess(fmt.Sprintf("API Key: Found for %s", cfg.Provider))
		}
		
		_, err = cache.Size()
		if err != nil {
			tui.PrintError("Cache Dir: Unreachable")
		} else {
			tui.PrintSuccess("Cache Dir: Writable")
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
