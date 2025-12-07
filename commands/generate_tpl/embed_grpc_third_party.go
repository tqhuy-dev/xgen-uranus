package generate_tpl

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed grpc_third_party/*
var grpcThirdPartyFS embed.FS

// CopyGrpcThirdParty copies the embedded grpc_third_party folder to the target directory
func CopyGrpcThirdParty(targetDir string) error {
	destDir := filepath.Join(targetDir, "grpc_third_party")

	err := fs.WalkDir(grpcThirdPartyFS, "grpc_third_party", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel("grpc_third_party", path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			// Create directory
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			return nil
		}

		// Read file from embedded FS
		content, err := grpcThirdPartyFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Ensure parent directory exists
		parentDir := filepath.Dir(destPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create parent directory %s: %w", parentDir, err)
		}

		// Write file to destination
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy grpc_third_party: %w", err)
	}

	fmt.Printf("  📁 Copied: %s\n", destDir)
	return nil
}

