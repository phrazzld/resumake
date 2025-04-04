package main

import (
	"testing"
	
	"github.com/phrazzld/resumake/tui"
)

// TestSignalHandlingFunctionExists tests that the signal handling function exists
func TestSignalHandlingFunctionExists(t *testing.T) {
	// Test that setupProgramWithSignalHandling function exists by calling it
	model := tui.NewModel()
	program := setupProgramWithSignalHandling(model)
	
	// Verify the function returned a non-nil program
	if program == nil {
		t.Error("setupProgramWithSignalHandling returned nil program")
	}
}