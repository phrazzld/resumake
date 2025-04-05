package tui

import (
	"strings"
	"testing"
	
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
)

func TestInputFocusFeedback(t *testing.T) {
	// Test textinput focus feedback
	t.Run("textinput focus feedback", func(t *testing.T) {
		// Create a focused and unfocused input
		focusedInput := textinput.New()
		focusedInput.Focus()
		focusedInput.SetValue("Focused input")
		
		unfocusedInput := textinput.New()
		unfocusedInput.SetValue("Unfocused input")
		
		// Create models with different focus states
		focusedModel := Model{
			sourcePathInput: focusedInput,
			state:           stateInputSourcePath,
			width:           80,
			height:          24,
		}
		
		unfocusedModel := Model{
			sourcePathInput: unfocusedInput,
			state:           stateInputSourcePath,
			width:           80,
			height:          24,
		}
		
		// Render both views
		focusedView := renderSourceFileInputView(focusedModel)
		unfocusedView := renderSourceFileInputView(unfocusedModel)
		
		// The focused view should contain visual indicators that aren't in the unfocused view
		// This is a simple test to ensure there's some difference in the rendered output
		// based on focus state
		if focusedView == unfocusedView {
			t.Error("Focused and unfocused views should have different styling")
		}
	})
	
	// Test textarea focus feedback
	t.Run("textarea focus feedback", func(t *testing.T) {
		// Create a focused and unfocused textarea
		focusedTA := textarea.New()
		focusedTA.Focus()
		focusedTA.SetValue("Focused textarea")
		
		unfocusedTA := textarea.New()
		unfocusedTA.SetValue("Unfocused textarea")
		
		// Create models with different focus states
		focusedModel := Model{
			stdinInput: focusedTA,
			state:      stateInputStdin,
			width:      80,
			height:     24,
		}
		
		unfocusedModel := Model{
			stdinInput: unfocusedTA,
			state:      stateInputStdin,
			width:      80,
			height:     24,
		}
		
		// Render both views
		focusedView := renderStdinInputView(focusedModel)
		unfocusedView := renderStdinInputView(unfocusedModel)
		
		// The focused view should contain visual indicators that aren't in the unfocused view
		if focusedView == unfocusedView {
			t.Error("Focused and unfocused views should have different styling")
		}
	})
	
	// Test for focus indicator elements
	t.Run("focus indicator elements", func(t *testing.T) {
		// Create a focused input
		focusedInput := textinput.New()
		focusedInput.Focus()
		focusedInput.SetValue("Test input")
		
		// Create a model with focused input
		model := Model{
			sourcePathInput: focusedInput,
			state:           stateInputSourcePath,
			width:           80,
			height:          24,
		}
		
		// Render the view
		view := renderSourceFileInputView(model)
		
		// Check for visual focus indicators
		// The focused input should have a distinctive styling that includes at least
		// one of these focus indicator terms
		focusIndicators := []string{"focused", "active", "selected", "highlight"}
		
		foundIndicator := false
		for _, indicator := range focusIndicators {
			if strings.Contains(strings.ToLower(view), indicator) {
				foundIndicator = true
				break
			}
		}
		
		if !foundIndicator {
			t.Error("Focused input should have a visual indicator in the view")
		}
	})
}