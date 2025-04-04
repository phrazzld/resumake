package input

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MaxFileSize is the maximum allowed file size in bytes (10MB).
// Files larger than this limit will be rejected to prevent memory issues.
const MaxFileSize = 10 * 1024 * 1024

// SupportedFileExtensions contains the allowed file extensions for resume files.
// The application will warn but not block if the file has a different extension.
var SupportedFileExtensions = []string{".txt", ".md", ".markdown"}

// ReadSourceFile reads the content of a file at the given path.
// It performs several validation checks before reading the file:
// - Verifies the file exists and is accessible
// - Confirms it's a regular file (not a directory or special file)
// - Ensures the file size is within the maximum allowed limit
// - Warns if the file extension is not in the supported list
//
// Parameters:
//   - filePath: The path to the file to read
//
// Returns:
//   - string: The file content as a string
//   - error: Any error that occurred during validation or reading
//
// Example:
//
//	content, err := input.ReadSourceFile("my_resume.md")
//	if err != nil {
//	    log.Fatalf("Error reading source file: %v", err)
//	}
func ReadSourceFile(filePath string) (string, error) {
	// Check if the file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file does not exist: %s", filePath)
		}
		return "", fmt.Errorf("error accessing file %s: %w", filePath, err)
	}
	
	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", filePath)
	}
	
	// Check file size
	if fileInfo.Size() > MaxFileSize {
		return "", fmt.Errorf("file size exceeds the maximum allowed size of %d bytes: %s", MaxFileSize, filePath)
	}
	
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	validExtension := false
	for _, supported := range SupportedFileExtensions {
		if ext == supported {
			validExtension = true
			break
		}
	}
	
	// Only warn about extension, don't block
	if !validExtension {
		fmt.Printf("Warning: %s has an unsupported file extension. Supported extensions are: %s\n", 
			filePath, strings.Join(SupportedFileExtensions, ", "))
	}
	
	// Read the file content
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	
	// Convert to string and return
	return string(contentBytes), nil
}

// ReadSourceFileFromFlags reads a source file if one is specified in the flags.
// It provides a convenient way to conditionally read a file based on command-line flags.
// If no source path is specified in the flags, it returns empty content.
//
// Parameters:
//   - flags: The parsed command-line flags
//
// Returns:
//   - string: The file content as a string (empty if no file specified)
//   - bool: True if a file was read, false if no file was specified
//   - error: Any error that occurred during file reading
//
// Example:
//
//	content, fileRead, err := input.ReadSourceFileFromFlags(flags)
//	if err != nil {
//	    log.Fatalf("Error reading source file: %v", err)
//	}
//	if fileRead {
//	    fmt.Printf("Successfully read source file: %s\n", flags.SourcePath)
//	}
func ReadSourceFileFromFlags(flags Flags) (string, bool, error) {
	// If no source file is specified, return empty content
	if flags.SourcePath == "" {
		return "", false, nil
	}
	
	// Read the source file
	content, err := ReadSourceFile(flags.SourcePath)
	if err != nil {
		return "", false, err
	}
	
	// Return the content and indicate a file was read
	return content, true, nil
}