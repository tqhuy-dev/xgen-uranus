package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var listRepoPath string

var listRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "List all repositories",
	Long: `List all repositories in the current directory or specified path.

This command scans for repository patterns (directories containing repository.go files)
and displays them in a formatted list.

Examples:
  uranus list repo
  uranus list repo --path ./internal`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listRepositories(listRepoPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing repos: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	listRepoCmd.Flags().StringVarP(&listRepoPath, "path", "p", ".", "Path to search for repositories")
}

func listRepositories(searchPath string) error {
	fmt.Printf("🔍 Searching for repositories in: %s\n\n", searchPath)

	repos := []string{}

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip directories we can't access
		}

		// Skip hidden directories and common non-repo directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// Look for repository.go files as indicators of a repository
		if info.Name() == "repository.go" {
			repoDir := filepath.Dir(path)
			repos = append(repos, repoDir)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(repos) == 0 {
		fmt.Println("📭 No repositories found.")
		fmt.Println("\nTip: Use 'uranus generate repo --name <name>' to create a new repository.")
		return nil
	}

	fmt.Printf("📦 Found %d repository(ies):\n\n", len(repos))
	fmt.Println(strings.Repeat("─", 60))

	for i, repo := range repos {
		repoName := filepath.Base(repo)
		
		// Try to get more info about the repo
		repoInfo := getRepoInfo(repo)
		
		fmt.Printf("  %d. %s\n", i+1, repoName)
		fmt.Printf("     📍 Path: %s\n", repo)
		if repoInfo != "" {
			fmt.Printf("     📋 %s\n", repoInfo)
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("\n✅ Total: %d repository(ies)\n", len(repos))

	return nil
}

func getRepoInfo(repoPath string) string {
	// Check what files exist in the repo directory
	files := []string{}
	
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			files = append(files, entry.Name())
		}
	}

	if len(files) == 0 {
		return ""
	}

	return fmt.Sprintf("Files: %s", strings.Join(files, ", "))
}


