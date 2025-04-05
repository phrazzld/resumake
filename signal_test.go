package main

import (
	"context"
	"testing"
	
	"github.com/phrazzld/resumake/tui"
)

// TestSignalHandlingFunctionExists tests that the signal handling function exists
func TestSignalHandlingFunctionExists(t *testing.T) {
	// Test that setupProgramWithSignalHandling function exists by calling it
	model := tui.NewModel()
	
	// Create a no-op cancel function for testing
	noopCancel := func() {}
	
	// Call the function with the required parameters
	program := setupProgramWithSignalHandling(model, noopCancel)
	
	// Verify the function returned a non-nil program
	if program == nil {
		t.Error("setupProgramWithSignalHandling returned nil program")
	}
}

// TestSignalHandlingWithCancel tests that setupProgramWithSignalHandling accepts a cancel function
func TestSignalHandlingWithCancel(t *testing.T) {
	// Create a model
	model := tui.NewModel()
	
	// Create a mock cancel function
	mockCancel := context.CancelFunc(func() {
		// This function won't be called in the test, but could be in actual use
		// We just need to verify setupProgramWithSignalHandling accepts it
	})
	
	// Call setupProgramWithSignalHandling with the cancel function
	p := setupProgramWithSignalHandling(model, mockCancel)
	
	// Basic check that the function returns a non-nil program
	if p == nil {
		t.Error("Expected setupProgramWithSignalHandling to return a non-nil program when given a cancel function")
	}
}