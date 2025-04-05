package tui

import (
	"testing"
	
	tea "github.com/charmbracelet/bubbletea"
)

// TestSpinnerTicking tests that the spinner ticking command is properly added
// during the generating state
func TestSpinnerTicking(t *testing.T) {
	// Create a new model in generating state
	m := NewModel()
	m.state = stateGenerating
	
	// Create a mock spinner message
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	
	// Update the model with the message
	newModel, cmd := m.Update(msg)
	
	// Verify that a command was returned (should include spinner.Tick)
	if cmd == nil {
		t.Fatal("Expected a command to be returned, got nil")
	}
	
	// Convert the model to the correct type
	updatedModel, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Expected model of type Model, got %T", newModel)
	}
	
	// Verify that the model is still in generating state
	if updatedModel.state != stateGenerating {
		t.Errorf("Expected model to remain in stateGenerating, got %v", updatedModel.state)
	}
	
	// This test will fail because we need to ensure the spinner command is always included
	// We'll fix this in the implementation
	
	// TODO: Enhance this test to verify the specific spinner command is included
}

// TestSpinnerStateTransitions tests that the spinner animation is properly
// handled during state transitions
func TestSpinnerStateTransitions(t *testing.T) {
	// Create a new model in a non-generating state
	m := NewModel()
	m.state = stateConfirmGenerate
	
	// Transition to generating state
	m.state = stateGenerating
	
	// Update model with a message that doesn't trigger specific case handlers
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	newModel, cmd := m.Update(msg)
	
	// Verify command was returned
	if cmd == nil {
		t.Fatal("Expected command to be returned for generating state, got nil")
	}
	
	// Transition to result success state
	updatedModel, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Expected model of type Model, got %T", newModel)
	}
	
	// This transition should stop spinner animation
	updatedModel.state = stateResultSuccess
	finalModel, finalCmd := updatedModel.Update(msg)
	
	// For non-generating states with no specific handlers,
	// we expect nil command when no other commands are queued
	if updatedModel.state != stateGenerating && finalCmd != nil {
		t.Errorf("Expected nil command after transition from generating state, got non-nil")
	}
	
	// Verify model type
	_, ok = finalModel.(Model)
	if !ok {
		t.Fatalf("Expected final model of type Model, got %T", finalModel)
	}
}

// TestSpinnerCleanup tests that the spinner is properly cleaned up when
// exiting the generating state
func TestSpinnerCleanup(t *testing.T) {
	// Create a new model in generating state
	m := NewModel()
	m.state = stateGenerating
	
	// Save initial spinner frames to check later
	initialFrames := m.spinner.Spinner.Frames
	
	// Create a mock API result message that will transition state
	msg := APIResultMsg{
		Success: true,
		Content: "Test content",
		OutputPath: "test.md",
	}
	
	// Update the model with this message (should transition to success state)
	newModel, _ := m.Update(msg)
	
	// Verify state transition
	updatedModel, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Expected model of type Model, got %T", newModel)
	}
	
	// Check that we transitioned away from generating state
	if updatedModel.state != stateResultSuccess {
		t.Errorf("Expected transition to stateResultSuccess, got %v", updatedModel.state)
	}
	
	// Verify spinner frames remain intact for potential future use
	if len(updatedModel.spinner.Spinner.Frames) != len(initialFrames) {
		t.Errorf("Spinner frames changed unexpectedly during state transition")
	}
}