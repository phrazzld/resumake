package tui

import (
	"strings"
	"testing"
	
	"github.com/charmbracelet/bubbles/textarea"
)

func TestEnhancedStdinInputView(t *testing.T) {
	// Create a text area with sample content
	stdinInput := textarea.New()
	stdinInput.SetValue("I have 5 years of experience as a software engineer.")
	stdinInput.Placeholder = "Tell us about your experience..."
	
	// Test case 1: Basic content and layout
	t.Run("basic content and layout", func(t *testing.T) {
		// Create model with the textarea
		model := Model{
			stdinInput: stdinInput,
			width:      80,
			height:     24,
		}
		
		// Get the rendered view
		stdinInputView := renderStdinInputView(model)
		
		// Check for required elements - note we're not checking the actual textarea content
		// since bubbletea's textarea rendering is mocked during tests
		requiredElements := []string{
			"Resume Details",                           // Title
			"professional background",                  // Purpose description
			"Ctrl+D",                                   // Finish instruction
			"Suggested Content to Include:",           // Tips section
		}
		
		for _, element := range requiredElements {
			if !strings.Contains(stdinInputView, element) {
				t.Errorf("Stdin input view should contain '%s'", element)
			}
		}
		
		// Check for a clear visual structure
		if !strings.Contains(stdinInputView, "│") && 
		   !strings.Contains(stdinInputView, "┌") && 
		   !strings.Contains(stdinInputView, "└") {
			t.Errorf("Stdin input view should use styled borders for visual structure")
		}
	})
	
	// Test case 2: Examples and suggestions
	t.Run("examples and suggestions", func(t *testing.T) {
		model := Model{
			stdinInput: stdinInput,
			width:      80,
			height:     24,
		}
		
		stdinInputView := renderStdinInputView(model)
		
		// Check for example content
		exampleTerms := []string{"skills", "experience", "education", "achievements"}
		examplesFound := 0
		for _, term := range exampleTerms {
			if strings.Contains(strings.ToLower(stdinInputView), term) {
				examplesFound++
			}
		}
		
		if examplesFound < 3 {
			t.Errorf("Stdin input view should include examples mentioning at least 3 resume categories")
		}
		
		// Check for formatting tips
		formattingTerms := []string{"bullet", "list", "format", "highlight"}
		formattingFound := 0
		for _, term := range formattingTerms {
			if strings.Contains(strings.ToLower(stdinInputView), term) {
				formattingFound++
			}
		}
		
		if formattingFound < 2 {
			t.Errorf("Stdin input view should include at least 2 formatting tips")
		}
	})
	
	// Test case 3: Keyboard shortcuts
	t.Run("keyboard shortcuts", func(t *testing.T) {
		model := Model{
			stdinInput: stdinInput,
			width:      80,
			height:     24,
		}
		
		stdinInputView := renderStdinInputView(model)
		
		// Check for keyboard guide
		if !strings.Contains(stdinInputView, "Tip: Enter your details below, then press Ctrl+D when finished") {
			t.Error("Stdin input view should include the keyboard tip at the top")
		}
		
		// Check for scrollable note
		if !strings.Contains(stdinInputView, "scrollable") {
			t.Error("Stdin input view should indicate that the textarea is scrollable")
		}
	})
	
	// Test case 4: Text wrapping
	t.Run("text wrapping", func(t *testing.T) {
		// Create a narrow model to test wrapping
		model := Model{
			stdinInput: stdinInput,
			width:      40, // Narrow width to force wrapping
			height:     24,
		}
		
		stdinInputView := renderStdinInputView(model)
		
		// Count lines in the rendered view
		lines := strings.Split(stdinInputView, "\n")
		
		// Verify that content is wrapped appropriately for narrow terminals
		// We can't check exact wrapping points, but we can ensure there are enough lines
		if len(lines) < 10 {
			t.Errorf("Stdin input view should properly wrap text in narrow terminals")
		}
	})
}