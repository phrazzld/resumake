package tui

import (
	"strings"
	"testing"
	
	"github.com/charmbracelet/bubbles/textinput"
)

func TestEnhancedSourceFileInputView(t *testing.T) {
	// Create a text input with a sample path
	sourceInput := textinput.New()
	sourceInput.Placeholder = "Enter path to existing resume (optional)"
	sourceInput.SetValue("/path/to/resume.md")
	
	// Test case 1: Model with flag source path (pre-filled from CLI)
	t.Run("with flag source path", func(t *testing.T) {
		// Create model with the source path input and flag path
		model := Model{
			sourcePathInput: sourceInput,
			flagSourcePath:  "/path/from/flags.md",
			width:           80,
			height:          24,
		}
		
		// Get the rendered view
		sourceInputView := renderSourceFileInputView(model)
		
		// Check for required elements
		requiredElements := []string{
			"Source File",                         // Title
			"/path/to/resume.md",                  // Input value
			"optional",                            // Optional indication
			"from command line",                   // Flag indication
			"Enter",                               // Continue instruction
		}
		
		for _, element := range requiredElements {
			if !strings.Contains(sourceInputView, element) {
				t.Errorf("Source file input view should contain '%s'", element)
			}
		}
		
		// Check for supported file formats
		if !strings.Contains(sourceInputView, ".md") || 
		   !strings.Contains(sourceInputView, ".txt") {
			t.Errorf("Source file input view should mention supported file formats")
		}
		
		// Check for an example path
		examplePatterns := []string{"/path/to/", "resume.md", "example"}
		foundExample := false
		for _, pattern := range examplePatterns {
			if strings.Contains(sourceInputView, pattern) {
				foundExample = true
				break
			}
		}
		if !foundExample {
			t.Errorf("Source file input view should include an example path")
		}
	})
	
	// Test case 2: Model without flag source path
	t.Run("without flag source path", func(t *testing.T) {
		// Create model without flag source path
		model := Model{
			sourcePathInput: sourceInput,
			flagSourcePath:  "",
			width:           80,
			height:          24,
		}
		
		// Get the rendered view
		sourceInputView := renderSourceFileInputView(model)
		
		// Should not mention flags if no flag source path was provided
		if strings.Contains(sourceInputView, "from command line flags") {
			t.Error("Source file input view should not mention flags if no flag source path is provided")
		}
		
		// Check for a help section or tips
		helpTerms := []string{"tip", "help", "suggestion", "hint"}
		foundHelp := false
		for _, term := range helpTerms {
			if strings.Contains(strings.ToLower(sourceInputView), term) {
				foundHelp = true
				break
			}
		}
		if !foundHelp {
			t.Errorf("Source file input view should include help text or tips")
		}
		
		// Should have some spacing between sections
		emptyLines := 0
		for _, line := range strings.Split(sourceInputView, "\n") {
			if strings.TrimSpace(line) == "" {
				emptyLines++
			}
		}
		if emptyLines < 2 {
			t.Errorf("Source file input view should have adequate spacing between sections")
		}
	})
	
	// Test case 3: Check for visual structure and formatting
	t.Run("visual elements", func(t *testing.T) {
		model := Model{
			sourcePathInput: sourceInput,
			width:           80,
			height:          24,
		}
		
		sourceInputView := renderSourceFileInputView(model)
		
		// Check for styled elements - at least one border
		if !strings.Contains(sourceInputView, "│") && 
		   !strings.Contains(sourceInputView, "┌") && 
		   !strings.Contains(sourceInputView, "└") {
			t.Errorf("Source file input view should use styled borders for visual structure")
		}
		
		// Check for purpose explanation
		purposeTerms := []string{"enhance", "existing", "improve", "resume"}
		purposeFound := 0
		for _, term := range purposeTerms {
			if strings.Contains(strings.ToLower(sourceInputView), term) {
				purposeFound++
			}
		}
		if purposeFound < 2 {
			t.Errorf("Source file input view should explain the purpose of providing a source file")
		}
	})
}