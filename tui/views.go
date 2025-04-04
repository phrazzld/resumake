package tui

import (
	"fmt"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
)

// Styles for different elements
var (
	// Define a consistent color palette with better contrast
	primaryColor   = lipgloss.Color("#2C8CFF") // Vibrant blue
	secondaryColor = lipgloss.Color("#15B097") // Teal
	accentColor    = lipgloss.Color("#F2C94C") // Gold
	successColor   = lipgloss.Color("#27AE60") // Green
	errorColor     = lipgloss.Color("#EB5757") // Red
	subtleColor    = lipgloss.Color("#BDBDBD") // Light gray
	neutralColor   = lipgloss.Color("#F2F2F2") // Off-white
	darkColor      = lipgloss.Color("#333333") // Dark gray
	
	// Title styles - make more prominent
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		MarginBottom(1)
	
	// Subtitle with better contrast
	subtitleStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(secondaryColor).
		MarginBottom(1)
	
	// Status styles
	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(successColor)
	
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(errorColor)
	
	infoStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
	
	// Keyboard hints with better visibility
	keyboardHintStyle = lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true)
	
	// Input styles
	inputLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor)
	
	flagValueStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
	
	// Help text styles
	tipStyle = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Italic(true)
	
	exampleStyle = lipgloss.NewStyle().
		Foreground(primaryColor)
	
	// Progress styles
	progressStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
	
	stepStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
	
	completedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(successColor)
	
	// Output path style
	pathStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(neutralColor).
		Background(darkColor).
		Padding(0, 1)
	
	// Error styles
	errorTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(neutralColor).
		Background(errorColor).
		Padding(0, 1).
		MarginBottom(1)
	
	errorMsgStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)
	
	troubleshootStyle = lipgloss.NewStyle().
		Foreground(primaryColor)
)

// renderWelcomeView generates the welcome screen content
func renderWelcomeView(m Model) string {
	var content string
	
	// Title - using prettier borders with proper styling
	welcomeTitle := titleStyle.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 3).
		Render("Welcome to Resumake!")
	
	content += welcomeTitle
	content += "\n\n"
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Application description - now with text wrapping
	descWidth := m.width
	if descWidth > 80 {
		descWidth = 80 // Cap at 80 chars for readability
	}
	
	description := "This tool helps you create a professional resume from your experience and qualifications."
	content += subtitleStyle.Render(wrap(description, descWidth-5)) // -5 for margin
	content += "\n\n"
	
	// How it works section with better styling
	content += stepStyle.Render("How it works:") + "\n"
	content += "1. " + wrap("Optionally provide an existing resume to enhance", descWidth-5) + "\n"
	content += "2. " + wrap("Tell us about your experience, skills, and qualifications", descWidth-5) + "\n"
	content += "3. " + wrap("We'll generate a polished resume in Markdown format", descWidth-5) + "\n\n"
	
	// API key status with better styling
	if m.apiKeyOk {
		content += successStyle.Render("✅ API key is valid and ready to use.")
		content += "\n\n"
		content += wrap("You're all set to create your professional resume!", descWidth-5)
	} else {
		content += errorStyle.Render("❌ API key is missing or invalid.")
		content += "\n\n"
		content += wrap("To use Resumake, you need to set the GEMINI_API_KEY environment variable:", descWidth-5) + "\n"
		content += "  export GEMINI_API_KEY=your_api_key_here\n\n"
		content += wrap("You can get an API key from: https://makersuite.google.com/app/apikey", descWidth-5) + "\n"
		content += errorStyle.Render(wrap("Note: Proceeding without a valid API key will result in errors.", descWidth-5))
	}
	
	// Keyboard shortcuts in a nice box
	content += "\n\n"
	shortcutsTitle := keyboardHintStyle.Bold(true).Render("Keyboard shortcuts:")
	shortcuts := "• Enter: Continue to next step\n• Ctrl+C: Quit application\n• Esc: Go back (when available)"
	
	shortcutsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtleColor).
		Padding(1, 2).
		Render(shortcutsTitle + "\n" + shortcuts)
		
	content += shortcutsBox
	content += "\n\n"
	
	// Call to action
	content += infoStyle.Render("Press Enter to continue...")
	
	return content
}

