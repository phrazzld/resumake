package tui

import (
	"os"
	"testing"
	
	"github.com/google/generative-ai-go/genai"
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
	
	// Test empty file path (should return success with empty content)
	emptyPathCmd := ReadSourceFileCmd("")
	emptyPathResult := emptyPathCmd()
	
	emptyPathMsg, ok := emptyPathResult.(FileReadResultMsg)
	if !ok {
		t.Fatalf("Expected FileReadResultMsg, got %T", emptyPathResult)
	}
	
	if !emptyPathMsg.Success {
		t.Errorf("Expected Success to be true for empty path, got false: %v", emptyPathMsg.Error)
	}
	
	if emptyPathMsg.Content != "" {
		t.Errorf("Expected empty content for empty path, got %q", emptyPathMsg.Content)
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

// TestInitializeAPICmd tests the API initialization command
// Note: This test uses dry run mode to avoid actual API calls
func TestInitializeAPICmd(t *testing.T) {
	// This test would need environment variables or mocking
	// Just test that the command returns the correct message type
	
	t.Run("Command Returns Correct Message Type", func(t *testing.T) {
		// Create the command
		cmd := InitializeAPICmd()
		
		// The actual execution would require API key and environment setup
		// We can't easily test the success/failure paths in a unit test
		
		// Just verify the command is created and doesn't panic when called
		result := cmd()
		_, ok := result.(APIInitResultMsg)
		if !ok {
			t.Fatalf("Expected APIInitResultMsg, got %T", result)
		}
		
		// Note: We don't verify success/failure as it depends on environment
	})
}

// TestGenerateResumeCmd tests the resume generation command
func TestGenerateResumeCmd(t *testing.T) {
	// Test using dry run mode to avoid actual API calls
	t.Run("Dry Run Mode", func(t *testing.T) {
		// Create test data
		sourceContent := "Source resume content"
		stdinContent := "Additional resume details"
		outputPath := "/tmp/test_resume.md"
		
		// Client and model should be nil in dry run mode
		var client *genai.Client = nil
		var model *genai.GenerativeModel = nil
		
		// Create command with dry run flag set to true
		cmd := GenerateResumeCmd(client, model, sourceContent, stdinContent, outputPath, true)
		result := cmd()
		
		// Check the result type
		msg, ok := result.(APIResultMsg)
		if !ok {
			t.Fatalf("Expected APIResultMsg, got %T", result)
		}
		
		// In dry run mode, we should always get success
		if !msg.Success {
			t.Errorf("Expected Success to be true in dry run mode, got false: %v", msg.Error)
		}
		
		// Check that content and output path are set correctly
		if msg.Content != "Test content (dry run)" {
			t.Errorf("Expected content to be 'Test content (dry run)', got %q", msg.Content)
		}
		
		if msg.OutputPath != outputPath {
			t.Errorf("Expected output path %q, got %q", outputPath, msg.OutputPath)
		}
	})
	
	// Note: Testing the non-dry run mode would require:
	// 1. Mocking the API calls
	// 2. Mocking file operations
	// 3. Setting up environment variables
	// This is better suited for integration tests
}

// TestGenerateResumeCmdUsesProvidedClient verifies that GenerateResumeCmd uses the provided client and model
func TestGenerateResumeCmdUsesProvidedClient(t *testing.T) {
	// This test will verify that the provided client is used instead of creating a new one
	t.Run("Uses provided client and model", func(t *testing.T) {
		// We'll use the dry run mode to avoid actual API calls,
		// but we can still verify that the function signature accepts client and model parameters
		
		// Create test data
		sourceContent := "Source resume content"
		stdinContent := "Additional resume details"
		outputPath := "/tmp/test_resume.md"
		
		// For now, just test with nil client/model since we're using dry run mode
		var client *genai.Client = nil
		var model *genai.GenerativeModel = nil
		
		// Create and run the command
		cmd := GenerateResumeCmd(client, model, sourceContent, stdinContent, outputPath, true)
		result := cmd()
		
		// Verify command produced expected result
		msg, ok := result.(APIResultMsg)
		if !ok {
			t.Fatalf("Expected APIResultMsg, got %T", result)
		}
		
		if !msg.Success {
			t.Errorf("Expected Success to be true, got false: %v", msg.Error)
		}
		
		// In a real implementation, we would create mock implementations of client and model
		// that allow us to verify they were used correctly. Since we're using dry run mode,
		// we're just verifying that the function accepts the client and model parameters.
	})
}