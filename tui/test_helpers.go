package tui

import (
	"testing"
	"strings"
	
	tea "github.com/charmbracelet/bubbletea"
)

// ModelAssertions provides helper methods to test the Model's behaviors and structure
// rather than exact string outputs
type ModelAssertions struct {
	Model *Model
	T     *testing.T
}

// NewModelAssertions creates a new assertion helper for a model
func NewModelAssertions(t *testing.T, m Model) *ModelAssertions {
	return &ModelAssertions{
		Model: &m,
		T:     t,
	}
}

// AssertView renders the view and asserts general structure properties
// rather than exact string matching
func (ma *ModelAssertions) AssertView(expectedElements []string, unexpectedElements []string) {
	view := ma.Model.View()
	
	// Check for expected elements
	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			ma.T.Errorf("View should contain '%s' but doesn't", element)
		}
	}
	
	// Check for unexpected elements
	for _, element := range unexpectedElements {
		if strings.Contains(view, element) {
			ma.T.Errorf("View should not contain '%s' but does", element)
		}
	}
}

// AssertState checks if the model state matches expected value
func (ma *ModelAssertions) AssertState(expectedState State) {
	if ma.Model.state != expectedState {
		ma.T.Errorf("Model state expected to be %v but was %v", expectedState, ma.Model.state)
	}
}

// AssertProgressStep ensures the progress information contains expected step information
func (ma *ModelAssertions) AssertProgressStep(expectedStep string) {
	if !strings.Contains(ma.Model.progressStep, expectedStep) {
		ma.T.Errorf("Progress step expected to contain '%s' but was '%s'", expectedStep, ma.Model.progressStep)
	}
}

// AssertProgressMessage ensures progress message contains expected information
func (ma *ModelAssertions) AssertProgressMessage(expectedMessageContent string) {
	if !strings.Contains(ma.Model.progressMsg, expectedMessageContent) {
		ma.T.Errorf("Progress message expected to contain '%s' but was '%s'", expectedMessageContent, ma.Model.progressMsg)
	}
}

// AssertNoError checks that the model isn't in an error state
func (ma *ModelAssertions) AssertNoError() {
	if ma.Model.errorMsg != "" {
		ma.T.Errorf("Model should not have error but has: %s", ma.Model.errorMsg)
	}
}

// AssertError checks that the model is in an error state containing expected message
func (ma *ModelAssertions) AssertError(expectedErrorContent string) {
	if !strings.Contains(ma.Model.errorMsg, expectedErrorContent) {
		ma.T.Errorf("Error message expected to contain '%s' but was '%s'", expectedErrorContent, ma.Model.errorMsg)
	}
}

// AssertResultContains checks that the result message contains expected content
func (ma *ModelAssertions) AssertResultContains(expectedContent string) {
	if !strings.Contains(ma.Model.resultMessage, expectedContent) {
		ma.T.Errorf("Result message expected to contain '%s' but was '%s'", expectedContent, ma.Model.resultMessage)
	}
}

// SimulatedInput helps test a sequence of inputs
type SimulatedInput struct {
	Model  Model
	T      *testing.T
	Inputs []tea.Msg
}

// NewSimulatedInput creates a test helper to simulate input sequences
func NewSimulatedInput(t *testing.T, initialModel Model) *SimulatedInput {
	return &SimulatedInput{
		Model:  initialModel,
		T:      t,
		Inputs: []tea.Msg{},
	}
}

// AddKeyPress adds a key press to the input sequence
func (si *SimulatedInput) AddKeyPress(key string) *SimulatedInput {
	// Add the key press message to the sequence
	// This will depend on the exact key message structure used in the app
	return si
}

// AddProgressUpdate adds a progress update to the input sequence
func (si *SimulatedInput) AddProgressUpdate(step, message string) *SimulatedInput {
	si.Inputs = append(si.Inputs, ProgressUpdateMsg{Step: step, Message: message})
	return si
}

// AddAPIResult adds an API result to the input sequence
func (si *SimulatedInput) AddAPIResult(content string, outputPath string, err error) *SimulatedInput {
	success := err == nil
	si.Inputs = append(si.Inputs, APIResultMsg{
		Success:    success,
		Content:    content,
		OutputPath: outputPath,
		Error:      err,
	})
	return si
}

// RunSequence runs the entire input sequence and returns the final model
func (si *SimulatedInput) RunSequence() Model {
	model := si.Model
	
	// Process each input and update the model
	for _, input := range si.Inputs {
		var cmd tea.Cmd
		updatedModel, cmd := model.Update(input)
		model, _ = updatedModel.(Model)
		
		// Handle commands that return additional messages
		if cmd != nil {
			msg := cmd()
			if msg != nil {
				updatedModel, _ := model.Update(msg)
				model, _ = updatedModel.(Model)
			}
		}
	}
	
	return model
}

// ContainsViewElement checks if the view contains an expected string element
func ContainsViewElement(view, element string) bool {
	return strings.Contains(view, element)
}

// ViewHasAllElements checks if the view contains all expected elements
func ViewHasAllElements(view string, elements []string) bool {
	for _, element := range elements {
		if !strings.Contains(view, element) {
			return false
		}
	}
	return true
}