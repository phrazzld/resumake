package tui

import (
	"strings"
	"testing"
	
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func TestRenderWelcomeView(t *testing.T) {
	// Test with valid API key
	model := Model{
		apiKeyOk: true,
	}
	
	validKeyView := renderWelcomeView(model)
	
	// The welcome view with a valid API key should:
	// 1. Contain a welcome message
	if !strings.Contains(validKeyView, "Welcome to Resumake") {
		t.Error("Welcome view should contain a welcome message")
	}
	
	// 2. Indicate that the API key is valid
	if !strings.Contains(validKeyView, "API key is valid") {
		t.Error("Welcome view should indicate API key is valid")
	}
	
	// 3. Include instructions on how to proceed
	if !strings.Contains(validKeyView, "Press Enter to continue") {
		t.Error("Welcome view should include instructions to proceed")
	}
	
	// 4. Include information about what the application does
	if !strings.Contains(validKeyView, "This tool helps you create a professional resume") {
		t.Error("Welcome view should include information about the application purpose")
	}
	
	// 5. Include keyboard shortcut info
	if !strings.Contains(validKeyView, "Ctrl+C") {
		t.Error("Welcome view should include keyboard shortcut information")
	}
	
	// Test with invalid API key
	model = Model{
		apiKeyOk: false,
	}
	
	invalidKeyView := renderWelcomeView(model)
	
	// The welcome view with an invalid API key should:
	// 1. Contain a welcome message
	if !strings.Contains(invalidKeyView, "Welcome to Resumake") {
		t.Error("Welcome view should contain a welcome message")
	}
	
	// 2. Indicate that the API key is invalid
	if !strings.Contains(invalidKeyView, "API key is missing or invalid") {
		t.Error("Welcome view should indicate API key is missing or invalid")
	}
	
	// 3. Include instructions on how to set the API key
	if !strings.Contains(invalidKeyView, "GEMINI_API_KEY") {
		t.Error("Welcome view should include instructions to set the API key")
	}
}

func TestRenderSourceFileInputView(t *testing.T) {
	// Initialize text input
	sourceInput := textinput.New()
	sourceInput.Placeholder = "Enter path to existing resume (optional)"
	sourceInput.SetValue("/path/to/file.md")
	
	// Create model with the source path input
	model := Model{
		sourcePathInput: sourceInput,
		flagSourcePath:  "/path/from/flags.md",
	}
	
	// Get the rendered view
	sourceInputView := renderSourceFileInputView(model)
	
	// The source file input view should:
	// 1. Contain a title or heading about source file
	if !strings.Contains(sourceInputView, "Source File") {
		t.Error("Source file input view should contain a title about source files")
	}
	
	// 2. Display the textinput component
	if !strings.Contains(sourceInputView, sourceInput.View()) {
		t.Error("Source file input view should display the textinput component")
	}
	
	// 3. Show if a file path was provided via flags
	if !strings.Contains(sourceInputView, "from command line flags") {
		t.Error("Source file input view should indicate if a path was provided via flags")
	}
	
	// 4. Include instructions about file input being optional
	if !strings.Contains(sourceInputView, "optional") {
		t.Error("Source file input view should indicate that file input is optional")
	}
	
	// 5. Include keyboard shortcut hints
	if !strings.Contains(sourceInputView, "Enter to continue") {
		t.Error("Source file input view should include Enter shortcut")
	}
	
	if !strings.Contains(sourceInputView, "Ctrl+C to quit") {
		t.Error("Source file input view should include quit shortcut")
	}
	
	// Test with empty flag source path
	emptyFlagModel := Model{
		sourcePathInput: sourceInput,
		flagSourcePath:  "",
	}
	
	emptyFlagView := renderSourceFileInputView(emptyFlagModel)
	
	// 6. Should not mention flags if no flag source path was provided
	if strings.Contains(emptyFlagView, "from command line flags") {
		t.Error("Source file input view should not mention flags if no flag source path is provided")
	}
}

func TestRenderStdinInputView(t *testing.T) {
	// Initialize textarea
	stdinTA := textarea.New()
	stdinTA.Placeholder = "Enter details about your experience, skills, etc."
	stdinTA.SetValue("My resume content")
	
	// Create model with the stdin input
	model := Model{
		stdinInput:     stdinTA,
		sourceContent:  "Sample source content",
		sourcePathInput: textinput.New(),
	}
	model.sourcePathInput.SetValue("/path/to/sample.md")
	
	// Get the rendered view
	stdinView := renderStdinInputView(model)
	
	// The stdin input view should:
	// 1. Contain a title or heading about resume details
	if !strings.Contains(stdinView, "Resume Details") {
		t.Error("Stdin input view should contain a title about resume details")
	}
	
	// 2. Display the textarea component
	if !strings.Contains(stdinView, stdinTA.View()) {
		t.Error("Stdin input view should display the textarea component")
	}
	
	// 3. Include instructions for providing resume details
	if !strings.Contains(stdinView, "experience") || !strings.Contains(stdinView, "skills") {
		t.Error("Stdin input view should include instructions about what content to provide")
	}
	
	// 4. Show source file information if a source file was provided
	if !strings.Contains(stdinView, "Source file") {
		t.Error("Stdin input view should show source file info when one is provided")
	}
	
	// 5. Include keyboard shortcut hints
	if !strings.Contains(stdinView, "Ctrl+D") {
		t.Error("Stdin input view should include Ctrl+D shortcut to finish input")
	}
	
	if !strings.Contains(stdinView, "Ctrl+C") {
		t.Error("Stdin input view should include Ctrl+C shortcut to quit")
	}
	
	// Test with empty source content
	emptySourceModel := Model{
		stdinInput:    stdinTA,
		sourceContent: "",
	}
	
	emptySourceView := renderStdinInputView(emptySourceModel)
	
	// 6. Should not mention source file if no source content exists
	if strings.Contains(emptySourceView, "Source file") {
		t.Error("Stdin input view should not mention source file if no source content exists")
	}
	
	// 7. Should include helpful tips or examples for resume content
	if !strings.Contains(stdinView, "Tips:") || !strings.Contains(stdinView, "Example:") {
		t.Error("Stdin input view should include helpful tips or examples")
	}
}

func TestRenderGeneratingView(t *testing.T) {
	// Initialize spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	// Create model with spinner and progress information
	model := Model{
		spinner:      sp,
		progressStep: "Processing",
		progressMsg:  "Analyzing your experience...",
		stdinContent: "Sample content for resume",
		sourceContent: "Sample source content",
	}
	
	// Get the rendered view
	generatingView := renderGeneratingView(model)
	
	// The generating view should:
	// 1. Contain a title or heading about generation process
	if !strings.Contains(generatingView, "Generating") {
		t.Error("Generating view should contain a title about the generation process")
	}
	
	// 2. Display the spinner component
	if !strings.Contains(generatingView, sp.View()) {
		t.Error("Generating view should display the spinner component")
	}
	
	// 3. Show progress step if provided
	if !strings.Contains(generatingView, "Processing") {
		t.Error("Generating view should show the current progress step")
	}
	
	// 4. Show progress message if provided
	if !strings.Contains(generatingView, "Analyzing your experience") {
		t.Error("Generating view should show the progress message")
	}
	
	// 5. Include information about inputs
	if !strings.Contains(generatingView, "characters") {
		t.Error("Generating view should include information about the input size")
	}
	
	// Test with empty progress information
	emptyProgressModel := Model{
		spinner:      sp,
		progressStep: "",
		progressMsg:  "",
		stdinContent: "Sample content for resume",
	}
	
	emptyProgressView := renderGeneratingView(emptyProgressModel)
	
	// 6. Should handle missing progress information gracefully
	if !strings.Contains(emptyProgressView, "Please wait") {
		t.Error("Generating view should show a waiting message when no progress info is available")
	}
	
	// 7. Include estimated time information
	if !strings.Contains(generatingView, "may take") {
		t.Error("Generating view should include information about estimated completion time")
	}
}

func TestRenderSuccessView(t *testing.T) {
	// Create model with success information
	model := Model{
		outputPath:    "/tmp/resume_out.md",
		resultMessage: "2500",
		stdinContent:  "Sample content for resume",
		sourceContent: "Sample source content",
	}
	
	// Get the rendered view
	successView := renderSuccessView(model)
	
	// The success view should:
	// 1. Contain a title or heading about successful completion
	if !strings.Contains(successView, "Success") || !strings.Contains(successView, "Complete") {
		t.Error("Success view should contain a title about successful completion")
	}
	
	// 2. Show the output file path
	if !strings.Contains(successView, "/tmp/resume_out.md") {
		t.Error("Success view should show the output file path")
	}
	
	// 3. Show the content length or size
	if !strings.Contains(successView, "2500") {
		t.Error("Success view should show the content length")
	}
	
	// 4. Include next steps or instructions
	if !strings.Contains(successView, "Next Steps") {
		t.Error("Success view should include next steps instructions")
	}
	
	// 5. Include keyboard shortcuts
	if !strings.Contains(successView, "Enter") {
		t.Error("Success view should include Enter shortcut to quit")
	}
	
	// 6. Include a celebratory message or visual element
	if !strings.Contains(successView, "✅") && !strings.Contains(successView, "✓") && !strings.Contains(successView, "congratulations") {
		t.Error("Success view should include a celebratory element")
	}
	
	// 7. Show formatting options
	if !strings.Contains(successView, "PDF") || !strings.Contains(successView, "DOCX") {
		t.Error("Success view should mention formatting options")
	}
}