package main_test

import (
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