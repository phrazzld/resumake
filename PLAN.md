```markdown
# Technical Plan: TUI Testing Overhaul

## 1. Overview

This plan outlines the strategy to overhaul the Terminal User Interface (TUI) testing for the Resumake Go application. The current tests, particularly those interacting with the Bubble Tea `Model`'s `View()` method, rely heavily on `strings.Contains` checks against rendered output. This makes them brittle and prone to breaking with minor UI text or styling changes.

The goal is to shift towards a more resilient testing approach that focuses on:

1.  **Behavior:** Verifying state transitions, command emissions, and data flow within the TUI model in response to user inputs and events.
2.  **Structure:** Asserting the presence and state of key UI elements (like input fields, spinners, error messages) based on the model's state, rather than checking exact rendered text.

This will involve refactoring existing TUI tests (`tui/*_test.go`) and potentially introducing helper methods on the `tui.Model` to facilitate these structural assertions.

## 2. Task Breakdown

| Task                                                                 | Description                                                                                                                                                              | Effort | Affected Files/Modules                     |
| :------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :----- | :----------------------------------------- |
| **1. Analyze Existing TUI Tests**                                    | Deep dive into `tui/model_test.go` and `tui/views_test.go` to pinpoint all instances of brittle string/component checking on `View()` output.                             | S      | `tui/model_test.go`, `tui/views_test.go`   |
| **2. Define Structural/Behavioral Assertions**                       | Define a set of assertions needed to verify TUI state without inspecting raw view strings (e.g., `IsShowingError`, `IsInputFocused`, `GetCurrentPrompt`).               | M      | Design Document / This Plan              |
| **3. Implement Model Helper Methods**                                | Add methods to `tui.Model` (or a test-only extension) that expose structural/state information based on the current `model.state` and component states.                | M      | `tui/model.go` (or new test helper file) |
| **4. Refactor `views_test.go`**                                      | Remove tests that directly check `renderXView` output strings. These views are tested implicitly via the model state tests. Keep tests for utility functions if any. | M      | `tui/views_test.go`, `tui/utils_test.go`   |
| **5. Refactor `model_test.go` - State & Command Assertions**         | Enhance tests to rigorously check `model.state`, relevant data fields (`errorMsg`, `sourceContent`, etc.), and emitted `tea.Cmd`s after `Update` calls.              | L      | `tui/model_test.go`                      |
| **6. Refactor `model_test.go` - Structural Assertions**              | Replace `strings.Contains(view, ...)` checks with calls to the new model helper methods implemented in Task 3.                                                            | L      | `tui/model_test.go`                      |
| **7. Implement Input Sequence Test Helpers (Optional but Recommended)** | Create helper functions to simulate sequences of user inputs (key presses, messages) and assert the final model state or structural properties.                       | M      | New test helper file (e.g., `tui/test_helpers_test.go`) |
| **8. Review & Refine Tests**                                         | Review all refactored TUI tests for clarity, coverage, and resilience. Ensure they focus on behavior and structure.                                                    | M      | `tui/*_test.go`                          |
| **9. Documentation Update**                                          | Briefly document the new testing approach and the purpose of any helper methods/functions.                                                                               | S      | README.md or CONTRIBUTING.md           |

**Effort Estimation:** S = Small (<= 1 day), M = Medium (1-3 days), L = Large (3-5 days)

## 3. Implementation Details

### 3.1. Problem Example (Current Brittle Test)

From `tui/views_test.go` or similar checks in `tui/model_test.go`:

```go
// In tui/views_test.go or tui/model_test.go
func TestRenderWelcomeView(t *testing.T) {
    model := Model{apiKeyOk: true, /* ... */}
    view := renderWelcomeView(model) // Or m.View() in model_test.go

    // BRITTLE CHECK: Relies on exact string
    if !strings.Contains(view, "API KEY STATUS: READY") {
        t.Error("Welcome view should indicate API key is valid")
    }
    // BRITTLE CHECK: Relies on exact string
    if !strings.Contains(view, "Press Enter to begin") {
        t.Error("Welcome view should include instructions to proceed")
    }
}
```

### 3.2. Proposed Solution: Model Helper Methods

Introduce methods on `tui.Model` (or potentially in a separate test utility package accessing the model) to query its state structurally.

**Affected File:** `tui/model.go` (or new `tui/model_test_helpers.go`)

```go
// Example Helper Methods in tui/model.go (or similar)

// IsShowingWelcome returns true if the model is in the welcome state.
func (m Model) IsShowingWelcome() bool {
	return m.state == stateWelcome
}

// HasValidAPIKeyIndication returns true if the model state reflects a valid API key check.
// This focuses on the *state* driving the view, not the view itself.
func (m Model) HasValidAPIKeyIndication() bool {
    // Assumes apiKeyOk is the source of truth for the view's indication
	return m.state == stateWelcome && m.apiKeyOk
}

