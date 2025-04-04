```markdown
# PLAN.md: Implement Bubble Tea TUI for Resumake

## 1. Overview

This plan outlines the steps to replace the current command-line interface (CLI) interaction of `resumake` with a lightweight Terminal User Interface (TUI) using the `bubbletea` library. The goal is to provide a more interactive and user-friendly experience for guiding the user through the resume generation process, displaying status updates, and handling inputs/outputs within the terminal window. The core logic for API interaction, prompt generation, and file handling will be retained but integrated into the `bubbletea` application lifecycle.

## 2. Task Breakdown

| Task                                                    | Description                                                                                                                               | Effort | Affected Files/Modules                     |
| :------------------------------------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------- | :----- | :----------------------------------------- |
| **1. Add Dependencies**                                 | Add `github.com/charmbracelet/bubbletea` and potentially related `charmbracelet` libraries (e.g., `lipgloss`, `bubbles`) to `go.mod`.        | S      | `go.mod`, `go.sum`                         |
| **2. Design TUI States & Model**                        | Define the different states of the TUI (e.g., Welcome, InputSource, InputStdin, Generating, Success, Error) and design the main `bubbletea` model struct to hold the application state. | M      | `tui/model.go` (new)                       |
| **3. Create Initial TUI Structure**                     | Refactor `main.go` to initialize and run the `bubbletea` program (`tea.NewProgram`). Create basic `Init`, `Update`, and `View` methods for the main model. | M      | `main.go`, `tui/model.go`                  |
| **4. Implement Welcome/Initial View**                   | Create the initial view displayed when the application starts. Check for API Key and display status/instructions.                         | S      | `tui/model.go`, `tui/view_welcome.go` (new) |
| **5. Refactor Input Handling - Source File**            | Integrate source file path input into the TUI. Replace direct flag reading for this purpose within the TUI flow. Use a `textinput` bubble. | M      | `tui/model.go`, `tui/cmd_input.go` (new), `input/file.go` |
| **6. Refactor Input Handling - Stdin**                  | Replace `input.ReadFromStdin` with a `textarea` bubble for multi-line text input within the TUI.                                          | M      | `tui/model.go`, `tui/cmd_input.go`, `input/stdin.go` |
| **7. Implement Generation State View & Logic**          | Create a view to show progress (e.g., spinner, status messages) during API call and processing.                                           | M      | `tui/model.go`, `tui/view_generating.go` (new), `tui/cmd_api.go` (new) |
| **8. Refactor Core Logic into Commands**                | Wrap existing logic (reading files, building prompt, calling API, processing response, writing output) into `tea.Cmd` functions that return `tea.Msg`. | L      | `tui/cmd_*.go`, `api/`, `prompt/`, `output/` |
| **9. Implement Message Handling in Update**             | Extend the `Update` method to handle messages returned by commands (e.g., file read result, API response, errors). Update the model state accordingly. | L      | `tui/model.go`, `tui/messages.go` (new)    |
| **10. Implement Success/Error Views**                   | Create views to display the final success message (output path, etc.) or detailed error messages if failures occur.                       | M      | `tui/model.go`, `tui/view_result.go` (new) |
| **11. Integrate Styling (lipgloss)**                    | Apply basic styling using `lipgloss` for better readability and presentation.                                                               | S      | `tui/view_*.go`, `tui/styles.go` (new)     |
| **12. Refactor/Update Tests**                           | Update existing tests (`main_test.go` might need significant changes or removal). Add unit tests for the TUI model's `Update` logic and command functions. | L      | `main_test.go`, `tui/*_test.go` (new)      |
| **13. Manual Testing & Refinement**                     | Perform thorough manual testing across different terminal environments and refine the UX based on feedback.                               | M      | N/A                                        |
| **14. (Optional) Handle Flags**                         | Decide how existing flags (`-source`, `-output`) interact with the TUI (e.g., pre-fill values, bypass TUI for non-interactive use).        | S      | `main.go`, `tui/model.go`                  |

**Effort Estimation:** S = Small (1-2 hours), M = Medium (2-5 hours), L = Large (5+ hours)

## 3. Implementation Details

### 3.1. TUI Package Structure

Create a new package `tui` to encapsulate all Bubble Tea related code.

```
resumake/
├── tui/
│   ├── model.go        # Main bubbletea model, Init, Update, View
│   ├── messages.go     # Custom tea.Msg types
│   ├── commands.go     # tea.Cmd functions wrapping core logic
│   ├── view_*.go       # View functions for different states
│   ├── styles.go       # lipgloss styles
│   └── tui_test.go     # Unit tests for TUI logic
├── api/
├── input/
├── output/
├── prompt/
├── main.go
├── go.mod
└── ...
```

### 3.2. Main Model (`tui/model.go`)

```go
package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	// other imports
)

type state int

const (
	stateWelcome state = iota
	stateInputSourcePath
	stateInputStdin
	stateConfirmGenerate
	stateGenerating
	stateResultSuccess
	stateResultError
)

type model struct {
	state         state
	apiKeyOk      bool
	sourcePathInput textinput.Model
	stdinInput    textarea.Model
	spinner       spinner.Model
	sourceContent string // Content read from file
	stdinContent  string // Content from textarea
	outputPath    string
	resultMessage string
	errorMsg      string
	// Potentially window dimensions: width, height
}

func InitialModel() model {
	// Initialize textinput, textarea, spinner, etc.
	// Check API Key here or via an initial command
	return model{
		state:    stateWelcome,
		apiKeyOk: checkAPIKey(), // Simple sync check initially
		// ... initialize components
	}
}

func (m model) Init() tea.Cmd {
	// Return initial commands, e.g., spinner.Tick, textinput.Blink
	return nil // Or spinner.Tick if starting in generating state, etc.
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key presses globally or delegate based on state
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case apiKeyCheckedMsg:
		m.apiKeyOk = msg.ok
		if !msg.ok {
			m.state = stateResultError
			m.errorMsg = "GEMINI_API_KEY not set or invalid."
		}
		// Potentially transition state
	case fileReadMsg:
		if msg.err != nil {
			m.state = stateResultError
			m.errorMsg = fmt.Sprintf("Error reading source file: %v", msg.err)
		} else {
			m.sourceContent = msg.content
			m.state = stateInputStdin // Transition to next state
			cmds = append(cmds, m.stdinInput.Focus()) // Focus the text area
		}
	case apiResultMsg:
		m.state = stateResultSuccess // Or stateResultError
		m.resultMessage = msg.outputMessage
		m.errorMsg = msg.errMsg
		m.outputPath = msg.outputPath
	// Handle other custom messages (e.g., stdinFinishedMsg)
	// Handle messages for embedded bubbles (textinput, textarea)
	}

	// Delegate updates to focused components based on state
	switch m.state {
	case stateInputSourcePath:
		m.sourcePathInput, cmd = m.sourcePathInput.Update(msg)
		cmds = append(cmds, cmd)
		// Handle Enter key to trigger file reading command
	case stateInputStdin:
		m.stdinInput, cmd = m.stdinInput.Update(msg)
		cmds = append(cmds, cmd)
		// Handle specific key (e.g., Ctrl+D) to finish input
	case stateGenerating:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Render UI based on m.state
	switch m.state {
	case stateWelcome:
		return viewWelcome(m)
	case stateInputSourcePath:
		return viewInputSourcePath(m)
	// ... other states
	case stateResultSuccess:
		return viewResultSuccess(m)
	case stateResultError:
		return viewResultError(m)
	default:
		return "Unknown state"
	}
}

// Helper function (can be moved)
func checkAPIKey() bool {
	_, err := api.GetAPIKey()
	return err == nil
}

```

### 3.3. Messages (`tui/messages.go`)

Define custom messages for communication between commands and the `Update` function.

```go
package tui

import tea "github.com/charmbracelet/bubbletea"

// Example messages
type apiKeyCheckedMsg struct{ ok bool }
type fileReadMsg struct {
	content string
	err     error
}
type stdinFinishedMsg struct{ content string }
type apiResultMsg struct {
	outputMessage string // Formatted success message
	outputPath    string // Path to the generated file
	errMsg        string // Error message if failed
}
type generationStepMsg struct{ step string } // For progress updates
```

### 3.4. Commands (`tui/commands.go`)

Wrap I/O and long-running operations in `tea.Cmd`.

```go
package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/phrazzld/resumake/api"
	"github.com/phrazzld/resumake/input"
	"github.com/phrazzld/resumake/output"
	"github.com/phrazzld/resumake/prompt"
	// other imports
)

func readFileCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		content, err := input.ReadSourceFile(filePath) // Use existing logic
		return fileReadMsg{content: content, err: err}
	}
}

func generateResumeCmd(sourceContent, stdinContent, outputFlagPath string) tea.Cmd {
	return func() tea.Msg {
		// This command now encapsulates the core logic from main.go
		apiKey, err := api.GetAPIKey()
		if err != nil {
			return apiResultMsg{errMsg: fmt.Sprintf("API Key Error: %v", err)}
		}

		ctx := context.Background()
		client, model, err := api.InitializeClient(ctx, apiKey)
		if err != nil {
			return apiResultMsg{errMsg: fmt.Sprintf("API Client Init Error: %v", err)}
		}
		defer client.Close()

		// Consider sending progress messages back using tea.Batch? More complex.
		// tea.Batch(tea.Printf("Building prompt..."))() // Example, needs better handling

		promptContent := prompt.GeneratePromptContent(sourceContent, stdinContent)

		response, err := api.ExecuteRequest(ctx, model, promptContent)
		if err != nil {
			return apiResultMsg{errMsg: fmt.Sprintf("API Request Error: %v", err)}
		}

		markdownContent, err := output.ProcessResponseContent(response)
		// Handle potential partial recovery from ProcessResponseContent if needed
		if err != nil {
			// Check for MaxTokens and try recovery
			if response != nil && len(response.Candidates) > 0 &&
				response.Candidates[0].FinishReason == genai.FinishReasonMaxTokens {
				partialContent, recoverErr := api.TryRecoverPartialContent(response)
				if recoverErr == nil && partialContent != "" {
					markdownContent = partialContent // Use partial content
					// Add warning to the success message later
				} else {
					return apiResultMsg{errMsg: fmt.Sprintf("API Response Processing Error: %v", err)}
				}
			} else {
				return apiResultMsg{errMsg: fmt.Sprintf("API Response Processing Error: %v", err)}
			}
		}

		outputPath, err := output.WriteOutput(markdownContent, outputFlagPath)
		if err != nil {
			return apiResultMsg{errMsg: fmt.Sprintf("File Write Error: %v", err)}
		}

		// Construct success message
		successMsg := fmt.Sprintf("----- RESUME GENERATION COMPLETE -----\n"+
			"Output file:      %s\n"+
			"Content length:   %d characters\n"+
			"\nNext steps:\n"+
			"  * Review your resume at %s\n"+
			"  * Make any necessary edits\n"+
			"  * Convert to other formats (PDF, DOCX)\n"+
			"---------------------------------------",
			outputPath, len(markdownContent), outputPath)

		// Add truncation warning if needed
		if response != nil && len(response.Candidates) > 0 &&
			response.Candidates[0].FinishReason == genai.FinishReasonMaxTokens {
			successMsg += "\n\nWarning: Output may be truncated due to token limits."
		}


		return apiResultMsg{outputMessage: successMsg, outputPath: outputPath}
	}
}
```

### 3.5. Main Function (`main.go`)

```go
package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/phrazzld/resumake/tui" // Import the new TUI package
)

func main() {
	// Optional: Handle flags for non-interactive mode or config before starting TUI
	// flags, err := input.ParseFlags() ...

	// Initialize the Bubble Tea model
	m := tui.InitialModel()

	// Start the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen()) // Use AltScreen for cleaner exit
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
		os.Exit(1)
	}

	fmt.Println("\nResumake finished.") // Optional message after TUI exits
}

```

## 4. Potential Challenges & Considerations

*   **Error Handling:** Errors from commands need to be gracefully handled and displayed in the TUI without crashing. The `apiResultMsg` includes an `errMsg` field for this. Need consistent error display across states.
*   **State Management:** The main model can become complex. Consider breaking down the model or using sub-models if needed, although for this scope, a single model with states might suffice.
*   **Testing TUIs:** Unit testing the `Update` function with mock messages is feasible. Testing the `View` output can be done using golden files or string comparisons. Full integration testing is harder; manual testing is essential. `main_test.go` will likely need significant changes as it tests the old CLI execution flow.
*   **Terminal Compatibility:** Ensure the TUI renders correctly across different terminals and platforms (Windows, macOS, Linux). `bubbletea` handles much of this, but testing is needed.
*   **Asynchronous Operations:** Commands run asynchronously. The UI needs to reflect loading/generating states correctly while waiting for messages. Spinners and status updates help.
*   **Input Handling:** Capturing multi-line input requires the `textarea` bubble. Defining the "finish" signal (e.g., Ctrl+D, specific key) needs clear user instructions.
*   **Flag Interaction:** Decide if flags should override TUI prompts (e.g., `-source file.md` skips the source input step) or if they are ignored when the TUI runs. A non-interactive mode triggered by specific flags might be useful.
*   **Dependencies:** Adding `bubbletea` and potentially `bubbles`, `lipgloss` increases binary size and dependencies.

## 5. Testing Strategy

*   **Unit Tests:**
    *   Test the `tui.model.Update` function with various `tea.Msg` inputs (key presses, custom messages) to verify correct state transitions and command generation. Mock messages and model states.
    *   Test individual `tea.Cmd` functions where possible, mocking dependencies (like API client or filesystem) if necessary, to ensure they return the correct `tea.Msg`.
    *   Test view helper functions if they contain complex logic.
*   **Integration Tests:**
    *   The existing `main_test.go` relies on executing the binary and checking stdout/stderr/exit codes. This approach is not suitable for testing the interactive TUI directly.
    *   Consider mocking the `tea.Program` runner or creating specific test entry points that run the TUI with controlled inputs/mocked commands for specific flows, although this can be complex.
    *   Focus unit tests on `Update` logic, as it's the core of the TUI behavior.
*   **Manual Testing:**
    *   Run the application interactively on different platforms (Linux, macOS, Windows Terminal).
    *   Test all user flows: with/without source file, empty inputs, valid inputs.
    *   Test error conditions: missing API key, invalid file path, API errors (mocking might be needed here or use temporary invalid keys), network issues.
    *   Test window resizing and different terminal emulators.
    *   Test quitting at different stages (Ctrl+C, Esc).

## 6. Open Questions

1.  **Flag Handling:** How should existing CLI flags (`-source`, `-output`) interact with the TUI?
    *   Option A: Ignore flags when TUI runs (TUI always prompts).
    *   Option B: Use flags to pre-fill TUI inputs.
    *   Option C: Add a `-non-interactive` flag to use the old CLI logic (requires maintaining both flows or adapting core logic).
    *   *Recommendation:* Start with Option B (pre-fill) for simplicity.
2.  **Input Termination:** What is the preferred way for the user to signal the end of stdin input in the `textarea`? (e.g., Ctrl+D, Ctrl+S, specific key combo?). Needs clear instruction in the UI.
3.  **Progress Granularity:** How detailed should the progress updates be during the "Generating" state? (e.g., "Building prompt...", "Calling API...", "Processing response..."). Requires commands to potentially send multiple messages.
4.  **File Picker:** Is a simple text input for the source file path sufficient, or should we investigate a TUI file picker bubble (might add complexity)?
    *   *Recommendation:* Start with `textinput`.
```