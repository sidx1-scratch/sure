package cmd

import (
	"fmt"
	"os"

	"github.com/sidx1/sure/internal/cache"
	"github.com/sidx1/sure/internal/config"
	"github.com/sidx1/sure/internal/keyring"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current status",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Config Error: %v\n", err)
			return
		}
		
		fmt.Printf("Enabled: %v\n", cfg.Enabled)
		fmt.Printf("Provider: %s\n", cfg.Provider)
		
		_, err = keyring.GetAPIKey(cfg.Provider)
		if err == nil {
			fmt.Printf("API Key: Set\n")
		} else {
			fmt.Printf("API Key: Not set (%v)\n", err)
		}
		
		size, _ := cache.Size()
		fmt.Printf("Cached Commands: %d\n", size)
		fmt.Printf("Config File: %s\n", config.ConfigDir()+"/config.yaml")
		
		if os.Getenv("SURE_SKIP") != "" {
			fmt.Printf("Shell Integration: Active\n")
		} else {
			fmt.Printf("Shell Integration: Not detected (or running standalone)\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
