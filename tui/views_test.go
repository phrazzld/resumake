package tui

import (
	"strings"
	"testing"
)

func TestRenderWelcomeView(t *testing.T) {
	// Test with valid API key
	model := Model{
		apiKeyOk: true,
	}
	
	validKeyView := renderWelcomeView(model)
	
	// The welcome view with a valid API key should:
	// 1. Contain a welcome message
	if !strings.Contains(validKeyView, "Welcome to Resumake") {
		t.Error("Welcome view should contain a welcome message")
	}
	
	// 2. Indicate that the API key is valid
	if !strings.Contains(validKeyView, "API key is valid") {
		t.Error("Welcome view should indicate API key is valid")
	}
	
	// 3. Include instructions on how to proceed
	if !strings.Contains(validKeyView, "Press Enter to continue") {
		t.Error("Welcome view should include instructions to proceed")
	}
	
	// 4. Include information about what the application does
	if !strings.Contains(validKeyView, "This tool helps you create a professional resume") {
		t.Error("Welcome view should include information about the application purpose")
	}
	
	// 5. Include keyboard shortcut info
	if !strings.Contains(validKeyView, "Ctrl+C") {
		t.Error("Welcome view should include keyboard shortcut information")
	}
	
	// Test with invalid API key
	model = Model{
		apiKeyOk: false,
	}
	
	invalidKeyView := renderWelcomeView(model)
	
	// The welcome view with an invalid API key should:
	// 1. Contain a welcome message
	if !strings.Contains(invalidKeyView, "Welcome to Resumake") {
		t.Error("Welcome view should contain a welcome message")
	}
	
	// 2. Indicate that the API key is invalid
	if !strings.Contains(invalidKeyView, "API key is missing or invalid") {
		t.Error("Welcome view should indicate API key is missing or invalid")
	}
	
	// 3. Include instructions on how to set the API key
	if !strings.Contains(invalidKeyView, "GEMINI_API_KEY") {
		t.Error("Welcome view should include instructions to set the API key")
	}
}

// Additional tests for other view functions would be added here