// renderSourceFileInputView generates the source file input view content
func renderSourceFileInputView(m Model) string {
	var content string
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Title
	content += titleStyle.Render("Source File Input")
	content += "\n\n"
	
	// Instructions and explanation
	content += wrap("You can optionally provide an existing resume file to use as a starting point.", displayWidth-5)
	content += "\n"
	content += wrap("This helps generate a better result by incorporating your existing content.", displayWidth-5)
	content += "\n\n"
	
	// Show if path was provided via flags
	if m.flagSourcePath != "" {
		content += infoStyle.Render("A file path was provided from command line flags:")
		content += "\n"
		content += flagValueStyle.Render(wrap(m.flagSourcePath, displayWidth-5))
		content += "\n\n"
		content += wrap("You can edit this path or leave it as is.", displayWidth-5)
		content += "\n\n"
	}
	
	// Input field label
	content += inputLabelStyle.Render("Enter path to existing resume file (optional):")
	content += "\n"
	
	// The actual text input component
	content += m.sourcePathInput.View()
	content += "\n\n"
	
	// Note about being optional
	content += wrap("Note: This step is optional. If you don't have an existing resume file, just press Enter to continue without selecting a file.", displayWidth-5)
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
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Title
	content += titleStyle.Render("Resume Details")
	content += "\n\n"
	
	// Instructions
	content += subtitleStyle.Render(wrap("Tell us about your experience, skills, and qualifications.", displayWidth-5))
	content += "\n\n"
	
	// Show source file info if one was provided
	if m.sourceContent != "" {
		content += "Source file: " + successStyle.Render(wrap(m.sourcePathInput.Value(), displayWidth-20))
		content += "\n"
		content += wrap(fmt.Sprintf("Content length: %d characters", len(m.sourceContent)), displayWidth-5)
		content += "\n\n"
		content += wrap("We'll combine this source content with the details you provide below.", displayWidth-5)
		content += "\n\n"
	}
	
	// Input field label
	content += inputLabelStyle.Render("Enter your resume details:")
	content += "\n\n"
	
	// Tips and examples
	content += tipStyle.Render("Tips:")
	content += "\n"
	content += tipStyle.Render(wrap("• Be specific about your achievements", displayWidth-10))
	content += "\n"
	content += tipStyle.Render(wrap("• Include relevant skills and technologies", displayWidth-10))
	content += "\n"
	content += tipStyle.Render(wrap("• Mention education, certifications, and years of experience", displayWidth-10))
	content += "\n\n"
	
	content += exampleStyle.Render("Example:")
	content += "\n"
	content += exampleStyle.Render(wrap("I have 5 years of experience as a software developer.", displayWidth-10))
	content += "\n"
	content += exampleStyle.Render(wrap("Skills: Go, Python, JavaScript, Docker, Kubernetes", displayWidth-10))
	content += "\n"
	content += exampleStyle.Render(wrap("Education: BS in Computer Science, University of Example", displayWidth-10))
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
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Title with more prominent styling
	generatingTitle := titleStyle.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Render("Generating Your Resume")
	
	content += generatingTitle
	content += "\n\n"
	
	// Create a styled box for the spinner and processing info
	spinnerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2).
		Render(m.spinner.View() + " " + progressStyle.Render("Processing your information"))
	
	content += spinnerBox
	content += "\n\n"
	
	// Input information with wrapping
	totalChars := len(m.stdinContent) + len(m.sourceContent)
	content += wrap(fmt.Sprintf("Processing %d characters of input...", totalChars), displayWidth-5)
	content += "\n"
	
	// Source file info if provided
	if m.sourceContent != "" {
		sourceInfo := "Source file: " + m.sourcePathInput.Value()
		content += wrap(sourceInfo, displayWidth-5)
		content += "\n"
	}
	
	// Estimated time with better styling
	content += "\n"
	timeEstimate := "This may take up to 60 seconds depending on the input size."
	content += infoStyle.Render(wrap(timeEstimate, displayWidth-5))
	content += "\n\n"
	
	// Progress information with proper wrapping
	if m.progressStep != "" && m.progressMsg != "" {
		// Create a progress box
		progressContent := stepStyle.Render("Step: " + m.progressStep) + "\n" +
			wrap(m.progressMsg, displayWidth-10)
		
		progressBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2).
			Render(progressContent)
			
		content += progressBox
		content += "\n\n"
	} else {
		content += wrap("Please wait while we generate your resume...", displayWidth-5)
		content += "\n\n"
	}
	
	// Additional status info with proper wrapping
	statusMsg1 := "The Gemini API is analyzing your experience and crafting a professional resume."
	statusMsg2 := "You'll be able to review and save the result when it's complete."
	
	content += wrap(statusMsg1, displayWidth-5)
	content += "\n"
	content += wrap(statusMsg2, displayWidth-5)
	
	return content
}

