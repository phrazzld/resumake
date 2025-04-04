package output

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestEnvironment creates a temporary directory for testing file operations
// and returns the clean-up function that removes the temporary directory
func setupTestEnvironment(t *testing.T) (string, func()) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "resumake-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	
	// Return the cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}
	
	return tempDir, cleanup
}

func TestWriteToFile(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	testContent := "# Test Resume\n\n## Skills\n\n- Go\n- Testing"
	
	tests := []struct {
		name        string
		path        string
		content     string
		setup       func(string) error
		shouldError bool
	}{
		{
			name:        "write to new file",
			path:        filepath.Join(tempDir, "new_file.md"),
			content:     testContent,
			setup:       nil,
			shouldError: false,
		},
		{
			name:    "overwrite existing file",
			path:    filepath.Join(tempDir, "existing_file.md"),
			content: testContent,
			setup: func(path string) error {
				return os.WriteFile(path, []byte("Old content"), 0644)
			},
			shouldError: false,
		},
		{
			name:    "write to read-only directory",
			path:    filepath.Join(tempDir, "readonly", "file.md"),
			content: testContent,
			setup: func(path string) error {
				dir := filepath.Dir(path)
				if err := os.MkdirAll(dir, 0444); err != nil {
					return err
				}
				return nil
			},
			shouldError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the test case
			if tt.setup != nil {
				if err := tt.setup(tt.path); err != nil {
					t.Fatalf("Failed to setup test: %v", err)
				}
			}
			
			err := WriteToFile(tt.path, tt.content)
			
			// Check if error matches expectation
			if (err != nil) != tt.shouldError {
				t.Errorf("WriteToFile() error = %v, shouldError = %v", err, tt.shouldError)
			}
			
			// If no error is expected, verify the file content
			if !tt.shouldError {
				content, err := os.ReadFile(tt.path)
				if err != nil {
					t.Fatalf("Failed to read test file: %v", err)
				}
				
				if string(content) != tt.content {
					t.Errorf("File content doesn't match. Got %q, want %q", string(content), tt.content)
				}
			}
		})
	}
}

func TestWriteOutput(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	defaultOutputPath := filepath.Join(tempDir, "resume_out.md")
	customOutputPath := filepath.Join(tempDir, "custom_output.md")
	testContent := "# Test Resume\n\n## Skills\n\n- Go\n- Testing"
	
	tests := []struct {
		name        string
		outputPath  string
		content     string
		shouldError bool
		checkPath   string
	}{
		{
			name:        "write to default output path",
			outputPath:  "",
			content:     testContent,
			shouldError: false,
			checkPath:   defaultOutputPath,
		},
		{
			name:        "write to custom output path",
			outputPath:  customOutputPath,
			content:     testContent,
			shouldError: false,
			checkPath:   customOutputPath,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override the DefaultOutputPath for testing
			origDefaultPath := DefaultOutputPath
			DefaultOutputPath = defaultOutputPath
			defer func() { DefaultOutputPath = origDefaultPath }()
			
			// Call the function
			outputPath, err := WriteOutput(tt.content, tt.outputPath)
			
			// Check if error matches expectation
			if (err != nil) != tt.shouldError {
				t.Errorf("WriteOutput() error = %v, shouldError = %v", err, tt.shouldError)
			}
			
			// Check the returned path
			if outputPath != tt.checkPath {
				t.Errorf("WriteOutput() returned path = %q, want %q", outputPath, tt.checkPath)
			}
			
			// If no error is expected, verify the file content
			if !tt.shouldError {
				content, err := os.ReadFile(tt.checkPath)
				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}
				
				if string(content) != tt.content {
					t.Errorf("File content doesn't match. Got %q, want %q", string(content), tt.content)
				}
			}
		})
	}
}

func TestEnsureDirectoryExists(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	tests := []struct {
		name        string
		path        string
		shouldError bool
	}{
		{
			name:        "existing directory",
			path:        tempDir,
			shouldError: false,
		},
		{
			name:        "new directory",
			path:        filepath.Join(tempDir, "new_dir"),
			shouldError: false,
		},
		{
			name:        "nested directory",
			path:        filepath.Join(tempDir, "nested", "dir"),
			shouldError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ensureDirectoryExists(tt.path)
			
			// Check if error matches expectation
			if (err != nil) != tt.shouldError {
				t.Errorf("ensureDirectoryExists() error = %v, shouldError = %v", err, tt.shouldError)
			}
			
			// Verify the directory exists
			if !tt.shouldError {
				info, err := os.Stat(tt.path)
				if err != nil {
					t.Fatalf("Failed to stat directory: %v", err)
				}
				
				if !info.IsDir() {
					t.Errorf("Path is not a directory: %s", tt.path)
				}
			}
		})
	}
}