// IsShowingAPIKeyError returns true if the model is showing an error related to the API key.
func (m Model) IsShowingAPIKeyError() bool {
    return (m.state == stateResultError && strings.Contains(m.errorMsg, "API key")) ||
           (m.state == stateWelcome && !m.apiKeyOk) // Covers the initial check failure message
}

// IsSourceInputFocused returns true if the source path input component should be focused.
func (m Model) IsSourceInputFocused() bool {
    // Check the state and potentially the focus state of the component
    return m.state == stateInputSourcePath // Simplified example
    // A more robust check might involve checking m.sourcePathInput.Focused()
    // if the focus logic is complex, but often state is sufficient.
}

// GetDisplayedError returns the error message intended for display, if any.
func (m Model) GetDisplayedError() (string, bool) {
    if m.state == stateResultError && m.errorMsg != "" {
        return m.errorMsg, true
    }
    return "", false
}

// IsShowingSpinner returns true if the spinner should be visible.
func (m Model) IsShowingSpinner() bool {
    return m.state == stateGenerating
}

// GetCurrentProgressStep returns the current progress step text.
func (m Model) GetCurrentProgressStep() (string, bool) {
    if m.state == stateGenerating && m.progressStep != "" {
        return m.progressStep, true
    }
    return "", false
}
```

### 3.3. Refactored Test Example

Replace brittle string checks with calls to the new helper methods.

**Affected File:** `tui/model_test.go`

```go
// In tui/model_test.go
func TestModelStateTransitions(t *testing.T) {
    t.Run("Welcome to Source Input with valid API key", func(t *testing.T) {
        m := NewModel()
        m.apiKeyOk = true // Simulate valid API key state

        // Send Enter key
        updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
        model := updatedModel.(Model)

        // Assert state transition
        if model.state != stateInputSourcePath {
            t.Errorf("Expected state stateInputSourcePath, got %v", model.state)
        }
        // Assert structural property using helper
        if !model.IsSourceInputFocused() { // Example helper usage
             t.Error("Expected source input to be focused")
        }
        // Assert command emission
        if cmd == nil { // Or check for specific command type if needed
            t.Error("Expected a command to be returned")
        }
        // Assert NO error is displayed
        if _, isError := model.GetDisplayedError(); isError {
            t.Error("Expected no error to be displayed")
        }
    })

    t.Run("Welcome shows API Key Error", func(t *testing.T) {
		m := NewModel()
		m.apiKeyOk = false // Simulate invalid API key state

		// Send Enter key (which triggers the error state transition in this case)
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model := updatedModel.(Model)

        // Assert state transition
		if model.state != stateResultError {
			t.Errorf("Expected state stateResultError, got %v", model.state)
		}
        // Assert structural property using helper
		if !model.IsShowingAPIKeyError() {
			t.Error("Expected model to indicate an API key error")
		}
        // Assert specific error message content (less brittle than checking full view)
        errMsg, ok := model.GetDisplayedError()
        if !ok || !strings.Contains(errMsg, "API key") {
             t.Errorf("Expected error message about API key, got: %q", errMsg)
        }
	})
}
```

### 3.4. Input Sequence Testing Helper (Conceptual)

**Affected File:** New test helper file (e.g., `tui/test_helpers_test.go`)

```go
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

// Simulate sends a sequence of messages to a model and returns the final model.
func Simulate(t *testing.T, initialModel tea.Model, messages ...tea.Msg) tea.Model {
	t.Helper()
	currentModel := initialModel
	var cmd tea.Cmd
	for _, msg := range messages {
		currentModel, cmd = currentModel.Update(msg)
		// In a real scenario, you might need to handle commands
		// by generating their corresponding result messages and feeding them back.
		// For simplicity here, we ignore commands, but a full implementation
		// would need to process them (e.g., if cmd produces FileReadResultMsg).
		if cmd != nil {
            // Basic example: If a command is returned, execute it and feed the result back
            // This needs careful implementation based on command types.
            // resultMsg := cmd() // This blocks, real implementation needs care
            // currentModel, cmd = currentModel.Update(resultMsg)
            t.Logf("Note: Command returned but not processed in this simple Simulate helper: %T", cmd)
		}
	}
	return currentModel
}

