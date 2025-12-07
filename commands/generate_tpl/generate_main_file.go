package generate_tpl

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed tpl/*
var templateFS embed.FS

// MainFileData contains data for main.go template
type MainFileData struct {
	AppName     string
	PackageName string
	ModuleName  string
	Port        int
}

// MainFile generates main.go file from template
func MainFile(outputDir string, data MainFileData) error {
	// Set default values
	if data.Port == 0 {
		data.Port = 10000
	}
	if data.PackageName == "" {
		data.PackageName = "main"
	}

	// Read template file
	tmplContent, err := templateFS.ReadFile("tpl/main.gotmpl")
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Create template with sprig functions
	tmpl, err := template.New("main.go").
		Funcs(sprig.TxtFuncMap()).
		Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Ensure output directory exists
	cmdDir := filepath.Join(outputDir, "cmd")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		return fmt.Errorf("failed to create cmd directory: %w", err)
	}

	// Write output file
	outputPath := filepath.Join(cmdDir, "main.go")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	fmt.Printf("  📄 Created: %s\n", outputPath)
	return nil
}

// CopyMakefile copies the static Makefile to the target directory
func CopyMakefile(outputDir string) error {
	// Read Makefile from embedded FS
	content, err := templateFS.ReadFile("tpl/Makefile")
	if err != nil {
		return fmt.Errorf("failed to read Makefile: %w", err)
	}

	// Write output file
	outputPath := filepath.Join(outputDir, "Makefile")
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write Makefile: %w", err)
	}

	fmt.Printf("  📄 Created: %s\n", outputPath)
	return nil
}

// GenerateFromTemplate is a generic function to generate file from any template
func GenerateFromTemplate(templateName string, outputPath string, data interface{}) error {
	// Read template file
	tmplContent, err := templateFS.ReadFile("tpl/" + templateName)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	// Create template with sprig functions
	tmpl, err := template.New(templateName).
		Funcs(sprig.TxtFuncMap()).
		Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("  📄 Created: %s\n", outputPath)
	return nil
}
