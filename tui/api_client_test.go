package tui

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/generative-ai-go/genai"
)


// TestAPIClientInitialization ensures the API client is initialized only once and at the right moment
// This test was updated after the API client refactoring to verify the new initialization approach
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
	// This test was added as part of the API client refactoring update
	t.Run("API client is initialized only once", func(t *testing.T) {
		// Create a model with valid API key
		m := NewModel()
		m.apiKeyOk = true // Simulate valid API key
		
		// Transition from welcome to sourceInput state to initialize the client
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model := updatedModel.(Model)
		
		// Save the client instance pointer
		originalClient := model.apiClient
		originalModel := model.apiModel
		
		if originalClient == nil {
			t.Fatal("Expected apiClient to be initialized, but it was nil")
		}
		
		// Transition to next state (stdinInput)
		model.state = stateInputSourcePath // Make sure we're in the right state first
		nextModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		modelAfterTransition := nextModel.(Model)
		
		// Verify client instance is the same (not re-initialized)
		if modelAfterTransition.apiClient != originalClient {
			t.Error("Expected apiClient to remain the same instance after state transition")
		}
		
		if modelAfterTransition.apiModel != originalModel {
			t.Error("Expected apiModel to remain the same instance after state transition")
		}
	})
}

// TestExitHandlersCallCleanup ensures different exit paths call cleanupAPIClient
func TestExitHandlersCallCleanup(t *testing.T) {
	// Track if cleanupAPIClient was called
	originalCleanupFunc := cleanupAPIClient
	
	// Test cases for different exit messages
	testCases := []struct {
		name string
		msg  tea.Msg
	}{
		{"QuitMsg", tea.QuitMsg{}},
		{"KeyCtrlC", tea.KeyMsg{Type: tea.KeyCtrlC}},
		{"KeyEsc", tea.KeyMsg{Type: tea.KeyEsc}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a counter to track if cleanup was called
			cleanupCalled := 0
			
			// Override the cleanup function with our instrumented version
			cleanupAPIClient = func(m Model) Model {
				cleanupCalled++
				return m
			}
			
			// Create a new model
			m := NewModel()
			
			// Send the test message to trigger exit logic
			_, _ = m.Update(tc.msg)
			
			// Verify that cleanup was called
			if cleanupCalled == 0 {
				t.Errorf("Expected cleanupAPIClient to be called for %s, but it wasn't", tc.name)
			}
		})
	}
	
	// Test cleanup on Enter key in final states
	finalStates := []struct {
		name  string
		state State
	}{
		{"Success State", stateResultSuccess},
		{"Error State", stateResultError},
	}
	
	for _, fs := range finalStates {
		t.Run("Enter key in "+fs.name, func(t *testing.T) {
			// Create a counter to track if cleanup was called
			cleanupCalled := 0
			
			// Override the cleanup function with our instrumented version
			cleanupAPIClient = func(m Model) Model {
				cleanupCalled++
				return m
			}
			
			// Create a new model in the specified final state
			m := NewModel()
			m.state = fs.state
			
			// Send Enter key message to trigger cleanup
			_, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			
			// Verify that cleanup was called
			if cleanupCalled == 0 {
				t.Errorf("Expected cleanupAPIClient to be called for Enter key in %s, but it wasn't", fs.name)
			}
		})
	}
	
	// Test cleanup after entering the error state due to file reading error
	t.Run("File read error state cleanup", func(t *testing.T) {
		// Create a counter to track if cleanup was called
		cleanupCalled := 0
		
		// Override the cleanup function with our instrumented version
		cleanupAPIClient = func(m Model) Model {
			cleanupCalled++
			return m
		}
		
		// Create a model and directly set it to the error state
		m := NewModel()
		m.state = stateResultError
		m.errorMsg = "Test error"
		
		// Send Enter key to exit from error state
		_, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		
		// Verify that cleanup was called
		if cleanupCalled == 0 {
			t.Error("Expected cleanupAPIClient to be called when exiting from error state")
		}
	})

	// Test cleanup when transitioning to error state after API client was initialized
	t.Run("API error after initialization", func(t *testing.T) {
		// Create a counter to track if cleanup was called
		cleanupCalled := 0
		
		// Override the cleanup function with our instrumented version
		cleanupAPIClient = func(m Model) Model {
			cleanupCalled++
			return m
		}
		
		// Create a model with initialized API client
		m := NewModel()
		// Manually set client (simulating successful initialization)
		m.apiClient = &genai.Client{} // Mock client
		
		// Simulate an API error by sending API failure result
		updatedModel, _ := m.Update(APIResultMsg{
			Success: false,
			Error:   fmt.Errorf("API request failed"),
		})
		// Now model should be in error state
		
		// Send Enter key to exit from error state
		_, _ = updatedModel.(Model).Update(tea.KeyMsg{Type: tea.KeyEnter})
		
		// Verify that cleanup was called
		if cleanupCalled == 0 {
			t.Error("Expected cleanupAPIClient to be called when exiting after API error")
		}
	})

	// Restore the original function
	cleanupAPIClient = originalCleanupFunc
}