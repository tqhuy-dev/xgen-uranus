package commands

import (
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate new resources (app, repo)",
	Long: `Generate new resources for your microservice project.

Available subcommands:
  app     Generate a new application
  repo    Generate a new repository

Examples:
  uranus generate app --name my_app --module github.com/my_org/my_app
  uranus generate repo --name my_repo
  uranus generate repo --name path/to/my_repo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	generateCmd.AddCommand(generateAppCmd)
	generateCmd.AddCommand(generateRepoCmd)
}
