package tui

// This file defines the message types used by the Bubble Tea commands.
// Messages are returned by commands to update the model state.

// FileReadResultMsg is returned when a file read operation completes.
type FileReadResultMsg struct {
	Success bool   // Whether the file read was successful
	Content string // The content of the file (if successful)
	Error   error  // The error that occurred (if unsuccessful)
}

// APIInitResultMsg is returned when API initialization completes.
type APIInitResultMsg struct {
	Success bool  // Whether the API initialization was successful
	Error   error // The error that occurred (if unsuccessful)
}

// APIResultMsg is returned when an API request completes.
type APIResultMsg struct {
	Success      bool   // Whether the API request was successful
	Content      string // The generated content (if successful)
	OutputPath   string // The path where the content was written
	TruncatedMsg string // Warning message if the output was truncated
	Error        error  // The error that occurred (if unsuccessful)
}

// StdinSubmitMsg is sent when the user submits stdin input.
type StdinSubmitMsg struct {
	Content string // The content entered by the user
}

// ProgressUpdateMsg is sent to update the UI with progress information.
type ProgressUpdateMsg struct {
	Step    string // The current step being executed
	Message string // Additional message about the progress
}