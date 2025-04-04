package output

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultOutputPath defines the default path for writing the generated resume
var DefaultOutputPath = "resume_out.md"

// WriteToFile writes content to a file at the specified path.
// It creates the file if it doesn't exist or overwrites it if it does.
// Returns an error if the file cannot be written.
func WriteToFile(path string, content string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := ensureDirectoryExists(dir); err != nil {
		return fmt.Errorf("failed to ensure directory exists: %w", err)
	}
	
	// Write the content to the file
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	
	return nil
}

// ensureDirectoryExists checks if the directory exists and creates it if it doesn't.
// Returns an error if the directory cannot be created.
func ensureDirectoryExists(dirPath string) error {
	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err == nil {
		// Path exists, check if it's a directory
		if !info.IsDir() {
			return fmt.Errorf("%s exists but is not a directory", dirPath)
		}
		return nil // Directory exists
	}
	
	// If the error is something other than "not exists", return it
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory: %w", err)
	}
	
	// Create the directory and any necessary parents
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	return nil
}

// WriteOutput writes content to the output file.
// If outputPath is empty, the DefaultOutputPath is used.
// Returns the path that was written to and any error that occurred.
func WriteOutput(content string, outputPath string) (string, error) {
	// Use default path if none provided
	if outputPath == "" {
		outputPath = DefaultOutputPath
	}
	
	// Write the content to the file
	err := WriteToFile(outputPath, content)
	if err != nil {
		return "", fmt.Errorf("failed to write output: %w", err)
	}
	
	return outputPath, nil
}