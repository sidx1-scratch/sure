package cmd

import (
	"fmt"

	"github.com/sidx1/sure/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.Load()
		out, _ := yaml.Marshal(cfg)
		fmt.Println(string(out))
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset to default configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.DefaultConfig()
		_ = config.Save(cfg)
		fmt.Println("Configuration reset to defaults.")
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}
