package tui

import (
	"fmt"
	
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/phrazzld/resumake/api"
)

// State represents the different states of the application.
type State int

const (
	// stateWelcome is the initial state, displaying welcome message and API key status.
	stateWelcome State = iota
	
	// stateInputSourcePath allows the user to input a source file path.
	stateInputSourcePath
	
	// stateInputStdin allows the user to input their resume details in a text area.
	stateInputStdin
	
	// stateConfirmGenerate asks for confirmation before generating the resume.
	stateConfirmGenerate
	
	// stateGenerating shows progress while calling the API and processing the response.
	stateGenerating
	
	// stateResultSuccess shows successful resume generation and output details.
	stateResultSuccess
	
	// stateResultError shows error details if something went wrong.
	stateResultError
)

// Model is the main model for the Bubble Tea application.
type Model struct {
	// Application state
	state         State
	apiKeyOk      bool
	errorMsg      string
	
	// Input components
	sourcePathInput textinput.Model
	stdinInput      textarea.Model
	
	// Content
	sourceContent string // Content read from file
	stdinContent  string // Content from stdin textarea
	
	// Output
	outputPath    string
	resultMessage string
	
	// UI components
	spinner       spinner.Model
	width         int
	height        int
	
	// Styling
	mainStyle     lipgloss.Style
}

// NewModel creates a new Model with default values.
func NewModel() Model {
	// Initialize text input for source file path
	sourceInput := textinput.New()
	sourceInput.Placeholder = "Enter path to existing resume (optional)"
	sourceInput.CharLimit = 150
	sourceInput.Width = 50
	
	// Initialize textarea for stdin input
	stdinTA := textarea.New()
	stdinTA.Placeholder = "Enter details about your experience, skills, etc."
	stdinTA.SetWidth(80)
	stdinTA.SetHeight(20)
	
	// Initialize spinner for loading state
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	// Check API key on startup
	apiKeyOk := checkAPIKey()
	
	return Model{
		state:          stateWelcome,
		apiKeyOk:       apiKeyOk,
		sourcePathInput: sourceInput,
		stdinInput:     stdinTA,
		spinner:        sp,
		mainStyle:      lipgloss.NewStyle().Bold(true),
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	// Initial commands like spinner spinning or cursor blinking
	return tea.Batch(
		tea.Cmd(m.spinner.Tick),
		m.sourcePathInput.Focus(),
	)
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global key handlers
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
		
		// State-specific key handling
		switch m.state {
		case stateWelcome:
			if msg.Type == tea.KeyEnter {
				if m.apiKeyOk {
					m.state = stateInputSourcePath
					cmds = append(cmds, m.sourcePathInput.Focus())
				} else {
					m.state = stateResultError
					m.errorMsg = "API key is missing or invalid. Set GEMINI_API_KEY environment variable."
				}
			}
		
		case stateInputSourcePath:
			// Update source input component
			var inputCmd tea.Cmd
			m.sourcePathInput, inputCmd = m.sourcePathInput.Update(msg)
			cmds = append(cmds, inputCmd)
			
			if msg.Type == tea.KeyEnter {
				// TODO: Implement file reading command
				m.sourceContent = m.sourcePathInput.Value() // Store the path for future file reading
				m.state = stateInputStdin
				cmds = append(cmds, m.stdinInput.Focus())
			}
		
		case stateInputStdin:
			// Update textarea component
			var textareaCmd tea.Cmd
			m.stdinInput, textareaCmd = m.stdinInput.Update(msg)
			cmds = append(cmds, textareaCmd)
			
			// Ctrl+D to finish input and proceed
			if msg.Type == tea.KeyCtrlD {
				m.stdinContent = m.stdinInput.Value() // Store the content
				m.state = stateConfirmGenerate
			}
		
		case stateConfirmGenerate:
			if msg.Type == tea.KeyEnter {
				m.state = stateGenerating
			} else if msg.Type == tea.KeyEsc {
				m.state = stateInputStdin
				cmds = append(cmds, m.stdinInput.Focus())
			}
		}
	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update component dimensions
		m.sourcePathInput.Width = msg.Width - 20
		m.stdinInput.SetWidth(msg.Width - 20)
		m.stdinInput.SetHeight(msg.Height - 10)
	}
	
	// If no commands were queued from above, don't include spinner.Tick
	if len(cmds) == 0 {
		return m, nil
	}
	
	// Only include spinner.Tick if we're in a state that shows the spinner
	if m.state == stateGenerating {
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
	}
	
	return m, tea.Batch(cmds...)
}

// View renders the model to a string.
func (m Model) View() string {
	var content string
	
	// Render different views based on the current state
	switch m.state {
	case stateWelcome:
		content = "Welcome to Resumake!\n\n"
		
		if m.apiKeyOk {
			content += "✅ API key is valid.\n\n"
			content += "Press Enter to continue..."
		} else {
			content += "❌ API key is missing or invalid.\n\n"
			content += "Please set the GEMINI_API_KEY environment variable and restart.\n"
			content += "Press Enter to continue anyway..."
		}
	
	case stateInputSourcePath:
		content = "Enter the path to an existing resume (optional):\n\n"
		content += m.sourcePathInput.View() + "\n\n"
		content += "Press Enter to continue, Ctrl+C to quit."
	
	case stateInputStdin:
		content = "Tell us about your experience, skills, and qualifications:\n\n"
		content += m.stdinInput.View() + "\n\n"
		content += "Ctrl+D when finished, Ctrl+C to quit."
	
	case stateConfirmGenerate:
		content = "Ready to generate your resume!\n\n"
		if m.sourceContent != "" {
			content += "Source file: " + m.sourceContent + "\n"
		}
		content += "Input length: " + fmt.Sprintf("%d", len(m.stdinContent)) + " characters\n\n"
		content += "Press Enter to confirm, Esc to go back."
	
	case stateGenerating:
		content = "Generating your resume...\n\n"
		content += m.spinner.View() + " Please wait, this may take a minute."
	
	case stateResultSuccess:
		content = "----- RESUME GENERATION COMPLETE -----\n"
		content += "Output file:      " + m.outputPath + "\n"
		content += "Content length:   " + m.resultMessage + "\n\n"
		content += "Next steps:\n"
		content += "  * Review your resume at " + m.outputPath + "\n"
		content += "  * Make any necessary edits\n"
		content += "  * Convert to other formats (PDF, DOCX)\n"
		content += "---------------------------------------\n\n"
		content += "Press Enter to quit."
	
	case stateResultError:
		content = "Error generating resume:\n\n"
		content += m.errorMsg + "\n\n"
		content += "Press Enter to quit."
	
	default:
		content = "Unknown state"
	}
	
	// Apply main style to the content
	return m.mainStyle.Render(content)
}

// Helper function to check if the API key is available and valid
func checkAPIKey() bool {
	_, err := api.GetAPIKey()
	return err == nil
}