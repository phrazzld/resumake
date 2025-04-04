package tui

import (
	"fmt"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
)

// renderWelcomeView generates the welcome screen content
func renderWelcomeView(m Model) string {
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
	
	// Title with border
	welcomeTitle := StyledTitle("Welcome to Resumake!", true)
	
	// Application description
	description := "This tool helps you create a professional resume from your experience and qualifications."
	descriptionContent := subtitleStyle.Render(wrap(description, displayWidth-5))
	
	// How it works section
	howItWorksTitle := stepStyle.Render("How it works:")
	howItWorksContent := 
		"1. " + wrap("Optionally provide an existing resume to enhance", displayWidth-5) + "\n" +
		"2. " + wrap("Tell us about your experience, skills, and qualifications", displayWidth-5) + "\n" +
		"3. " + wrap("We'll generate a polished resume in Markdown format", displayWidth-5)
	
	howItWorksSection := lipgloss.JoinVertical(lipgloss.Left, 
		howItWorksTitle,
		howItWorksContent)
	
	// API key status section
	var apiStatusSection string
	if m.apiKeyOk {
		statusContent := successStyle.Render("✅ API key is valid and ready to use.") + "\n\n" +
			wrap("You're all set to create your professional resume!", displayWidth-5)
			
		apiStatusSection = successBoxStyle.Render(statusContent)
	} else {
		errorContent := errorStyle.Render("❌ API key is missing or invalid.") + "\n\n" +
			wrap("To use Resumake, you need to set the GEMINI_API_KEY environment variable:", displayWidth-5) + "\n" +
			"  export GEMINI_API_KEY=your_api_key_here\n\n" +
			wrap("You can get an API key from: https://makersuite.google.com/app/apikey", displayWidth-5) + "\n\n" +
			errorStyle.Render(wrap("Note: Proceeding without a valid API key will result in errors.", displayWidth-5))
			
		apiStatusSection = errorBoxStyle.Render(errorContent)
	}
	
	// Keyboard shortcuts box
	shortcutsMap := map[string]string{
		"Enter": "Continue to next step",
		"Ctrl+C": "Quit application",
		"Esc": "Go back (when available)",
	}
	shortcutsContent := KeyboardShortcuts(shortcutsMap)
	
	shortcutsTitle := keyboardHintStyle.Bold(true).Render("Keyboard shortcuts:")
	shortcutsSection := secondaryBoxStyle.Render(shortcutsTitle + "\n\n" + shortcutsContent)
	
	// Call to action
	callToAction := infoStyle.Render("Press Enter to continue...")
	
	// Join all sections vertically
	return lipgloss.JoinVertical(lipgloss.Left,
		welcomeTitle,
		"",
		descriptionContent,
		"",
		howItWorksSection,
		"",
		apiStatusSection,
		"",
		shortcutsSection,
		"",
		callToAction,
	)
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

// renderConfirmGenerateView generates the confirmation screen before generating the resume
func renderConfirmGenerateView(m Model) string {
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
	title := StyledTitle("Ready to generate your resume", true)
	
	// Summary section
	summaryContent := ""
	
	// Show source file info if it was provided
	if m.sourceContent != "" {
		summaryContent += "Source file: " + flagValueStyle.Render(m.sourcePathInput.Value()) + "\n"
		summaryContent += wrap(fmt.Sprintf("Source content: %d characters", len(m.sourceContent)), displayWidth-5) + "\n\n"
	}
	
	// Show input content info
	summaryContent += wrap(fmt.Sprintf("Input content: %d characters", len(m.stdinContent)), displayWidth-5) + "\n"
	
	// Show output path if it was provided via flags
	if m.flagOutputPath != "" {
		summaryContent += "\nOutput will be written to: " + pathStyle.Render(m.flagOutputPath)
	} else {
		summaryContent += "\nOutput will be written to the default path: " + pathStyle.Render("resume_out.md")
	}
	
	summarySection := primaryBoxStyle.Render(summaryContent)
	
	// Keyboard shortcuts
	shortcutsMap := map[string]string{
		"Enter": "Confirm and generate resume",
		"Esc": "Go back to edit details",
	}
	shortcutsContent := KeyboardShortcuts(shortcutsMap)
	shortcutsSection := secondaryBoxStyle.Render(shortcutsContent)
	
	// Call to action
	callToAction := infoStyle.Render("Press Enter to confirm and generate your resume...")
	
	// Join all sections vertically
	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		summarySection,
		"",
		shortcutsSection,
		"",
		callToAction,
	)
}

