package input

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadSourceFileFromFlags(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "resumake-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after tests
	
	// Test case 1: No source path in flags
	t.Run("No source path", func(t *testing.T) {
		// Create empty flags
		flags := Flags{SourcePath: ""}
		
		// Read source file from flags
		content, fileRead, err := ReadSourceFileFromFlags(flags)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify no file was read
		if fileRead {
			t.Error("Expected fileRead to be false")
		}
		
		// Verify content is empty
		if content != "" {
			t.Errorf("Expected empty content, got %q", content)
		}
	})
	
	// Test case 2: Valid source path in flags
	t.Run("Valid source path", func(t *testing.T) {
		// Create a test file
		testContent := "Test content for flags"
		filePath := filepath.Join(tempDir, "test-flags.md")
		err := os.WriteFile(filePath, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		
		// Create flags with source path
		flags := Flags{SourcePath: filePath}
		
		// Read source file from flags
		content, fileRead, err := ReadSourceFileFromFlags(flags)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify file was read
		if !fileRead {
			t.Error("Expected fileRead to be true")
		}
		
		// Verify content matches
		if content != testContent {
			t.Errorf("Expected content %q, got %q", testContent, content)
		}
	})
	
	// Test case 3: Invalid source path in flags
	t.Run("Invalid source path", func(t *testing.T) {
		// Create flags with invalid source path
		flags := Flags{SourcePath: filepath.Join(tempDir, "non-existent.md")}
		
		// Read source file from flags
		_, fileRead, err := ReadSourceFileFromFlags(flags)
		
		// Verify error occurred
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
		
		// Verify no file was read
		if fileRead {
			t.Error("Expected fileRead to be false")
		}
	})
}

func TestReadSourceFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "resumake-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after tests
	
	// Save original stdout and restore after tests
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Test case 1: Successfully read existing file
	t.Run("Read existing file", func(t *testing.T) {
		// Create a temporary test file
		testContent := "This is a test resume content.\nIt has multiple lines.\n"
		filePath := filepath.Join(tempDir, "test-resume.md")
		err := os.WriteFile(filePath, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		
		// Read the file using ReadSourceFile
		content, err := ReadSourceFile(filePath)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify content matches
		if content != testContent {
			t.Errorf("Expected content %q, got %q", testContent, content)
		}
	})

	// Test case 2: File not found
	t.Run("File not found", func(t *testing.T) {
		// Try to read a non-existent file
		nonExistentPath := filepath.Join(tempDir, "non-existent-file.md")
		
		// Read the file using ReadSourceFile
		_, err := ReadSourceFile(nonExistentPath)
		
		// Verify appropriate error is returned
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
		
		// Verify error message contains file path
		if !strings.Contains(err.Error(), nonExistentPath) {
			t.Errorf("Expected error message to contain file path, got: %v", err)
		}
	})

	// Test case 3: Empty file
	t.Run("Empty file", func(t *testing.T) {
		// Create an empty temporary file
		emptyFilePath := filepath.Join(tempDir, "empty-resume.md")
		err := os.WriteFile(emptyFilePath, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create empty test file: %v", err)
		}
		
		// Read the file using ReadSourceFile
		content, err := ReadSourceFile(emptyFilePath)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error for empty file, got %v", err)
		}
		
		// Verify empty content is returned
		if content != "" {
			t.Errorf("Expected empty content, got %q", content)
		}
	})

	// Test case 4: File with special characters
	t.Run("File with special characters", func(t *testing.T) {
		// Create a temporary file with special characters
		specialContent := "Unicode: 😊 ñ é\nSymbols: © ® ™\nTabs:\t\t\t"
		specialFilePath := filepath.Join(tempDir, "special-resume.md")
		err := os.WriteFile(specialFilePath, []byte(specialContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file with special characters: %v", err)
		}
		
		// Read the file using ReadSourceFile
		content, err := ReadSourceFile(specialFilePath)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify content matches
		if content != specialContent {
			t.Errorf("Expected content %q, got %q", specialContent, content)
		}
	})
	
	// Test case 5: Unsupported file extension (should still work but print warning)
	t.Run("Unsupported file extension", func(t *testing.T) {
		// Redirect stdout to capture warning
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		// Create a temporary file with unsupported extension
		testContent := "This is a test resume content."
		unsupportedFilePath := filepath.Join(tempDir, "resume.unsupported")
		err := os.WriteFile(unsupportedFilePath, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		
		// Read the file using ReadSourceFile
		content, err := ReadSourceFile(unsupportedFilePath)
		
		// Close the writer to get the output
		w.Close()
		outputBytes := make([]byte, 1024)
		n, _ := r.Read(outputBytes)
		output := string(outputBytes[:n])
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error despite unsupported extension, got %v", err)
		}
		
		// Verify content matches
		if content != testContent {
			t.Errorf("Expected content %q, got %q", testContent, content)
		}
		
		// Verify warning was printed
		if !strings.Contains(output, "Warning") || !strings.Contains(output, "unsupported file extension") {
			t.Errorf("Expected warning about unsupported file extension, got: %q", output)
		}
	})
	
	// Test case 6: Non-regular file
	t.Run("Non-regular file", func(t *testing.T) {
		// Create a subdirectory to test with
		dirPath := filepath.Join(tempDir, "subdir")
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}
		
		// Try to read the directory as a file
		_, err = ReadSourceFile(dirPath)
		
		// Verify error about non-regular file
		if err == nil {
			t.Error("Expected error for non-regular file, got nil")
		}
		
		if !strings.Contains(err.Error(), "not a regular file") {
			t.Errorf("Expected error about non-regular file, got: %v", err)
		}
	})
}