package cmd

import (
	"fmt"
	"strings"

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

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value (e.g. provider, model, sensitivity, local.endpoint, local.model)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := strings.ToLower(args[0])
		val := args[1]

		cfg, err := config.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}

		switch key {
		case "provider":
			cfg.Provider = val
		case "model":
			cfg.Model = val
		case "sensitivity":
			cfg.Sensitivity = val
		case "enabled":
			cfg.Enabled = (val == "true" || val == "1")
		case "local.endpoint", "endpoint":
			cfg.Local.Endpoint = val
		case "local.model":
			cfg.Local.Model = val
		case "web_port", "port":
			var p int
			fmt.Sscanf(val, "%d", &p)
			if p > 0 {
				cfg.WebPort = p
			}
		default:
			fmt.Printf("Unknown config key %q. Valid keys: provider, model, sensitivity, enabled, local.endpoint, local.model, web_port\n", key)
			return
		}

		if err := config.Save(cfg); err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return
		}
		fmt.Printf("Configuration updated: %s = %s\n", key, val)
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
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}
