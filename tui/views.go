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
	if displayWidth > 100 {
		displayWidth = 100 // Cap at 100 chars for readability
	}
	if displayWidth < 80 {
		displayWidth = 80 // Minimum width for the logo
	}
	
	// Container for our welcome screen
	docStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Bold(false).
		Width(displayWidth)
	
	// ASCII art logo
	logo := LogoText()
	
	// Version info (positioned below the logo)
	version := VersionInfo(m.appVersion)
	
	// Create a header with logo and version
	header := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		version,
	)
	
	// Tagline
	tagline := StyledTitle("Create Professional Resumes with AI", false, lipgloss.Center)
	
	// Application description with better styling
	description := "Resumake helps you create polished resumes from your experience and qualifications. " +
		"Powered by Google's Gemini AI, it transforms your input into a professional-looking document."
	
	descriptionStyle := subtitleStyle.Copy().
		Width(displayWidth - 20).
		Align(lipgloss.Center)
	
	descriptionContent := descriptionStyle.Render(description)
	
	// Features section
	featuresTitle := stepStyle.Render("Key Features:")
	
	featuresContent := 
		"• " + boldStyle.Render("AI-Powered Resume Generation") + ": Transform your experience into a professional resume\n" +
		"• " + boldStyle.Render("Enhance Existing Resumes") + ": Improve or update your current resume\n" +
		"• " + boldStyle.Render("Markdown Format") + ": Easy to edit or convert to other formats\n" +
		"• " + boldStyle.Render("Simple Terminal Interface") + ": No need for complex GUI applications"
	
	featuresBox := primaryBoxStyle.Copy().
		BorderForeground(secondaryColor).
		Width(displayWidth - 20).
		AlignHorizontal(lipgloss.Center).
		Render(featuresTitle + "\n\n" + featuresContent)
	
	// How it works section with clearer steps and icons
	howItWorksTitle := stepStyle.Render("How It Works:")
	
	step1 := successStyle.Render("Step 1 ") + boldStyle.Render("➤ Input") + 
		": Optionally provide an existing resume to enhance"
	step2 := successStyle.Render("Step 2 ") + boldStyle.Render("➤ Details") + 
		": Tell us about your experience, skills, and qualifications"
	step3 := successStyle.Render("Step 3 ") + boldStyle.Render("➤ Generate") + 
		": We'll use AI to create a polished resume in Markdown format"
	
	stepsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		wrap(step1, displayWidth-25),
		"",
		wrap(step2, displayWidth-25),
		"",
		wrap(step3, displayWidth-25),
	)
	
	stepsBox := accentBoxStyle.Copy().
		Width(displayWidth - 20).
		AlignHorizontal(lipgloss.Center).
		Render(howItWorksTitle + "\n\n" + stepsContent)
	
	// API key status section with clear visual indicator
	var apiStatusBox string
		
	if m.apiKeyOk {
		statusContent := lipgloss.JoinVertical(
			lipgloss.Left,
			successStyle.Render("✅  API KEY STATUS: READY"),
			"",
			wrap("Your Google Gemini API key is valid and ready to use.", displayWidth-25),
			wrap("You're all set to create your professional resume!", displayWidth-25),
		)
		
		apiStatusBox = successBoxStyle.Copy().
			Width(displayWidth - 20).
			AlignHorizontal(lipgloss.Center).
			Render(statusContent)
	} else {
		errorTitle := errorStyle.Render("❌  API KEY STATUS: MISSING")
		
		instructionsTitle := boldStyle.Render("To use Resumake, you need a Google Gemini API key:")
		instructions := "1. Visit: https://makersuite.google.com/app/apikey\n" +
			"2. Create a free API key\n" +
			"3. Set the environment variable:\n\n" +
			"   " + pathStyle.Render("export GEMINI_API_KEY=your_key_here")
		
		warningText := errorStyle.Render("Note: ") + 
			wrap("Proceeding without a valid API key will result in errors.", displayWidth-30)
		
		errorContent := lipgloss.JoinVertical(
			lipgloss.Left,
			errorTitle,
			"",
			instructionsTitle,
			"",
			instructions,
			"",
			warningText,
		)
		
		apiStatusBox = errorBoxStyle.Copy().
			Width(displayWidth - 20).
			AlignHorizontal(lipgloss.Center).
			Render(errorContent)
	}
	
	// Keyboard shortcuts in a cleaner format
	shortcutsTitle := keyboardHintStyle.Bold(true).Render("Keyboard Shortcuts")
	
	shortcutItems := []string{
		boldStyle.Render("Enter") + ": Continue to next step",
		boldStyle.Render("Ctrl+C") + ": Quit application",
		boldStyle.Render("Esc") + ": Go back (when available)",
	}
	
	shortcutsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		shortcutItems...,
	)
	
	shortcutsBox := secondaryBoxStyle.Copy().
		Width(displayWidth - 20).
		AlignHorizontal(lipgloss.Center).
		Render(shortcutsTitle + "\n\n" + shortcutsContent)
	
	// Call to action with more emphasis
	callToAction := infoStyle.Copy().
		Background(darkColor).
		Padding(0, 1).
		MarginTop(1).
		Render("Press Enter to begin...")
	
	// Compose the final document with all elements aligned
	return docStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			"",
			tagline,
			"",
			descriptionContent,
			"",
			featuresBox,
			"",
			stepsBox,
			"",
			apiStatusBox,
			"",
			shortcutsBox,
			"",
			callToAction,
		),
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
	title := StyledTitle("Ready to generate your resume", true, lipgloss.Left)
	
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