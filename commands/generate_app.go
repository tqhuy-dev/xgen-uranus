package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tqhuy-dev/xgen-uranus/commands/generate_tpl"
)

// AppConfig holds configuration for app generation
type AppConfig struct {
	Name   string
	Module string
	Dir    string
	Port   int
}

var (
	appName    string
	moduleName string
	skipInit   string
)

var generateAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Generate a new application",
	Long: `Generate a new application with the specified name.

This command creates a new application scaffold with all necessary files
and directory structure for a microservice application.

Examples:
  uranus generate app --name simple_app --module github.com/user/simple_app
  uranus generate app -n my_service -m github.com/user/my_service`,
	Run: func(cmd *cobra.Command, args []string) {
		if appName == "" {
			fmt.Println("Error: --name flag is required")
			cmd.Help()
			os.Exit(1)
		}

		if err := generateApp(appName, moduleName, skipInit == "true"); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating app: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	generateAppCmd.Flags().StringVarP(&appName, "name", "n", "", "Name of the application to generate (required)")
	generateAppCmd.MarkFlagRequired("name")

	generateAppCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Module of the application to generate (required)")
	generateAppCmd.MarkFlagRequired("module")

	generateAppCmd.Flags().StringVarP(&skipInit, "skip_init", "s", "", "skip make init")
}

// generateApp orchestrates the app generation process
func generateApp(name string, module string, skipInit bool) error {
	config, err := newAppConfig(name, module)
	if err != nil {
		return err
	}

	fmt.Printf("🚀 Generating application: %s\n", config.Name)

	// Step 1: Create directory structure
	if err := createDirectories(config); err != nil {
		return err
	}

	// Step 2: Generate files from templates
	if err := generateTemplateFiles(config); err != nil {
		return err
	}

	// Step 3: Copy static files
	if err := copyStaticFiles(config); err != nil {
		return err
	}

	// Step 4: Generate config files
	if err := generateConfigFiles(config); err != nil {
		return err
	}

	fmt.Printf("\n✅ Application '%s' generated successfully!\n", config.Name)

	// Step 5: Run post-generation commands
	if !skipInit {
		if err := runPostGeneration(config); err != nil {
			return err
		}
	}

	printNextSteps(config)
	return nil
}

// newAppConfig creates and validates app configuration
func newAppConfig(name, module string) (*AppConfig, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("app name cannot be empty")
	}

	dir := filepath.Dir(name)
	baseName := filepath.Base(name)

	appDir := name
	if dir == "." {
		appDir = baseName
	}

	return &AppConfig{
		Name:   baseName,
		Module: module,
		Dir:    appDir,
		Port:   10000,
	}, nil
}

// createDirectories creates the app directory structure
func createDirectories(config *AppConfig) error {
	dirs := []string{
		config.Dir,
		filepath.Join(config.Dir, "cmd"),
		filepath.Join(config.Dir, "internal"),
		filepath.Join(config.Dir, "internal", "services"),
		filepath.Join(config.Dir, "internal", "dto"),
		filepath.Join(config.Dir, "internal", "repository"),
		filepath.Join(config.Dir, "pkg"),
		filepath.Join(config.Dir, "pkg", "proto"),
		filepath.Join(config.Dir, "configs"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
		fmt.Printf("  📁 Created: %s\n", d)
	}

	return nil
}

// generateTemplateFiles generates files from templates
func generateTemplateFiles(config *AppConfig) error {
	// Generate main.go
	mainData := generate_tpl.MainFileData{
		AppName:    config.Name,
		ModuleName: config.Module,
		Port:       config.Port,
	}
	if err := generate_tpl.MainFile(config.Dir, mainData); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}

	return nil
}

// copyStaticFiles copies static files to the app directory
func copyStaticFiles(config *AppConfig) error {
	// Copy grpc_third_party folder
	if err := generate_tpl.CopyGrpcThirdParty(config.Dir); err != nil {
		return fmt.Errorf("failed to copy grpc_third_party: %w", err)
	}

	// Copy Makefile
	if err := generate_tpl.CopyMakefile(config.Dir); err != nil {
		return fmt.Errorf("failed to copy Makefile: %w", err)
	}

	return nil
}

// generateConfigFiles generates configuration files
func generateConfigFiles(config *AppConfig) error {
	if err := generateGoMod(config); err != nil {
		return err
	}

	if err := generateReadme(config); err != nil {
		return err
	}

	if err := generateConfigYaml(config); err != nil {
		return err
	}

	return nil
}

// generateGoMod creates the go.mod file
func generateGoMod(config *AppConfig) error {
	content := fmt.Sprintf(`module %s

go 1.25

require (
)
`, config.Module)

	path := filepath.Join(config.Dir, "go.mod")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}
	fmt.Printf("  📄 Created: %s\n", path)
	return nil
}

// generateReadme creates the README.md file
func generateReadme(config *AppConfig) error {
	content := fmt.Sprintf(`# %s

A microservice application generated by Uranus CLI.

## Getting Started

### Prerequisites
- Go 1.21 or higher
- protoc (Protocol Buffer Compiler)
- make

### Running the application

`+"```bash"+`
cd %s
make run
`+"```"+`

### Build

`+"```bash"+`
make build          # Build for current OS
make build-all      # Build for all platforms
`+"```"+`

### Generate Protobuf

`+"```bash"+`
make init           # Install protoc plugins (first time)
make proto          # Generate protobuf code
`+"```"+`

Run `+"`make help`"+` to see all available commands.

## Project Structure

`+"```"+`
%s/
├── cmd/                  # Application entry points
│   └── main.go           # Entry point
├── internal/             # Private application code
│   ├── services/         # Business logic
│   ├── dto/              # Data transfer objects
│   └── repository/       # Data access layer
├── pkg/                  # Public libraries
│   └── proto/            # Generated protobuf code
├── configs/              # Configuration files
├── grpc_third_party/     # Third-party proto definitions
└── Makefile              # Build automation
`+"```"+`

## License
MIT
`, config.Name, config.Dir, config.Name)

	path := filepath.Join(config.Dir, "README.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}
	fmt.Printf("  📄 Created: %s\n", path)
	return nil
}

// generateConfigYaml creates the config.yaml file
func generateConfigYaml(config *AppConfig) error {
	content := fmt.Sprintf(`# Configuration for %s

server:
  host: "0.0.0.0"
  port: 8080
  grpc_port: %d

database:
  host: "localhost"
  port: 5432
  name: "%s_db"
  user: "postgres"
  password: ""

logging:
  level: "info"
  format: "json"
`, config.Name, config.Port, config.Name)

	path := filepath.Join(config.Dir, "configs", "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create config.yaml: %w", err)
	}
	fmt.Printf("  📄 Created: %s\n", path)
	return nil
}

// runPostGeneration runs post-generation commands (make init)
func runPostGeneration(config *AppConfig) error {
	absDir, err := filepath.Abs(config.Dir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	fmt.Printf("\n🔧 Running 'make init' in %s...\n", absDir)

	cmd := exec.Command("make", "init")
	cmd.Dir = absDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run 'make init': %w", err)
	}

	return nil
}

// printNextSteps prints the next steps for the user
func printNextSteps(config *AppConfig) {
	fmt.Printf("\n🎉 All done! Your app is ready.\n")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", config.Dir)
	fmt.Printf("  make run\n")
}
