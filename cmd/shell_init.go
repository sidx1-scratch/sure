package cmd

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed shell_scripts
var shellScripts embed.FS

var shellInitCmd = &cobra.Command{
	Use:   "shell-init",
	Short: "Output shell integration script",
	Long:  "Output shell integration script for bash or zsh. Add to your shell config:\n  Bash: eval \"$(sure shell-init --shell bash)\"\n  Zsh:  eval \"$(sure shell-init --shell zsh)\"",
	Run: func(cmd *cobra.Command, args []string) {
		shell, _ := cmd.Flags().GetString("shell")

		switch strings.ToLower(shell) {
		case "bash":
			data, err := shellScripts.ReadFile("shell_scripts/sure.bash")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading bash integration script: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(data))
		case "zsh":
			data, err := shellScripts.ReadFile("shell_scripts/sure.zsh")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading zsh integration script: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(data))
		default:
			fmt.Fprintf(os.Stderr, "Unsupported shell: %s (supported: bash, zsh)\n", shell)
			os.Exit(1)
		}
	},
}

func init() {
	shellInitCmd.Flags().String("shell", "bash", "Shell type (bash or zsh)")
	rootCmd.AddCommand(shellInitCmd)
}