// Example Usage in a test:
func TestFullUserInputFlow(t *testing.T) {
    initialModel := NewModel()
    initialModel.apiKeyOk = true // Assume valid key

    // Define the sequence of inputs
    inputs := []tea.Msg{
        tea.KeyMsg{Type: tea.KeyEnter},         // Welcome -> Source Input
        // Simulate typing a path (can send individual KeyRunes or set value directly for simplicity)
        // tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/path/to/file.md")},
        tea.KeyMsg{Type: tea.KeyEnter},         // Source Input -> Stdin Input (triggers FileReadCmd)
        // Simulate FileReadResultMsg (assuming file read succeeds)
        FileReadResultMsg{Success: true, Content: "Existing data"},
        // Simulate typing in textarea
        // tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("New input")},
        tea.KeyMsg{Type: tea.KeyCtrlD},         // Stdin Input -> Confirm Generate (triggers StdinSubmitMsg)
        // Simulate StdinSubmitMsg
        StdinSubmitMsg{Content: "New input"},
        tea.KeyMsg{Type: tea.KeyEnter},         // Confirm Generate -> Generating (triggers GenerateResumeCmd)
        // Simulate APIResultMsg (assuming success)
        APIResultMsg{Success: true, Content: "Generated resume", OutputPath: "out.md"},
    }

    // Simulate the flow
    finalModelInterface := Simulate(t, initialModel, inputs...)
    finalModel := finalModelInterface.(Model) // Type assertion

    // Assert final state
    if finalModel.state != stateResultSuccess {
        t.Errorf("Expected final state stateResultSuccess, got %v", finalModel.state)
    }
    if finalModel.outputPath != "out.md" {
        t.Errorf("Expected output path 'out.md', got %q", finalModel.outputPath)
    }
    // Add more assertions using helpers
    if _, isError := finalModel.GetDisplayedError(); isError {
        t.Error("Expected no final error")
    }
}

```
*Note: The `Simulate` helper needs careful implementation to handle commands correctly, potentially by executing them and sending their resulting messages back into the `Update` loop.*

## 4. Potential Challenges & Considerations

1.  **Complexity of Helpers:** Helper methods on the model must accurately reflect the logic driving the view without becoming overly complex or tightly coupled to specific view rendering details.
2.  **Asynchronous Commands:** Testing flows involving asynchronous commands (`tea.Cmd`) requires careful handling. The `Simulate` helper needs to potentially execute commands and feed their resulting `tea.Msg` back into the `Update` loop, which can be complex to implement robustly for all command types.
3.  **Focus Management:** Accurately testing which component has focus might require exposing more state from the underlying Bubble components or relying on state-machine logic (e.g., "in state X, component Y *should* be focused").
4.  **Refactoring Effort:** Refactoring existing tests, especially in `model_test.go`, will require careful work to replace view checks with state/structural checks without losing coverage.
5.  **Maintaining Helpers:** As the TUI evolves, the helper methods must be maintained alongside the model and view logic.

## 5. Testing Strategy

1.  **Unit Tests (Primary Focus):**
    *   Test individual `tui.Model` helper methods thoroughly.
    *   Test the `tui.Model.Update` function extensively:
        *   Verify correct state transitions for various `tea.Msg` inputs in each state.
        *   Verify correct data updates within the model (`errorMsg`, `sourceContent`, `progressStep`, etc.).
        *   Verify correct `tea.Cmd` emissions for relevant state transitions.
        *   Use the new helper methods to assert structural properties post-update (e.g., `IsShowingError`, `IsSpinnerVisible`).
    *   Continue unit testing commands (`tui/commands_test.go`) using mocks/stubs where necessary (especially for API/file interactions not covered by dry runs).
    *   Remove direct testing of `renderXView` functions (`tui/views_test.go`) as they are implicitly tested via the model state tests. Keep tests for utility functions like `wrapText`.
2.  **Integration Tests (Via Input Sequences):**
    *   Use the `Simulate` helper (or similar) to test common user flows (e.g., welcome -> generate -> success, welcome -> file error, welcome -> API error).
    *   Assert the final `Model` state and key structural properties after the sequence.
3.  **Manual Testing (Reduced Scope):**
    *   Perform minimal manual checks after significant changes to ensure the overall look and feel remain correct, but rely primarily on automated tests for functional correctness.

## 6. Open Questions

1.  **Definition of "Structure":** How granular should the structural helper methods be? Is checking for the *presence* of an error message sufficient, or should helpers sometimes return specific error *types* or key phrases if they represent distinct states? (Initial proposal: Focus on presence/absence of elements like errors, spinners, specific inputs, and key state indicators like API key validity).
2.  **Command Simulation:** How comprehensively should the `Simulate` helper handle `tea.Cmd` execution? Should it fully mock file reads, API calls, etc., or focus only on commands generating simple messages? (Initial proposal: Start simple, handling only commands that return basic messages like `StdinSubmitMsg`, and expand if needed).
3.  **Helper Location:** Should helper methods be directly on `tui.Model` or in a separate test-only package/file to avoid polluting the main model? (Initial proposal: Add directly to `tui.Model` for simplicity, reconsider if it becomes cluttered).
```