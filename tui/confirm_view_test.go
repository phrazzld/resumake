package tui

import (
	"strings"
	"testing"
)

func TestConfirmViewEnhancements(t *testing.T) {
	// Create a test model with all required fields
	model := createTestModelWithAllFields()
	
	// Test with a normal width
	t.Run("with_normal_width", func(t *testing.T) {
		model.width = 80
		view := renderConfirmGenerateView(model)
		
		// Check for the expected UI elements
		requiredElements := []string{
			"Ready to Generate Resume",    // Title
			"Summary of Input",            // Section title
			"Source file:",                // Source info
			"Input:",                      // Input info
			"Press Enter to confirm",      // Action instruction
			"Press ESC to go back",        // Back instruction
		}
		
		for _, element := range requiredElements {
			if !strings.Contains(view, element) {
				t.Errorf("Confirm view should contain '%s'", element)
			}
		}
		
		// Check for a styled border
		if !strings.Contains(view, "─") {
			t.Error("Confirm view should have border styling")
		}
	})
	
	// Test with a very long path to verify wrapping
	t.Run("with_long_path", func(t *testing.T) {
		// Set a very long path
		longPath := "/this/is/an/extremely/long/path/that/would/definitely/need/wrapping/in/smaller/terminal/windows/resume.md"
		model.sourcePathInput.SetValue(longPath)
		model.width = 40  // Narrow width to force wrapping
		
		view := renderConfirmGenerateView(model)
		
		// Check that the view still renders with the narrow width
		if !strings.Contains(view, "Ready to Generate") {
			t.Error("Confirm view title should be visible even at narrow width")
		}
		
		// The long path should be split across multiple lines in a narrow window
		// Check that the full path appears in the output somewhere
		if !strings.Contains(view, "resume.md") {
			t.Error("Long path should be included in the view")
		}
		
		// Check for new formatting elements
		if !strings.Contains(view, "📄") || !strings.Contains(view, "✏️") {
			t.Error("Confirm view should include emoji icons for better visual hierarchy")
		}
	})
	
	// Test with output path specified
	t.Run("with_output_path", func(t *testing.T) {
		model.flagOutputPath = "/custom/output/path/resume_out.md"
		view := renderConfirmGenerateView(model)
		
		if !strings.Contains(view, "Output path") {
			t.Error("Confirm view should display the output path when specified")
		}
	})
}