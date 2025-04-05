package tui

import (
	"errors"
	"testing"
)

func TestFileReadResultMsg(t *testing.T) {
	// Test successful file read
	successMsg := FileReadResultMsg{
		Success: true,
		Content: "File content",
		Error:   nil,
	}
	
	if !successMsg.Success {
		t.Error("Expected Success to be true")
	}
	
	if successMsg.Content != "File content" {
		t.Errorf("Expected Content to be 'File content', got %q", successMsg.Content)
	}
	
	if successMsg.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", successMsg.Error)
	}
	
	// Test failed file read
	testErr := errors.New("file not found")
	failMsg := FileReadResultMsg{
		Success: false,
		Content: "",
		Error:   testErr,
	}
	
	if failMsg.Success {
		t.Error("Expected Success to be false")
	}
	
	if failMsg.Content != "" {
		t.Errorf("Expected Content to be empty, got %q", failMsg.Content)
	}
	
	if failMsg.Error != testErr {
		t.Errorf("Expected Error to be %v, got %v", testErr, failMsg.Error)
	}
}

func TestAPIResultMsg(t *testing.T) {
	// Test API success
	successMsg := APIResultMsg{
		Success:      true,
		Content:      "Generated content",
		OutputPath:   "output.md",
		Error:        nil,
	}
	
	if !successMsg.Success {
		t.Error("Expected Success to be true")
	}
	
	if successMsg.Content != "Generated content" {
		t.Errorf("Expected Content to be 'Generated content', got %q", successMsg.Content)
	}
	
	if successMsg.OutputPath != "output.md" {
		t.Errorf("Expected OutputPath to be 'output.md', got %q", successMsg.OutputPath)
	}
	
	if successMsg.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", successMsg.Error)
	}
	
	// Test API failure
	testErr := errors.New("API error")
	failMsg := APIResultMsg{
		Success:      false,
		Content:      "",
		OutputPath:   "",
		Error:        testErr,
	}
	
	if failMsg.Success {
		t.Error("Expected Success to be false")
	}
	
	if failMsg.Content != "" {
		t.Errorf("Expected Content to be empty, got %q", failMsg.Content)
	}
	
	if failMsg.OutputPath != "" {
		t.Errorf("Expected OutputPath to be empty, got %q", failMsg.OutputPath)
	}
	
	if failMsg.Error != testErr {
		t.Errorf("Expected Error to be %v, got %v", testErr, failMsg.Error)
	}
}


func TestStdinSubmitMsg(t *testing.T) {
	msg := StdinSubmitMsg{
		Content: "User input",
	}
	
	if msg.Content != "User input" {
		t.Errorf("Expected Content to be 'User input', got %q", msg.Content)
	}
}