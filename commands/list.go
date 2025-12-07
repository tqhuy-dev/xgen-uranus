package commands

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Long: `List available resources in the project.

Available subcommands:
  repo    List all repositories

Examples:
  uranus list repo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	listCmd.AddCommand(listRepoCmd)
}


