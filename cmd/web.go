package cmd

import (
	"fmt"
	"os"

	"github.com/sidx1/sure/internal/config"
	"github.com/sidx1/sure/internal/web"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Launch the web dashboard to see mistakes caught and typing improvement",
	Long:  "Launch an interactive local web dashboard at http://localhost:7873 to view caught mistakes, command accuracy, typing improvement trends, and get an SVG badge to showcase on GitHub.",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		if port <= 0 {
			cfg, err := config.Load()
			if err == nil && cfg.WebPort > 0 {
				port = cfg.WebPort
			} else {
				port = 7873
			}
		}

		server := web.NewServer(port)
		if err := server.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Web server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	webCmd.Flags().IntP("port", "p", 7873, "Port to run the dashboard server on")
	rootCmd.AddCommand(webCmd)
}
