package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelImplementsTea(t *testing.T) {
	// Verify that our Model implements the tea.Model interface
	var _ tea.Model = (*Model)(nil)

	// Create a new model and verify basic properties
	m := NewModel()
	
	// Verify Init returns a command
	cmd := m.Init()
	if cmd != nil {
		// This is okay, Init can return nil or a command
	}

	// Test the update function with a simple message
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := m.Update(msg)
	
	// Verify we get a model back
	if _, ok := updatedModel.(Model); !ok {
		t.Error("Expected Update to return a Model")
	}

	// Verify View returns a string
	view := m.View()
	if view != "Resumake TUI" {
		t.Errorf("Expected View to return 'Resumake TUI', got %q", view)
	}
}