package tui

import (
	"strings"
	"testing"
	
	"github.com/charmbracelet/bubbles/textinput"
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

func TestRenderSourceFileInputView(t *testing.T) {
	// Initialize text input
	sourceInput := textinput.New()
	sourceInput.Placeholder = "Enter path to existing resume (optional)"
	sourceInput.SetValue("/path/to/file.md")
	
	// Create model with the source path input
	model := Model{
		sourcePathInput: sourceInput,
		flagSourcePath:  "/path/from/flags.md",
	}
	
	// Get the rendered view
	sourceInputView := renderSourceFileInputView(model)
	
	// The source file input view should:
	// 1. Contain a title or heading about source file
	if !strings.Contains(sourceInputView, "Source File") {
		t.Error("Source file input view should contain a title about source files")
	}
	
	// 2. Display the textinput component
	if !strings.Contains(sourceInputView, sourceInput.View()) {
		t.Error("Source file input view should display the textinput component")
	}
	
	// 3. Show if a file path was provided via flags
	if !strings.Contains(sourceInputView, "from command line flags") {
		t.Error("Source file input view should indicate if a path was provided via flags")
	}
	
	// 4. Include instructions about file input being optional
	if !strings.Contains(sourceInputView, "optional") {
		t.Error("Source file input view should indicate that file input is optional")
	}
	
	// 5. Include keyboard shortcut hints
	if !strings.Contains(sourceInputView, "Enter to continue") {
		t.Error("Source file input view should include Enter shortcut")
	}
	
	if !strings.Contains(sourceInputView, "Ctrl+C to quit") {
		t.Error("Source file input view should include quit shortcut")
	}
	
	// Test with empty flag source path
	emptyFlagModel := Model{
		sourcePathInput: sourceInput,
		flagSourcePath:  "",
	}
	
	emptyFlagView := renderSourceFileInputView(emptyFlagModel)
	
	// 6. Should not mention flags if no flag source path was provided
	if strings.Contains(emptyFlagView, "from command line flags") {
		t.Error("Source file input view should not mention flags if no flag source path is provided")
	}
}