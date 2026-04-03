package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pack-compose",
	Short: "A tool to bundle Docker Compose services for offline use",
	Long:  `pack-compose parses Docker Compose files, pulls images, and creates offline bundles.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")
}
