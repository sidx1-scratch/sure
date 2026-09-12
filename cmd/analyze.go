package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sidx1/sure/internal/analyzer"
	"github.com/sidx1/sure/internal/warning"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze a command",
	Run: func(cmd *cobra.Command, args []string) {
		command, _ := cmd.Flags().GetString("command")
		jsonOut, _ := cmd.Flags().GetBool("json")

		if command == "" {
			if jsonOut {
				fmt.Println(`{"error": "no command provided"}`)
			} else {
				fmt.Println("Error: no command provided")
			}
			os.Exit(1)
		}

		result, err := analyzer.Analyze(command)
		if err != nil {
			if jsonOut {
				fmt.Printf(`{"error": "%s"}\n`, err.Error())
			} else {
				fmt.Fprintf(os.Stderr, "Analysis error: %v\n", err)
			}
			os.Exit(0)
		}

		if jsonOut {
			data, _ := json.Marshal(result)
			fmt.Println(string(data))
			if result.Warn {
				os.Exit(1)
			}
			os.Exit(0)
		}

		proceed := warning.Display(result, command)
		if !proceed {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	analyzeCmd.Flags().String("command", "", "Command to analyze")
	analyzeCmd.Flags().Bool("json", false, "Output JSON result")
	rootCmd.AddCommand(analyzeCmd)
}
