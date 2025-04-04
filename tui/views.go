package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles for different elements
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0099FF")).
		MarginBottom(1)
	
	subtitleStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#00CCFF"))
	
	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00CC00"))
	
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))
	
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFCC00"))
	
	keyboardHintStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA")).
		Italic(true)
		
	inputLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))
		
	flagValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFAA"))
)

// renderWelcomeView generates the welcome screen content
func renderWelcomeView(m Model) string {
	var content string
	
	// Title
	content += titleStyle.Render("╔═════════════════════════════════╗")
	content += "\n"
	content += titleStyle.Render("║       Welcome to Resumake!      ║")
	content += "\n"
	content += titleStyle.Render("╚═════════════════════════════════╝")
	content += "\n\n"
	
	// Application description
	content += subtitleStyle.Render("This tool helps you create a professional resume from your experience and qualifications.")
	content += "\n\n"
	
	// How it works
	content += "How it works:\n"
	content += "1. Optionally provide an existing resume to enhance\n"
	content += "2. Tell us about your experience, skills, and qualifications\n"
	content += "3. We'll generate a polished resume in Markdown format\n\n"
	
	// API key status
	if m.apiKeyOk {
		content += successStyle.Render("✅ API key is valid and ready to use.")
		content += "\n\n"
		content += "You're all set to create your professional resume!"
	} else {
		content += errorStyle.Render("❌ API key is missing or invalid.")
		content += "\n\n"
		content += "To use Resumake, you need to set the GEMINI_API_KEY environment variable:\n"
		content += "  export GEMINI_API_KEY=your_api_key_here\n\n"
		content += "You can get an API key from: https://makersuite.google.com/app/apikey\n"
		content += errorStyle.Render("Note: Proceeding without a valid API key will result in errors.")
	}
	
	// Keyboard shortcuts and instructions
	content += "\n\n"
	content += keyboardHintStyle.Render("Keyboard shortcuts:")
	content += "\n"
	content += keyboardHintStyle.Render("• Enter: Continue to next step")
	content += "\n"
	content += keyboardHintStyle.Render("• Ctrl+C: Quit application")
	content += "\n"
	content += keyboardHintStyle.Render("• Esc: Go back (when available)")
	content += "\n\n"
	
	// Call to action
	content += infoStyle.Render("Press Enter to continue...")
	
	return content
}

// renderSourceFileInputView generates the source file input view content
func renderSourceFileInputView(m Model) string {
	var content string
	
	// Title
	content += titleStyle.Render("Source File Input")
	content += "\n\n"
	
	// Instructions and explanation
	content += "You can optionally provide an existing resume file to use as a starting point."
	content += "\n"
	content += "This helps generate a better result by incorporating your existing content."
	content += "\n\n"
	
	// Show if path was provided via flags
	if m.flagSourcePath != "" {
		content += infoStyle.Render("A file path was provided from command line flags:")
		content += "\n"
		content += flagValueStyle.Render(m.flagSourcePath)
		content += "\n\n"
		content += "You can edit this path or leave it as is."
		content += "\n\n"
	}
	
	// Input field label
	content += inputLabelStyle.Render("Enter path to existing resume file (optional):")
	content += "\n"
	
	// The actual text input component
	content += m.sourcePathInput.View()
	content += "\n\n"
	
	// Note about being optional
	content += "Note: This step is optional. If you don't have an existing resume file,"
	content += "\n"
	content += "just press Enter to continue without selecting a file."
	content += "\n\n"
	
	// Keyboard shortcuts
	content += keyboardHintStyle.Render("• Enter: Continue to next step")
	content += "\n"
	content += keyboardHintStyle.Render("• Ctrl+C to quit")
	
	return content
}