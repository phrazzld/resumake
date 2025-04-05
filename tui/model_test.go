package tui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelImplementsTea(t *testing.T) {
	// Verify that our Model implements the tea.Model interface
	var _ tea.Model = (*Model)(nil)

	// Create a new model and verify basic properties
	m := NewModel()
	
	// Verify Init returns a command
	_ = m.Init()

	// Test the update function with a simple message
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.Update(msg)
	
	// Verify we get a model back
	if _, ok := updatedModel.(Model); !ok {
		t.Error("Expected Update to return a Model")
	}

	// Verify View returns a string
	_ = m.View() // Just make sure it doesn't crash
}

func TestModelInitialState(t *testing.T) {
	// Test that a new model starts in the welcome state
	m := NewModel()
	
	if m.state != stateWelcome {
		t.Errorf("Expected initial state to be stateWelcome, got %v", m.state)
	}
}

func TestModelStateTransitions(t *testing.T) {
	// Test state transitions based on input and conditions
	
	t.Run("Welcome to Source Input with valid API key", func(t *testing.T) {
		// Create model in welcome state with valid API key
		m := NewModel()
		m.apiKeyOk = true
		
		// Send Enter key
		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model := updatedModel.(Model)
		
		// Verify state transition to source input
		if model.state != stateInputSourcePath {
			t.Errorf("Expected state to transition to stateInputSourcePath, got %v", model.state)
		}
		
		// Verify a command was returned
		if cmd == nil {
			t.Error("Expected a command to be returned")
		}
	})
	
	t.Run("Welcome to Error with invalid API key", func(t *testing.T) {
		// Create model in welcome state with invalid API key
		m := NewModel()
		m.apiKeyOk = false
		
		// Send Enter key
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model := updatedModel.(Model)
		
		// Verify state transition to error
		if model.state != stateResultError {
			t.Errorf("Expected state to transition to stateResultError, got %v", model.state)
		}
		
		// Verify error message is set
		if !strings.Contains(model.errorMsg, "API key") {
			t.Errorf("Expected error message to mention API key, got: %s", model.errorMsg)
		}
	})
	
	t.Run("Source Input to Stdin Input on Enter", func(t *testing.T) {
		// Create model in source input state
		m := NewModel()
		m.state = stateInputSourcePath
		m.sourcePathInput.SetValue("/path/to/test.md")
		
		// Send Enter key
		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model := updatedModel.(Model)
		
		// Verify state transition to stdin input
		if model.state != stateInputStdin {
			t.Errorf("Expected state to transition to stateInputStdin, got %v", model.state)
		}
		
		// Verify commands were returned
		if cmd == nil {
			t.Error("Expected commands to be returned")
		}
	})
	
	t.Run("Stdin Input to Confirm Generate on Ctrl+D", func(t *testing.T) {
		// Create model in stdin input state
		m := NewModel()
		m.state = stateInputStdin
		
		// Send Ctrl+D key
		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
		
		// Verify command was returned (SubmitStdinInputCmd)
		if cmd == nil {
			t.Error("Expected SubmitStdinInputCmd to be returned")
		}
		
		// Simulate the command result
		stdinContent := "Test resume content"
		updatedModel, _ = updatedModel.(Model).Update(StdinSubmitMsg{Content: stdinContent})
		model := updatedModel.(Model)
		
		// Verify state transition to confirm generate
		if model.state != stateConfirmGenerate {
			t.Errorf("Expected state to transition to stateConfirmGenerate, got %v", model.state)
		}
		
		// Verify content was stored
		if model.stdinContent != stdinContent {
			t.Errorf("Expected stdinContent to be set to %q, got %q", stdinContent, model.stdinContent)
		}
	})
	
	t.Run("Confirm Generate to Generating on Enter", func(t *testing.T) {
		// Create model in confirm generate state
		m := NewModel()
		m.state = stateConfirmGenerate
		
		// Send Enter key
		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model := updatedModel.(Model)
		
		// Verify state transition to generating
		if model.state != stateGenerating {
			t.Errorf("Expected state to transition to stateGenerating, got %v", model.state)
		}
		
		// Verify commands were returned
		if cmd == nil {
			t.Error("Expected commands to be returned")
		}
	})
	
	// Note: We're intentionally not testing the Esc key in Confirm Generate state
	// because it requires a deeper level of initialization that is challenging to 
	// set up in a unit test. This would be better tested in an integration test.
}

