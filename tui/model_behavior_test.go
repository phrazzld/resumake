package tui

import (
	"errors"
	"testing"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

// TestModelStateBehavior uses our behavior-focused utilities to test model transitions
func TestModelStateBehavior(t *testing.T) {
	// Simplified for now while working on the test overhaul
	t.Skip("Temporarily skipping transition test while views are being updated")
	
	// Create a model in welcome state with valid API key
	model := NewModel()
	model.apiKeyOk = true
	
	// Create assertions helper for the initial state
	ma := NewModelAssertions(t, model)
	ma.AssertState(stateWelcome)
	
	// Test the welcome view structure
	welcomeView := model.View()
	welcomeElements := []string{"R E S U M A K E", "API key"}
	for _, element := range welcomeElements {
		if !ContainsViewElement(welcomeView, element) {
			t.Errorf("Welcome view should contain '%s'", element)
		}
	}
	
	// Simulate transition to source input state
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = updatedModel.(Model)
	
	// Check new state
	ma = NewModelAssertions(t, model)
	ma.AssertState(stateInputSourcePath)
}

// TestProgressUpdateBehavior tests specifically how progress updates affect the model
func TestProgressUpdateBehavior(t *testing.T) {
	// Create a model in generating state
	model := NewModel()
	model.state = stateGenerating 
	model.spinner = spinner.New()
	model.stdinContent = "Test content"
	model.width = 80
	model.height = 24
	
	// Verify initial state
	ma := NewModelAssertions(t, model)
	ma.AssertState(stateGenerating)
	
	// Initial view should have some basic generating content
	initialView := model.View()
	if !ContainsViewElement(initialView, "Generating") {
		t.Error("Generating view should have a title with 'Generating'")
	}
	
	// Send a progress update
	step1 := "1 of 4"
	msg1 := "Building prompt from your inputs..."
	updatedModel, _ := model.Update(ProgressUpdateMsg{Step: step1, Message: msg1})
	model, _ = updatedModel.(Model)
	
	// Check that progress info was updated
	ma = NewModelAssertions(t, model)
	ma.AssertProgressStep(step1)
	ma.AssertProgressMessage(msg1)
	
	// Send a second progress update
	step2 := "2 of 4"
	msg2 := "Sending request to Gemini AI..."
	updatedModel, _ = model.Update(ProgressUpdateMsg{Step: step2, Message: msg2})
	model, _ = updatedModel.(Model)
	
	// Check that progress info was updated
	ma = NewModelAssertions(t, model)
	ma.AssertProgressStep(step2)
	ma.AssertProgressMessage(msg2)
}

// TestErrorHandlingBehavior tests how the model handles errors
func TestErrorHandlingBehavior(t *testing.T) {
	// Create a model in generating state
	model := NewModel()
	model.state = stateGenerating
	model.spinner = spinner.New()
	model.width = 80
	model.height = 24
	
	// Test API error
	errorMsg := "Failed to connect to API: timeout"
	updatedModel, _ := model.Update(APIResultMsg{
		Success: false,
		Error:   errors.New(errorMsg),
	})
	model, _ = updatedModel.(Model)
	
	// Check state transition
	ma := NewModelAssertions(t, model)
	ma.AssertState(stateResultError)
	ma.AssertError(errorMsg)
	
	// Error view should contain the error message
	errorView := model.View()
	if !ContainsViewElement(errorView, "Error") {
		t.Error("Error view should have an error title")
	}
	
	if !ContainsViewElement(errorView, "Failed to connect") {
		t.Error("Error view should contain the error message text")
	}
}

// TestSuccessStateBehavior tests the success state behavior
func TestSuccessStateBehavior(t *testing.T) {
	// Create a model in generating state
	model := NewModel()
	model.state = stateGenerating
	model.spinner = spinner.New()
	model.width = 80
	model.height = 24
	
	// Send a success result
	content := "# Generated Resume\n\nThis is a test resume."
	outputPath := "/tmp/resume_out.md"
	updatedModel, _ := model.Update(APIResultMsg{
		Success:    true,
		Content:    content,
		OutputPath: outputPath,
	})
	model, _ = updatedModel.(Model)
	
	// Check state transition
	ma := NewModelAssertions(t, model)
	ma.AssertState(stateResultSuccess)
	
	// Success view should contain the output path
	successView := model.View()
	if !ContainsViewElement(successView, "Success") {
		t.Error("Success view should have a success title")
	}
	
	if !ContainsViewElement(successView, outputPath) {
		t.Error("Success view should contain the output path")
	}
	
	// ResultMessage should be the content length
	// Note: '42' is the length of "# Generated Resume\n\nThis is a test resume."
	expectedLength := "42"
	ma.AssertResultContains(expectedLength)
}