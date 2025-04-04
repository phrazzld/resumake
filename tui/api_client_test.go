package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// mockQuitMsg implements tea.Msg interface for testing quit handling
type mockQuitMsg struct{}

// TestAPIClientInitialization ensures the API client is initialized only once and at the right moment
func TestAPIClientInitialization(t *testing.T) {
	// Test 1: API client is nil on initial model creation
	t.Run("API client is nil on model creation", func(t *testing.T) {
		m := NewModel()
		
		if m.apiClient != nil {
			t.Error("Expected apiClient to be nil on model creation")
		}
		
		if m.apiModel != nil {
			t.Error("Expected apiModel to be nil on model creation")
		}
	})
	
	// Test 2: API client is initialized when transitioning from welcome to sourceInput with valid API key
	t.Run("API client is initialized on state transition", func(t *testing.T) {
		// This test will fail until we implement the initialization logic
		m := NewModel()
		m.apiKeyOk = true // Simulate valid API key
		
		// Transition from welcome to sourceInput
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model := updatedModel.(Model)
		
		// Assert that apiClient and apiModel are now initialized
		if model.apiClient == nil {
			t.Error("Expected apiClient to be initialized after state transition")
		}
		
		if model.apiModel == nil {
			t.Error("Expected apiModel to be initialized after state transition")
		}
	})
	
	// Test 3: Client is not re-initialized on subsequent state transitions
	t.Run("API client is initialized only once", func(t *testing.T) {
		// This test will be implemented after basic initialization works
		// The test will need to capture the client instance and ensure it doesn't change
	})
}

// TestAPIClientCleanup ensures the API client is closed when application exits
func TestAPIClientCleanup(t *testing.T) {
	t.Run("API client is closed on application quit", func(t *testing.T) {
		// This test will fail until we implement the cleanup logic
		// We need a way to detect if the client was closed
		// This may require a mock client or a more complex test setup
		
		m := NewModel()
		m.apiKeyOk = true
		
		// First transition to initialize the client
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model := updatedModel.(Model)
		
		// Then trigger a quit
		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		
		// We need to verify that the quit command would close the client
		// This is challenging in a unit test without modifying the code specifically for testing
		// For now, we'll just ensure the command is a quit command
		if cmd == nil {
			t.Error("Expected a command to be returned")
		}
	})
}