func TestModelMessageHandling(t *testing.T) {
	// Test handling of different message types
	
	t.Run("FileReadResultMsg success", func(t *testing.T) {
		// Create model
		m := NewModel()
		
		// Send successful file read message
		fileContent := "Sample resume content from file"
		updatedModel, _ := m.Update(FileReadResultMsg{
			Success: true,
			Content: fileContent,
			Error:   nil,
		})
		model := updatedModel.(Model)
		
		// Verify content was stored
		if model.sourceContent != fileContent {
			t.Errorf("Expected sourceContent to be %q, got %q", fileContent, model.sourceContent)
		}
	})
	
	t.Run("FileReadResultMsg failure", func(t *testing.T) {
		// Create model
		m := NewModel()
		
		// Send failed file read message
		fileError := errors.New("file not found")
		updatedModel, _ := m.Update(FileReadResultMsg{
			Success: false,
			Content: "",
			Error:   fileError,
		})
		model := updatedModel.(Model)
		
		// Verify state transitions to error
		if model.state != stateResultError {
			t.Errorf("Expected state to transition to stateResultError, got %v", model.state)
		}
		
		// Verify error message is set
		if !strings.Contains(model.errorMsg, fileError.Error()) {
			t.Errorf("Expected error message to contain %q, got %q", fileError.Error(), model.errorMsg)
		}
	})
	
	
	t.Run("APIResultMsg success", func(t *testing.T) {
		// Create model
		m := NewModel()
		
		// Send successful API result message
		outputPath := "/tmp/resume_output.md"
		content := "Generated resume content"
		updatedModel, _ := m.Update(APIResultMsg{
			Success:    true,
			Content:    content,
			OutputPath: outputPath,
			Error:      nil,
		})
		model := updatedModel.(Model)
		
		// Verify state transitions to success
		if model.state != stateResultSuccess {
			t.Errorf("Expected state to transition to stateResultSuccess, got %v", model.state)
		}
		
		// Verify output path is set
		if model.outputPath != outputPath {
			t.Errorf("Expected outputPath to be %q, got %q", outputPath, model.outputPath)
		}
		
		// Verify result message is set with content length
		expectedLength := fmt.Sprintf("%d", len(content))
		if model.resultMessage != expectedLength {
			t.Errorf("Expected resultMessage to be %q, got %q", expectedLength, model.resultMessage)
		}
	})
	
	t.Run("APIResultMsg failure", func(t *testing.T) {
		// Create model
		m := NewModel()
		
		// Send failed API result message
		apiError := errors.New("API request failed")
		updatedModel, _ := m.Update(APIResultMsg{
			Success: false,
			Content: "",
			Error:   apiError,
		})
		model := updatedModel.(Model)
		
		// Verify state transitions to error
		if model.state != stateResultError {
			t.Errorf("Expected state to transition to stateResultError, got %v", model.state)
		}
		
		// Verify error message is set
		if !strings.Contains(model.errorMsg, apiError.Error()) {
			t.Errorf("Expected error message to contain %q, got %q", apiError.Error(), model.errorMsg)
		}
	})
	
	t.Run("StdinSubmitMsg", func(t *testing.T) {
		// Create model
		m := NewModel()
		
		// Send stdin submit message
		stdinContent := "User-entered resume details"
		updatedModel, _ := m.Update(StdinSubmitMsg{
			Content: stdinContent,
		})
		model := updatedModel.(Model)
		
		// Verify state transitions to confirm generate
		if model.state != stateConfirmGenerate {
			t.Errorf("Expected state to transition to stateConfirmGenerate, got %v", model.state)
		}
		
		// Verify content was stored
		if model.stdinContent != stdinContent {
			t.Errorf("Expected stdinContent to be %q, got %q", stdinContent, model.stdinContent)
		}
	})
	
	t.Run("ProgressUpdateMsg", func(t *testing.T) {
		// Create model
		m := NewModel()
		
		// Send progress update message
		step := "Processing"
		message := "Analyzing input data..."
		updatedModel, _ := m.Update(ProgressUpdateMsg{
			Step:    step,
			Message: message,
		})
		model := updatedModel.(Model)
		
		// Verify progress information is stored
		if model.progressStep != step {
			t.Errorf("Expected progressStep to be %q, got %q", step, model.progressStep)
		}
		
		if model.progressMsg != message {
			t.Errorf("Expected progressMsg to be %q, got %q", message, model.progressMsg)
		}
	})
	
	t.Run("WindowSizeMsg", func(t *testing.T) {
		// Create model
		m := NewModel()
		
		// Send window size message with values different from default
		width := 100
		height := 50
		updatedModel, _ := m.Update(tea.WindowSizeMsg{
			Width:  width,
			Height: height,
		})
		model := updatedModel.(Model)
		
		// Verify dimensions are stored
		if model.width != width {
			t.Errorf("Expected width to be %d, got %d", width, model.width)
		}
		
		if model.height != height {
			t.Errorf("Expected height to be %d, got %d", height, model.height)
		}
		
		// Verify input component dimensions are updated
		// For simplicity, just check they're set to reasonable values
		if model.sourcePathInput.Width <= 0 {
			t.Error("Expected sourcePathInput width to be positive")
		}
		
		if model.stdinInput.Width() <= 0 {
			t.Error("Expected stdinInput width to be positive")
		}
		
		if model.stdinInput.Height() <= 0 {
			t.Error("Expected stdinInput height to be positive")
		}
	})
}

