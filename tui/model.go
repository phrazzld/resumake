package tui

import (
	"context"
	"fmt"
	
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/generative-ai-go/genai"
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
	
	// Flag-provided values
	flagSourcePath string
	flagOutputPath string
	
	// Status messages
	progressStep  string
	progressMsg   string
	
	// API client instances
	apiClient     *genai.Client       // Initialized API client instance
	apiModel      *genai.GenerativeModel // Initialized model instance
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
	
	// Initialize spinner for loading state with more visible spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2C94C")).Bold(true)
	// Important: use a spinner with more visible animation frames
	sp.Spinner = spinner.Spinner{
		Frames: []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
		FPS:    12, // Faster animation
	}
	
	// Check API key on startup
	apiKeyOk := checkAPIKey()
	
	return Model{
		state:          stateWelcome,
		apiKeyOk:       apiKeyOk,
		sourcePathInput: sourceInput,
		stdinInput:     stdinTA,
		spinner:        sp,
		mainStyle:      lipgloss.NewStyle().Bold(true),
		// Flag values will be populated with WithSourcePath/WithOutputPath
		flagSourcePath: "",
		flagOutputPath: "",
		// API client instances start as nil and will be initialized as needed
		apiClient:      nil,
		apiModel:       nil,
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
	// Handle tea.QuitMsg to clean up resources
	case tea.QuitMsg:
		m = cleanupAPIClient(m)
		return m, tea.Quit
		
	// Handle custom messages from commands
	case FileReadResultMsg:
		if msg.Success {
			m.sourceContent = msg.Content
		} else {
			m.state = stateResultError
			m.errorMsg = msg.Error.Error()
			return m, nil
		}
		
	case APIInitResultMsg:
		if !msg.Success {
			m.state = stateResultError
			m.errorMsg = msg.Error.Error()
			return m, nil
		}
		
	case APIResultMsg:
		if msg.Success {
			m.state = stateResultSuccess
			m.outputPath = msg.OutputPath
			m.resultMessage = fmt.Sprintf("%d", len(msg.Content))
		} else {
			m.state = stateResultError
			m.errorMsg = msg.Error.Error()
		}
		return m, nil
		
	case StdinSubmitMsg:
		m.stdinContent = msg.Content
		m.state = stateConfirmGenerate
		return m, nil
		
	case ProgressUpdateMsg:
		m.progressStep = msg.Step
		m.progressMsg = msg.Message
		
	case tea.KeyMsg:
		// Global key handlers
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m = cleanupAPIClient(m)
			return m, tea.Quit
		}
		
		// State-specific key handling
		switch m.state {
		case stateWelcome:
			if msg.Type == tea.KeyEnter {
				if m.apiKeyOk {
					// Initialize API client here when we confirm a valid API key
					// This is the earliest point where we need the API client
					var err error
					m, err = initializeAPIClient(m)
					if err != nil {
						m.state = stateResultError
						m.errorMsg = err.Error()
						return m, nil
					}
					
					// If a source path was provided via flags, we can pre-fill it
					if m.flagSourcePath != "" {
						// We'll still go to the input screen but with pre-filled value
						m.state = stateInputSourcePath
						cmds = append(cmds, m.sourcePathInput.Focus())
					} else {
						// Otherwise, prompt for source path
						m.state = stateInputSourcePath
						cmds = append(cmds, m.sourcePathInput.Focus())
					}
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
				// Use the file reading command to read the source file
				filePath := m.sourcePathInput.Value()
				m.state = stateInputStdin
				cmds = append(cmds, 
					ReadSourceFileCmd(filePath),  // Read the file asynchronously
					m.stdinInput.Focus(),         // Focus the text area
				)
			}
		
		case stateInputStdin:
			// Update textarea component
			var textareaCmd tea.Cmd
			m.stdinInput, textareaCmd = m.stdinInput.Update(msg)
			cmds = append(cmds, textareaCmd)
			
			// Ctrl+D to finish input and proceed
			if msg.Type == tea.KeyCtrlD {
				// Submit the stdin input using our command
				cmds = append(cmds, SubmitStdinInputCmd(m.stdinInput.Value()))
			}
		
		case stateConfirmGenerate:
			if msg.Type == tea.KeyEnter {
				m.state = stateGenerating
				
				// Use provided output path from flags if available
				outputPath := ""
				if m.flagOutputPath != "" {
					outputPath = m.flagOutputPath
				}
				
				// Add progress update and API commands
				cmds = append(cmds, 
					SendProgressUpdateCmd("Starting", "Initializing resume generation..."),
					GenerateResumeCmd(m.sourceContent, m.stdinContent, outputPath, false),
				)
			} else if msg.Type == tea.KeyEsc {
				m.state = stateInputStdin
				cmds = append(cmds, m.stdinInput.Focus())
			}
			
		case stateResultSuccess, stateResultError:
			// Any key in final states quits the application
			if msg.Type == tea.KeyEnter {
				m = cleanupAPIClient(m)
				return m, tea.Quit
			}
		}
	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update component dimensions with minimum widths to prevent text truncation
		inputWidth := msg.Width - 20
		if inputWidth < 60 {
			inputWidth = 60
		}
		
		textareaHeight := msg.Height - 10
		if textareaHeight < 10 {
			textareaHeight = 10
		}
		
		m.sourcePathInput.Width = inputWidth
		m.stdinInput.SetWidth(inputWidth)
		m.stdinInput.SetHeight(textareaHeight)
	}
	
	// If no commands were queued from above, don't include spinner.Tick
	if len(cmds) == 0 {
		return m, nil
	}
	
	// Always include spinner.Tick if we're in the generating state
	// This ensures the spinner animates properly
	if m.state == stateGenerating {
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
		
		// Always ensure the spinner keeps ticking by adding the tick command
		// This is crucial to keep animation going
		cmds = append(cmds, m.spinner.Tick)
	}
	
	return m, tea.Batch(cmds...)
}

