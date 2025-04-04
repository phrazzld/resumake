package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestGoModExists verifies that a valid go.mod file exists in the project root
func TestGoModExists(t *testing.T) {
	// Check if go.mod file exists
	goModContent, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("Failed to read go.mod file: %v", err)
	}
	
	// Verify go.mod contains a valid module declaration
	if !strings.Contains(string(goModContent), "module github.com/phrazzld/resumake") {
		t.Error("go.mod file does not contain the correct module declaration")
	}
}

// TestModuleCompiles verifies that the module can compile
func TestModuleCompiles(t *testing.T) {
	// Create a temporary main.go file if it doesn't exist
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		tempMain := []byte("package main\n\nfunc main() {\n\t// Placeholder\n}\n")
		err = os.WriteFile("main.go", tempMain, 0644)
		if err != nil {
			t.Fatalf("Failed to create temporary main.go file: %v", err)
		}
		defer os.Remove("main.go") // Clean up temporary file
	}
	
	// Run go build to test compilation
	cmd := exec.Command("go", "build", "-o", "resumake-test")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Module fails to compile: %v", err)
	}
	
	// Clean up the binary
	os.Remove("resumake-test")
}

// TestHelpFlagExitsCleanly verifies that running the application with
// a help flag (-h or --help) displays usage information and exits with code 0
func TestHelpFlagExitsCleanly(t *testing.T) {
	// Build the binary for testing
	cmd := exec.Command("go", "build", "-o", "resumake-test-help")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer os.Remove("resumake-test-help")
	
	// Test cases for different help flags
	testCases := []struct {
		name     string
		args     []string
	}{
		{"long help flag", []string{"--help"}},
		{"short help flag", []string{"-h"}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run the binary with help flag
			cmd := exec.Command("./resumake-test-help", tc.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			
			// Run the command and check its exit status
			err := cmd.Run()
			
			// It should exit with status code 0 (success)
			if err != nil {
				t.Errorf("Expected clean exit (code 0) with %v flag, got error: %v", tc.args[0], err)
			}
			
			// Combine stdout and stderr for checks
			output := stdout.String() + stderr.String()
			
			// It should contain usage information
			if !strings.Contains(output, "Usage of") {
				t.Errorf("Help output should contain usage information, got: %s", output)
			}
			
			// It should contain flag descriptions
			if !strings.Contains(output, "-source") || !strings.Contains(output, "-output") {
				t.Errorf("Help output should list command flags, got: %s", output)
			}
			
			// It should NOT contain an error message
			if strings.Contains(output, "Error parsing flags:") {
				t.Errorf("Help output should not contain error messages, got: %s", output)
			}
		})
	}
}

// TestInvalidFlagShowsError verifies that the program still correctly
// handles invalid flags by showing an error and exiting with non-zero code
func TestInvalidFlagShowsError(t *testing.T) {
	// Build the binary for testing if it doesn't exist
	if _, err := os.Stat("resumake-test-help"); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", "resumake-test-help")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to build test binary: %v", err)
		}
	}
	defer os.Remove("resumake-test-help")
	
	// Run the binary with an invalid flag
	cmd := exec.Command("./resumake-test-help", "--invalid-flag")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Run the command and check its exit status
	err := cmd.Run()
	
	// It should exit with non-zero status code
	if err == nil {
		t.Error("Expected non-zero exit code with invalid flag, got success (code 0)")
	}
	
	// Combine stdout and stderr for checks
	output := stdout.String() + stderr.String()
	
	// It should contain an error message
	if !strings.Contains(output, "Error parsing flags:") {
		t.Errorf("Invalid flag should produce error message, got: %s", output)
	}
}