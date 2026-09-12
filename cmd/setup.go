package cmd

import (
	"fmt"

	"github.com/sidx1/sure/internal/config"
	"github.com/sidx1/sure/internal/keyring"
	"github.com/sidx1/sure/internal/tui"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run the setup wizard",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := tui.SetupWizard()
		if err != nil {
			tui.PrintError(fmt.Sprintf("Setup failed: %v", err))
			return
		}

		if key != "" {
			err = keyring.StoreAPIKey("gemini", key)
			if err != nil {
				tui.PrintError(fmt.Sprintf("Failed to store API key: %v", err))
				return
			}
			tui.PrintSuccess("API key stored securely.")
		}

		cfg := config.DefaultConfig()
		if err := config.Save(cfg); err != nil {
			tui.PrintError(fmt.Sprintf("Failed to save default config: %v", err))
		} else {
			tui.PrintSuccess("Default configuration created.")
		}

		fmt.Println("\nTo enable sure in your shell, add this to your ~/.bashrc or ~/.zshrc:")
		fmt.Println(`eval "$(sure shell-init --shell bash)"  # or --shell zsh`)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
