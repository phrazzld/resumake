package input

import (
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Test case 1: No flags provided
	t.Run("No source flag provided", func(t *testing.T) {
		// Parse flags with no arguments
		flags, err := ParseFlagsWithArgs([]string{})
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify source is empty
		if flags.SourcePath != "" {
			t.Errorf("Expected empty source path, got %q", flags.SourcePath)
		}
	})

	// Test case 2: Source flag provided
	t.Run("Source flag provided", func(t *testing.T) {
		// Setup test with source flag
		expectedPath := "/path/to/resume.md"
		args := []string{"-source", expectedPath}
		
		// Parse flags
		flags, err := ParseFlagsWithArgs(args)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify source matches expected value
		if flags.SourcePath != expectedPath {
			t.Errorf("Expected source path %q, got %q", expectedPath, flags.SourcePath)
		}
	})

	// Test case 3: Source flag with empty value
	t.Run("Source flag with empty value", func(t *testing.T) {
		// Setup test with empty source flag
		args := []string{"-source", ""}
		
		// Parse flags
		flags, err := ParseFlagsWithArgs(args)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify source is empty
		if flags.SourcePath != "" {
			t.Errorf("Expected empty source path, got %q", flags.SourcePath)
		}
	})
	
	// Test case 4: Invalid flag
	t.Run("Invalid flag", func(t *testing.T) {
		// Setup test with invalid flag
		args := []string{"-invalid-flag"}
		
		// Parse flags
		_, err := ParseFlagsWithArgs(args)
		
		// Verify an error occurred
		if err == nil {
			t.Error("Expected error for invalid flag, got nil")
		}
	})
	
	// Test case 5: Output flag provided
	t.Run("Output flag provided", func(t *testing.T) {
		// Setup test with output flag
		expectedPath := "/path/to/output.md"
		args := []string{"-output", expectedPath}
		
		// Parse flags
		flags, err := ParseFlagsWithArgs(args)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify output matches expected value
		if flags.OutputPath != expectedPath {
			t.Errorf("Expected output path %q, got %q", expectedPath, flags.OutputPath)
		}
	})
}