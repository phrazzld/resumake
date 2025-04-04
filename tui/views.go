package tui

import (
	"fmt"
	"strings"
	
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
		
	tipStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88FF88")).
		Italic(true)
		
	exampleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BBBBFF"))
		
	progressStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFAA44")).
		Bold(true)
		
	stepStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44AAFF")).
		Bold(true)
		
	completedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#33FF33"))
		
	pathStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#555555")).
		Padding(0, 1)
		
	errorTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF0000")).
		Padding(0, 1).
		MarginBottom(1)
		
	errorMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5555")).
		Bold(true)
		
	troubleshootStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AACCFF"))
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

// renderStdinInputView generates the stdin textarea input view content
func renderStdinInputView(m Model) string {
	var content string
	
	// Title
	content += titleStyle.Render("Resume Details")
	content += "\n\n"
	
	// Instructions
	content += subtitleStyle.Render("Tell us about your experience, skills, and qualifications.")
	content += "\n\n"
	
	// Show source file info if one was provided
	if m.sourceContent != "" {
		content += "Source file: " + successStyle.Render(m.sourcePathInput.Value())
		content += "\n"
		content += fmt.Sprintf("Content length: %d characters", len(m.sourceContent))
		content += "\n\n"
		content += "We'll combine this source content with the details you provide below."
		content += "\n\n"
	}
	
	// Input field label
	content += inputLabelStyle.Render("Enter your resume details:")
	content += "\n\n"
	
	// Tips and examples
	content += tipStyle.Render("Tips:")
	content += "\n"
	content += tipStyle.Render("• Be specific about your achievements")
	content += "\n"
	content += tipStyle.Render("• Include relevant skills and technologies")
	content += "\n"
	content += tipStyle.Render("• Mention education, certifications, and years of experience")
	content += "\n\n"
	
	content += exampleStyle.Render("Example:")
	content += "\n"
	content += exampleStyle.Render("I have 5 years of experience as a software developer.")
	content += "\n"
	content += exampleStyle.Render("Skills: Go, Python, JavaScript, Docker, Kubernetes")
	content += "\n"
	content += exampleStyle.Render("Education: BS in Computer Science, University of Example")
	content += "\n\n"
	
	// The actual textarea component
	content += m.stdinInput.View()
	content += "\n\n"
	
	// Keyboard shortcuts
	content += keyboardHintStyle.Render("• Ctrl+D when finished")
	content += "\n"
	content += keyboardHintStyle.Render("• Ctrl+C to quit")
	
	return content
}

// renderGeneratingView generates the view shown during resume generation
func renderGeneratingView(m Model) string {
	var content string
	
	// Title with animated spinner
	content += titleStyle.Render("Generating Your Resume")
	content += "\n\n"
	
	// Spinner and status
	content += m.spinner.View() + " " + progressStyle.Render("Processing your information")
	content += "\n\n"
	
	// Input information
	totalChars := len(m.stdinContent) + len(m.sourceContent)
	content += fmt.Sprintf("Processing %d characters of input...", totalChars)
	content += "\n"
	
	// Source file info if provided
	if m.sourceContent != "" {
		content += "Source file: " + m.sourcePathInput.Value()
		content += "\n"
	}
	
	// Estimated time
	content += "\n"
	content += infoStyle.Render("This may take up to 60 seconds depending on the input size.")
	content += "\n\n"
	
	// Progress information
	if m.progressStep != "" && m.progressMsg != "" {
		content += stepStyle.Render("Step: " + m.progressStep)
		content += "\n"
		content += m.progressMsg
		content += "\n\n"
	} else {
		content += "Please wait while we generate your resume..."
		content += "\n\n"
	}
	
	// Additional status info
	content += "The Gemini API is analyzing your experience and crafting a professional resume."
	content += "\n"
	content += "You'll be able to review and save the result when it's complete."
	
	return content
}

// renderSuccessView generates the view shown after successful resume generation
func renderSuccessView(m Model) string {
	var content string
	
	// Title with success indicator
	content += completedStyle.Render("✅ Success! Resume Generation Complete! ✅")
	content += "\n\n"
	
	// Border for output path
	content += "Your resume has been successfully generated and saved to:"
	content += "\n"
	content += pathStyle.Render(m.outputPath)
	content += "\n\n"
	
	// Show result information
	content += fmt.Sprintf("Content length: %s characters", m.resultMessage)
	content += "\n\n"
	
	// Next steps with a box
	content += titleStyle.Render("╔═════════════════════════╗")
	content += "\n"
	content += titleStyle.Render("║ Next Steps              ║")
	content += "\n"
	content += titleStyle.Render("╚═════════════════════════╝")
	content += "\n\n"
	
	// Detailed next steps list
	content += successStyle.Render("1. Review your resume at ") + m.outputPath
	content += "\n"
	content += successStyle.Render("2. Make any necessary edits or refinements")
	content += "\n"
	content += successStyle.Render("3. Convert to other formats:")
	content += "\n"
	content += "   • PDF: Use a markdown converter like pandoc"
	content += "\n"
	content += "   • DOCX: Use a markdown converter or import into your word processor"
	content += "\n\n"
	
	// Congratulatory message
	content += infoStyle.Render("Congratulations on creating your professional resume!")
	content += "\n\n"
	
	// Keyboard shortcuts
	content += keyboardHintStyle.Render("Press Enter to quit")
	
	return content
}

// renderErrorView generates the view shown when an error occurs
func renderErrorView(m Model) string {
	var content string
	
	// Title with error indicator
	content += errorTitleStyle.Render("❌ ERROR: Resume Generation Failed ❌")
	content += "\n\n"
	
	// Error message
	content += "The following error occurred while generating your resume:"
	content += "\n\n"
	content += errorMsgStyle.Render(m.errorMsg)
	content += "\n\n"
	
	// Troubleshooting section
	content += titleStyle.Render("Troubleshooting")
	content += "\n\n"
	
	// Different troubleshooting suggestions based on error type
	if strings.Contains(strings.ToLower(m.errorMsg), "api") {
		content += troubleshootStyle.Render("API-related issues:")
		content += "\n"
		content += "• Check your internet connection"
		content += "\n"
		content += "• Verify your GEMINI_API_KEY environment variable is set correctly"
		content += "\n"
		content += "• Try again later as the API service might be temporarily unavailable"
		content += "\n"
		content += "• Check the API usage quota in your Google Cloud Console"
	} else if strings.Contains(strings.ToLower(m.errorMsg), "file") {
		content += troubleshootStyle.Render("File-related issues:")
		content += "\n"
		content += "• Check if the file exists at the specified path"
		content += "\n"
		content += "• Verify you have read permissions for the file"
		content += "\n"
		content += "• Try using an absolute path instead of a relative path"
	} else {
		content += troubleshootStyle.Render("General troubleshooting:")
		content += "\n"
		content += "• Try running the application again"
		content += "\n"
		content += "• Check the application logs for more details"
		content += "\n"
		content += "• Verify you have sufficient disk space and memory"
	}
	
	// What to try next
	content += "\n\n"
	content += infoStyle.Render("What to try next:")
	content += "\n"
	content += "1. Address the issue mentioned above"
	content += "\n"
	content += "2. Run the application again"
	content += "\n"
	content += "3. If the problem persists, try with simplified input"
	content += "\n\n"
	
	// Additional help
	content += "If you continue to experience issues, check the project documentation"
	content += "\n"
	content += "or report this problem in the GitHub repository."
	content += "\n\n"
	
	// Keyboard shortcuts
	content += keyboardHintStyle.Render("Press Enter to quit")
	
	return content
}