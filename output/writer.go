package output

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultOutputPath defines the default path for writing the generated resume.
// This path is used when the user doesn't specify an output path via command-line flags.
var DefaultOutputPath = "resume_out.md"

// WriteToFile writes content to a file at the specified path.
// It creates the file if it doesn't exist or overwrites it if it does.
// This function also ensures the target directory exists, creating it if necessary.
//
// Parameters:
//   - path: The absolute or relative path where the file should be written
//   - content: The string content to write to the file
//
// Returns:
//   - error: An error if directory creation or file writing fails, nil otherwise
//
// Example:
//
//	err := output.WriteToFile("./resumes/my_resume.md", markdownContent)
//	if err != nil {
//	    log.Fatalf("Failed to write file: %v", err)
//	}
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

// WriteOutput writes content to the output file, handling path selection logic.
// It's a higher-level function that decides which path to use (provided or default)
// and then calls WriteToFile to perform the actual writing.
//
// Parameters:
//   - content: The string content to write to the file
//   - outputPath: The path where the file should be written, or empty to use default
//
// Returns:
//   - string: The actual path where the content was written (useful for reporting)
//   - error: An error if file writing fails, nil otherwise
//
// Example:
//
//	path, err := output.WriteOutput(markdownContent, flags.OutputPath)
//	if err != nil {
//	    log.Fatalf("Failed to write output: %v", err)
//	}
//	fmt.Printf("Resume written to: %s\n", path)
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