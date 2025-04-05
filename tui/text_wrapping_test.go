package tui

import (
	"strings"
	"testing"
	
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

func TestWrapTextUsage(t *testing.T) {
	// Create a test model
	model := createTestModelWithAllFields()
	
	// Set a long path that would need wrapping
	longPath := "/this/is/an/extremely/long/path/that/would/definitely/need/wrapping/in/smaller/terminal/windows/resume.md"
	model.sourcePathInput.SetValue(longPath)
	model.width = 60  // Narrow width to force wrapping
	
	// Test each view function
	views := []struct {
		name      string
		renderFn  func(Model) string
	}{
		{"Welcome", renderWelcomeView},
		{"SourceFileInput", renderSourceFileInputView},
		{"StdinInput", renderStdinInputView},
		{"ConfirmGenerate", renderConfirmGenerateView},
		{"Generating", renderGeneratingView},
		{"Success", renderSuccessView},
		{"Error", renderErrorView},
	}
	
	for _, v := range views {
		t.Run(v.name, func(t *testing.T) {
			// Render the view
			output := v.renderFn(model)
			
			// Basic validation
			if output == "" {
				t.Errorf("%s view returned empty string", v.name)
			}
			
			// Check specific views for wrapping behavior
			if v.name == "ConfirmGenerate" {
				// The confirm view specifically should include the file path
				if !strings.Contains(output, "Source file") {
					t.Errorf("%s view doesn't properly display source file path", v.name)
				}
			}
		})
	}
}

// Helper function to create a test model with all fields populated
func createTestModelWithAllFields() Model {
	sourceInput := textinput.New()
	sourceInput.SetValue("/path/to/resume.md")
	
	stdinInput := textarea.New()
	stdinInput.SetValue("Experience: Senior Software Engineer")
	
	s := spinner.New()
	s.Spinner = spinner.Dot
	
	return Model{
		sourcePathInput: sourceInput,
		stdinInput:      stdinInput,
		spinner:         s,
		state:           stateConfirmGenerate,
		width:           80,
		height:          24,
		apiKeyOk:        true,
		flagSourcePath:  "/path/from/flags.md",
		sourceContent:   "Sample source content",
		stdinContent:    "Sample stdin content",
		progressStep:    "Processing",
		progressMsg:     "Analyzing your resume content",
		errorMsg:        "Failed to process content: API error",
		resultMessage:   "2048",
		outputPath:      "/path/to/output/resume_out.md",
	}
}