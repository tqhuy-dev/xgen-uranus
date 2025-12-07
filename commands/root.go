package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "uranus",
	Short: "Uranus is a CLI tool to generate microservice components",
	Long: `Uranus is a powerful CLI tool designed to help developers 
quickly scaffold and manage microservice applications and repositories.

Available Commands:
  generate    Generate new app or repository
  list        List existing resources

Use "uranus [command] --help" for more information about a command.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(listCmd)
}
