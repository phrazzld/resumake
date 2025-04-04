# TODO

## Bubble Tea TUI Implementation

### Setup & Structure
- [x] Add Bubble Tea and related dependencies
  - Description: Add github.com/charmbracelet/bubbletea, lipgloss, and bubbles to go.mod
  - Dependencies: None
  - Priority: High

- [x] Create basic TUI package structure
  - Description: Create tui directory and basic file structure (model.go, messages.go, etc.)
  - Dependencies: Added dependencies
  - Priority: High

- [x] Design TUI model and states
  - Description: Define the state enum and model struct with all required fields
  - Dependencies: TUI package structure
  - Priority: High

### Main Application Flow
- [x] Refactor main.go for TUI
  - Description: Update main.go to initialize and run Bubble Tea program
  - Dependencies: TUI model design
  - Priority: High

- [x] Implement initial welcome view
  - Description: Create welcome screen with API key check and instructions
  - Dependencies: TUI model implementation
  - Priority: Medium

- [x] Implement command pattern for core logic
  - Description: Create tea.Cmd functions to wrap existing logic (file reading, API calls)
  - Dependencies: TUI model implementation
  - Priority: High

### Input Handling
- [x] Implement source file input view
  - Description: Create view for inputting source file path with textinput bubble
  - Dependencies: TUI model and welcome view
  - Priority: Medium

- [x] Implement stdin text area input
  - Description: Create view for multi-line text input using textarea bubble
  - Dependencies: TUI model implementation
  - Priority: Medium

- [x] Handle flag pre-filling
  - Description: Use flag values to pre-fill TUI inputs if provided
  - Dependencies: Source file and stdin input views
  - Priority: Low

### API Integration & Processing
- [x] Implement generating state view
  - Description: Create view with spinner and status messages during API call
  - Dependencies: Input handling views
  - Priority: Medium

- [x] Wrap API call logic in tea.Cmd
  - Description: Create command function to handle API requests asynchronously
  - Dependencies: Input handling implementation
  - Priority: High

- [x] Implement progress updates during generation
  - Description: Add ability to show progress steps during generation process
  - Dependencies: Generating state view
  - Priority: Low

### Result Handling
- [x] Implement success view
  - Description: Create view to show successful resume generation with output path
  - Dependencies: API integration
  - Priority: Medium

- [x] Implement error view
  - Description: Create view to display detailed error messages
  - Dependencies: API integration
  - Priority: Medium

### Styling & Polish
- [x] Create consistent styling with lipgloss
  - Description: Define styles for UI components in styles.go
  - Dependencies: All views implemented
  - Priority: Low

- [x] Add keyboard shortcuts help
  - Description: Add visible keyboard shortcut hints in views (Ctrl+C to quit, etc.)
  - Dependencies: All views implemented
  - Priority: Low

### Testing
- [x] Update existing tests
  - Description: Modify main_test.go to work with new TUI approach
  - Dependencies: All implementation complete
  - Priority: Medium

- [x] Add unit tests for TUI model update logic
  - Description: Create tests for state transitions and message handling
  - Dependencies: TUI model implementation
  - Priority: Medium

- [x] Add tests for command functions
  - Description: Test tea.Cmd functions with mocked dependencies
  - Dependencies: Command implementation
  - Priority: Medium

## Assumptions & Clarifications
- The existing API integration will remain unchanged, just wrapped in tea.Cmd functions
- Flags will be used to pre-fill TUI inputs rather than bypassing the TUI completely
- Basic keyboard controls will include Ctrl+C to quit, arrow keys for navigation, Enter to submit
- Stdin input will be terminated with Ctrl+D (standard convention)
- Error handling will be comprehensive with helpful messages in the TUI