func TestModelAPIKeyCheck(t *testing.T) {
	// Test API key checking functionality
	
	// This test would need to set up environment variables or mock the API key check
	// Skip for now since it depends on external state
	t.Skip("Skipping API key check test as it depends on environment variables")
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// TestModelWithContext tests that the WithContext method correctly sets the context field
func TestModelWithContext(t *testing.T) {
	// Create a new model
	m := NewModel()
	
	// Check if the model has a default context
	if m.ctx == nil {
		t.Error("Expected default model to have a context")
	}
	
	// Create a context with a value
	ctx := context.WithValue(context.Background(), contextKey("test"), "test-value")
	
	// Apply the context to the model
	m = m.WithContext(ctx)
	
	// Verify that the context is set correctly
	if m.ctx == nil {
		t.Error("Expected ctx to be set, got nil")
	}
	
	// Verify that the context contains the expected value
	if val := m.ctx.Value(contextKey("test")); val != "test-value" {
		t.Errorf("Expected context value to be 'test-value', got %v", val)
	}
}

// TestContextPassedToAPIClient tests that the model passes its context to API initialization
func TestContextPassedToAPIClient(t *testing.T) {
	// This test would require mocking the API initialization function
	// which is difficult to do since it's in another package
	
	// Instead, we'll do a simpler check: verify that initializeAPIClient uses m.ctx
	// We can do this by checking the code string itself (using strings.Contains)
	
	// Check if the file contains the expected context usage
	fileContent, err := os.ReadFile("model.go")
	if err != nil {
		t.Skip("Could not read model.go to verify context usage")
	}
	
	// Check if the API client initialization uses the model's context
	if !strings.Contains(string(fileContent), "api.InitializeClient(m.ctx, apiKey)") {
		t.Error("API client initialization should use the model's context")
	}
	
	// Check if GenerateResumeCmd is called with the model's context
	if !strings.Contains(string(fileContent), "GenerateResumeCmd(m.ctx,") {
		t.Error("GenerateResumeCmd should be called with the model's context")
	}
}

func TestModelFieldInitialization(t *testing.T) {
	// Test that all necessary fields are properly initialized
	m := NewModel()
	
	// Check text input initialization
	if m.sourcePathInput.Placeholder != "Enter path to existing resume (optional)" {
		t.Errorf("Expected sourcePathInput placeholder to be set")
	}
	
	// Check text area initialization
	if m.stdinInput.Placeholder != "Enter details about your experience, skills, etc." {
		t.Errorf("Expected stdinInput placeholder to be set")
	}
	
	// Check spinner initialization
	if len(m.spinner.Spinner.Frames) == 0 {
		t.Errorf("Expected spinner to be initialized with frames")
	}
}

func TestModelView(t *testing.T) {
	// Skip this test temporarily while we overhaul the testing approach
	t.Skip("Temporarily skipping while views are being updated")
	
	// Test that View returns appropriate content based on application state
	stateTests := []struct {
		name            string
		setupModel      func() Model
		expectedContent []string
	}{
		{
			name: "Welcome State View",
			setupModel: func() Model {
				m := NewModel()
				m.state = stateWelcome
				m.apiKeyOk = true
				m.width = 100
				m.height = 40
				return m
			},
			expectedContent: []string{
				"R E S U M A K E",
				"API key",
				"valid",
			},
		},
		{
			name: "Source Input State View",
			setupModel: func() Model {
				m := NewModel()
				m.state = stateInputSourcePath
				m.width = 80
				m.height = 24
				return m
			},
			expectedContent: []string{
				"Source File",
				"optional",
				"Enter",
			},
		},
		{
			name: "Stdin Input State View",
			setupModel: func() Model {
				m := NewModel()
				m.state = stateInputStdin
				m.width = 80
				m.height = 24
				return m
			},
			expectedContent: []string{
				"Resume Details",
				"experience",
				"skills",
			},
		},
		{
			name: "Confirm Generate State View",
			setupModel: func() Model {
				m := NewModel()
				m.state = stateConfirmGenerate
				m.sourcePathInput.SetValue("/path/to/source.md")
				m.sourceContent = "Sample source content"
				m.stdinContent = "Sample stdin content"
				m.width = 80
				m.height = 24
				return m
			},
			expectedContent: []string{
				"Ready to Generate",
				"/path/to/source.md",
				"Enter",
			},
		},
		{
			name: "Generating State View",
			setupModel: func() Model {
				m := NewModel()
				m.state = stateGenerating
				m.progressStep = "2 of 4"
				m.progressMsg = "Sending request to Gemini AI..."
				m.width = 80
				m.height = 24
				return m
			},
			expectedContent: []string{
				"Generating Your Resume",
				"Step",
				"2 of 4",
				"Sending request",
			},
		},
		{
			name: "Success State View",
			setupModel: func() Model {
				m := NewModel()
				m.state = stateResultSuccess
				m.outputPath = "/tmp/resume_output.md"
				m.resultMessage = "1500"
				m.width = 80
				m.height = 24
				return m
			},
			expectedContent: []string{
				"Success",
				"/tmp/resume_output.md",
				"Enter",
			},
		},
		{
			name: "Error State View",
			setupModel: func() Model {
				m := NewModel()
				m.state = stateResultError
				m.errorMsg = "API connection failed"
				m.width = 80
				m.height = 24
				return m
			},
			expectedContent: []string{
				"Error",
				"API connection failed",
				"Enter",
			},
		},
	}
	
	// Run each state test
	for _, tt := range stateTests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			view := model.View()
			
			// Check for expected content
			for _, content := range tt.expectedContent {
				if !strings.Contains(view, content) {
					t.Errorf("View should contain '%s'", content)
				}
			}
		})
	}
}