package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mingeme/imgood/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "imgood",
	Short: "imgood - Image processing and S3 management tool",
	Long: `imgood is a command-line tool for processing images and managing them in S3.
It supports uploading, copying, and listing images with various processing options.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	// Initialize configuration
	if err := config.Init(); err != nil {
		fmt.Printf("Warning: %s\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed

	// Add completion command
	rootCmd.AddCommand(completionCmd)
}