// renderSuccessView generates the view shown after successful resume generation
func renderSuccessView(m Model) string {
	var content string
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Title with success indicator
	content += completedStyle.Render("✅ Success! Resume Generation Complete! ✅")
	content += "\n\n"
	
	// Border for output path
	content += wrap("Your resume has been successfully generated and saved to:", displayWidth-5)
	content += "\n"
	content += pathStyle.Render(wrap(m.outputPath, displayWidth-10))
	content += "\n\n"
	
	// Show result information
	content += wrap(fmt.Sprintf("Content length: %s characters", m.resultMessage), displayWidth-5)
	content += "\n\n"
	
	// Next steps with a box
	content += titleStyle.Render("╔═════════════════════════╗")
	content += "\n"
	content += titleStyle.Render("║ Next Steps              ║")
	content += "\n"
	content += titleStyle.Render("╚═════════════════════════╝")
	content += "\n\n"
	
	// Detailed next steps list
	content += successStyle.Render("1. Review your resume at ") + wrap(m.outputPath, displayWidth-30)
	content += "\n"
	content += successStyle.Render(wrap("2. Make any necessary edits or refinements", displayWidth-5))
	content += "\n"
	content += successStyle.Render(wrap("3. Convert to other formats:", displayWidth-5))
	content += "\n"
	content += wrap("   • PDF: Use a markdown converter like pandoc", displayWidth-5)
	content += "\n"
	content += wrap("   • DOCX: Use a markdown converter or import into your word processor", displayWidth-5)
	content += "\n\n"
	
	// Congratulatory message
	content += infoStyle.Render(wrap("Congratulations on creating your professional resume!", displayWidth-5))
	content += "\n\n"
	
	// Keyboard shortcuts
	content += keyboardHintStyle.Render("Press Enter to quit")
	
	return content
}

// renderErrorView generates the view shown when an error occurs
func renderErrorView(m Model) string {
	var content string
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Title with error indicator
	content += errorTitleStyle.Render("❌ ERROR: Resume Generation Failed ❌")
	content += "\n\n"
	
	// Error message
	content += "The following error occurred while generating your resume:"
	content += "\n\n"
	content += errorMsgStyle.Render(wrap(m.errorMsg, displayWidth-5))
	content += "\n\n"
	
	// Troubleshooting section
	content += titleStyle.Render("Troubleshooting")
	content += "\n\n"
	
	// Different troubleshooting suggestions based on error type
	if strings.Contains(strings.ToLower(m.errorMsg), "api") {
		content += troubleshootStyle.Render("API-related issues:")
		content += "\n"
		content += wrap("• Check your internet connection", displayWidth-5)
		content += "\n"
		content += wrap("• Verify your GEMINI_API_KEY environment variable is set correctly", displayWidth-5)
		content += "\n"
		content += wrap("• Try again later as the API service might be temporarily unavailable", displayWidth-5)
		content += "\n"
		content += wrap("• Check the API usage quota in your Google Cloud Console", displayWidth-5)
	} else if strings.Contains(strings.ToLower(m.errorMsg), "file") {
		content += troubleshootStyle.Render("File-related issues:")
		content += "\n"
		content += wrap("• Check if the file exists at the specified path", displayWidth-5)
		content += "\n"
		content += wrap("• Verify you have read permissions for the file", displayWidth-5)
		content += "\n"
		content += wrap("• Try using an absolute path instead of a relative path", displayWidth-5)
	} else {
		content += troubleshootStyle.Render("General troubleshooting:")
		content += "\n"
		content += wrap("• Try running the application again", displayWidth-5)
		content += "\n"
		content += wrap("• Check the application logs for more details", displayWidth-5)
		content += "\n"
		content += wrap("• Verify you have sufficient disk space and memory", displayWidth-5)
	}
	
	// What to try next
	content += "\n\n"
	content += infoStyle.Render("What to try next:")
	content += "\n"
	content += wrap("1. Address the issue mentioned above", displayWidth-5)
	content += "\n"
	content += wrap("2. Run the application again", displayWidth-5)
	content += "\n"
	content += wrap("3. If the problem persists, try with simplified input", displayWidth-5)
	content += "\n\n"
	
	// Additional help
	content += wrap("If you continue to experience issues, check the project documentation or report this problem in the GitHub repository.", displayWidth-5)
	content += "\n\n"
	
	// Keyboard shortcuts
	content += keyboardHintStyle.Render("Press Enter to quit")
	
	return content
}