package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ReadFile reads a markdown file from the given path
func ReadFile(path string) ([]byte, error) {
	if path == "" || path == "-" {
		// Read from stdin
		return io.ReadAll(os.Stdin)
	}

	// Read from file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return content, nil
}

// WriteFile writes content to a file
func WriteFile(path string, content []byte) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsMarkdownFile checks if the file has a markdown extension
func IsMarkdownFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".md" || ext == ".markdown" || ext == ".mdown" || ext == ".mkd"
}
