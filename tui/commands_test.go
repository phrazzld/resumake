package tui

import (
	"os"
	"testing"
)

// TestReadSourceFileCmd tests the file reading command
func TestReadSourceFileCmd(t *testing.T) {
	// Create a temporary test file
	tmpfile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	
	// Write test content to the file
	testContent := "Test file content"
	if _, err := tmpfile.Write([]byte(testContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}
	
	// Get a command for a valid file
	validCmd := ReadSourceFileCmd(tmpfile.Name())
	validResult := validCmd()
	
	// Test the result for a valid file
	fileMsg, ok := validResult.(FileReadResultMsg)
	if !ok {
		t.Fatalf("Expected FileReadResultMsg, got %T", validResult)
	}
	
	if !fileMsg.Success {
		t.Errorf("Expected Success to be true, got false: %v", fileMsg.Error)
	}
	
	if fileMsg.Content != testContent {
		t.Errorf("Expected content %q, got %q", testContent, fileMsg.Content)
	}
	
	// Get a command for an invalid file
	invalidCmd := ReadSourceFileCmd("nonexistent-file.md")
	invalidResult := invalidCmd()
	
	// Test the result for an invalid file
	invalidMsg, ok := invalidResult.(FileReadResultMsg)
	if !ok {
		t.Fatalf("Expected FileReadResultMsg, got %T", invalidResult)
	}
	
	if invalidMsg.Success {
		t.Error("Expected Success to be false, got true")
	}
	
	if invalidMsg.Error == nil {
		t.Error("Expected non-nil error for nonexistent file")
	}
}

// TestSubmitStdinInputCmd tests the stdin input command
func TestSubmitStdinInputCmd(t *testing.T) {
	content := "Test stdin input"
	cmd := SubmitStdinInputCmd(content)
	result := cmd()
	
	// Check the result type
	msg, ok := result.(StdinSubmitMsg)
	if !ok {
		t.Fatalf("Expected StdinSubmitMsg, got %T", result)
	}
	
	// Check the message content
	if msg.Content != content {
		t.Errorf("Expected content %q, got %q", content, msg.Content)
	}
}

// TestSendProgressUpdateCmd tests the progress update command
func TestSendProgressUpdateCmd(t *testing.T) {
	cmd := SendProgressUpdateCmd("Testing", "Progress update test")
	result := cmd()
	
	// Check the result type
	msg, ok := result.(ProgressUpdateMsg)
	if !ok {
		t.Fatalf("Expected ProgressUpdateMsg, got %T", result)
	}
	
	// Check the message content
	if msg.Step != "Testing" {
		t.Errorf("Expected step %q, got %q", "Testing", msg.Step)
	}
	
	if msg.Message != "Progress update test" {
		t.Errorf("Expected message %q, got %q", "Progress update test", msg.Message)
	}
}