// View renders the model to a string.
func (m Model) View() string {
	var content string
	
	// Render different views based on the current state
	switch m.state {
	case stateWelcome:
		content = renderWelcomeView(m)
	
	case stateInputSourcePath:
		content = renderSourceFileInputView(m)
	
	case stateInputStdin:
		content = renderStdinInputView(m)
	
	case stateConfirmGenerate:
		content = "Ready to generate your resume!\n\n"
		if m.sourceContent != "" {
			content += "Source file: " + m.sourcePathInput.Value() + "\n"
		}
		content += "Input length: " + fmt.Sprintf("%d", len(m.stdinContent)) + " characters\n\n"
		
		// Show output path if it was provided via flags
		if m.flagOutputPath != "" {
			content += "Output will be written to: " + m.flagOutputPath + "\n\n"
		}
		
		content += "Press Enter to confirm, Esc to go back."
	
	case stateGenerating:
		content = renderGeneratingView(m)
	
	case stateResultSuccess:
		content = renderSuccessView(m)
	
	case stateResultError:
		content = renderErrorView(m)
	
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

// initializeAPIClient initializes the API client and model if needed
// Returns the modified model and any error that occurred
func initializeAPIClient(m Model) (Model, error) {
	// Skip initialization if already done
	if m.apiClient != nil && m.apiModel != nil {
		return m, nil
	}
	
	// Get API key
	apiKey, err := api.GetAPIKey()
	if err != nil {
		return m, fmt.Errorf("API key error: %w", err)
	}
	
	// Initialize client and model
	ctx := context.Background()
	client, model, err := api.InitializeClient(ctx, apiKey)
	if err != nil {
		return m, fmt.Errorf("failed to initialize API client: %w", err)
	}
	
	// Store the instances in the model
	m.apiClient = client
	m.apiModel = model
	
	return m, nil
}

// cleanupAPIClient closes the API client if it was initialized
func cleanupAPIClient(m Model) Model {
	if m.apiClient != nil {
		m.apiClient.Close()
		m.apiClient = nil
		m.apiModel = nil
	}
	return m
}

// WithSourcePath returns a copy of the model with the source path set
// Used when the source path is provided via command-line flags
func (m Model) WithSourcePath(path string) Model {
	m.flagSourcePath = path
	m.sourcePathInput.SetValue(path)
	return m
}

// WithOutputPath returns a copy of the model with the output path set
// Used when the output path is provided via command-line flags
func (m Model) WithOutputPath(path string) Model {
	m.flagOutputPath = path
	return m
}