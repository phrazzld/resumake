package input

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MaxFileSize is the maximum allowed file size in bytes (10MB)
const MaxFileSize = 10 * 1024 * 1024

// SupportedFileExtensions contains the allowed file extensions
var SupportedFileExtensions = []string{".txt", ".md", ".markdown"}

// ReadSourceFile reads the content of a file at the given path.
// Returns the file content as a string and any error that occurred.
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
// Returns the file content as a string, a boolean indicating if a file was read,
// and any error that occurred.
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