// renderErrorView generates the view shown when an error occurs
func renderErrorView(m Model) string {
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
	errorTitle := errorTitleStyle.Render("❌ ERROR: Resume Generation Failed ❌")
	
	// Error message section
	errorMessageTitle := "The following error occurred while generating your resume:"
	errorMessageContent := errorMsgStyle.Render(wrap(m.errorMsg, displayWidth-5))
	
	errorMessageSection := errorBoxStyle.Render(
		errorMessageTitle + "\n\n" + errorMessageContent)
	
	// Troubleshooting section
	troubleshootingTitle := titleStyle.Render("Troubleshooting")
	var troubleshootingContent string
	
	// Different troubleshooting suggestions based on error type
	errorLower := strings.ToLower(m.errorMsg)
	if strings.Contains(errorLower, "api") || strings.Contains(errorLower, "key") {
		troubleshootingContent += troubleshootStyle.Render("API-related issues:") + "\n"
		troubleshootingContent += wrap("• Check your internet connection", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Verify your GEMINI_API_KEY environment variable is set correctly", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Try again later as the API service might be temporarily unavailable", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Check the API usage quota in your Google Cloud Console", displayWidth-5)
	} else if strings.Contains(errorLower, "file") || strings.Contains(errorLower, "path") || strings.Contains(errorLower, "directory") {
		troubleshootingContent += troubleshootStyle.Render("File-related issues:") + "\n"
		troubleshootingContent += wrap("• Check if the file exists at the specified path", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Verify you have read permissions for the file", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Try using an absolute path instead of a relative path", displayWidth-5)
	} else if strings.Contains(errorLower, "timeout") || strings.Contains(errorLower, "context") || strings.Contains(errorLower, "deadline") {
		troubleshootingContent += troubleshootStyle.Render("Timeout-related issues:") + "\n"
		troubleshootingContent += wrap("• The request might have taken too long to complete", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Check your internet connection", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Try again with a shorter input", displayWidth-5)
	} else {
		troubleshootingContent += troubleshootStyle.Render("General troubleshooting:") + "\n"
		troubleshootingContent += wrap("• Try running the application again", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Check the application logs for more details", displayWidth-5) + "\n"
		troubleshootingContent += wrap("• Verify you have sufficient disk space and memory", displayWidth-5)
	}
	
	troubleshootingSection := secondaryBoxStyle.Render(
		troubleshootingTitle + "\n\n" + troubleshootingContent)
	
	// Next steps section
	nextStepsTitle := infoStyle.Render("What to try next:")
	nextStepsContent := 
		wrap("1. Address the issue mentioned above", displayWidth-5) + "\n" +
		wrap("2. Run the application again", displayWidth-5) + "\n" +
		wrap("3. If the problem persists, try with simplified input", displayWidth-5)
	
	nextStepsSection := secondaryBoxStyle.Render(
		nextStepsTitle + "\n\n" + nextStepsContent)
	
	// Additional help
	additionalHelp := wrap("If you continue to experience issues, check the project documentation or report this problem in the GitHub repository.", displayWidth-5)
	
	// Keyboard shortcuts
	keyboardHint := keyboardHintStyle.Render("Press Enter to quit")
	
	// Join all sections vertically
	return lipgloss.JoinVertical(lipgloss.Left,
		errorTitle,
		"",
		errorMessageSection,
		"",
		troubleshootingSection,
		"",
		nextStepsSection,
		"",
		additionalHelp,
		"",
		keyboardHint,
	)
}