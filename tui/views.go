package tui

import (
	"fmt"
	
	"github.com/charmbracelet/lipgloss"
)

// Helper function to constrain display width within reasonable bounds
func getConstrainedWidth(width int) int {
	// Set reasonable bounds for the width
	if width > 100 {
		width = 100 // Cap at 100 chars for readability
	}
	if width < 40 {
		width = 40 // Minimum width
	}
	return width
}

// renderWelcomeView generates the welcome screen content
func renderWelcomeView(m Model) string {
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Calculate display width
	displayWidth := getConstrainedWidth(m.width)
	
	// Container for our welcome screen
	docStyle := lipgloss.NewStyle().
		Width(displayWidth)
	
	// Logo text
	logo := LogoText()
	
	// Use inline styles with higher contrast
	titleText := lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		Background(bgAccentColor).
		Padding(1).
		Width(displayWidth-10).
		Align(lipgloss.Center).
		Render("Create Professional Resumes with AI")
		
	// API key status
	var apiStatus string
	if m.apiKeyOk {
		apiStatus = successStyle.Render("✓ API key is valid and ready to use")
	} else {
		apiStatus = errorStyle.Render("✗ API key is missing")
		apiStatus += "\n\n" + errorStyle.Render("To use Resumake, you need a Google Gemini API key")
		apiStatus += "\n" + pathStyle.Render("export GEMINI_API_KEY=your_key_here")
	}
	
	// Choose border color based on API key status
	var borderColor lipgloss.AdaptiveColor
	if m.apiKeyOk {
		borderColor = successColor
	} else {
		borderColor = errorColor
	}
	
	apiBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1).
		Width(displayWidth-20).
		Render(apiStatus)
		
	// Steps section
	stepsText := lipgloss.NewStyle().Bold(true).Render("How it works:") + "\n\n" +
		"1. " + wrap("Optionally provide an existing resume to enhance", displayWidth-20) + "\n\n" +
		"2. " + wrap("Tell us about your experience and skills", displayWidth-20) + "\n\n" +
		"3. " + wrap("Get your polished resume in markdown format", displayWidth-20)
	
	stepsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1).
		Width(displayWidth-20).
		Render(stepsText)
	
	// Call to action
	callToAction := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(accentColor).
		Padding(1).
		Render(" Press Enter to begin... ")
	
	// Join all elements vertically
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		titleText,
		"",
		apiBox,
		"",
		stepsBox,
		"",
		callToAction,
	)
	
	return docStyle.Render(content)
}

// renderSourceFileInputView generates the source file input view content
func renderSourceFileInputView(m Model) string {
	// Calculate display width (unused in this view currently)
	_ = getConstrainedWidth(m.width)
	
	// Create a title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(primaryColor).
		Padding(1).
		Render(" Source File Input ")
	
	// Create instructions with good contrast
	instructions := lipgloss.NewStyle().Bold(true).Render("Enter path to existing resume file (optional):")
	
	// Add text input view
	inputView := m.sourcePathInput.View()
	
	// Add hint about optional nature
	hint := italicStyle.Render("This step is optional. Press Enter to continue.")
	
	// Compose the view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		instructions,
		"",
		inputView,
		"",
		hint,
	)
}

// renderStdinInputView generates the stdin input view content
func renderStdinInputView(m Model) string {
	// Calculate display width (unused in this view currently)
	_ = getConstrainedWidth(m.width)
	
	// Create a title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(primaryColor).
		Padding(1).
		Render(" Enter Resume Details ")
	
	// Instructions
	instructions := lipgloss.NewStyle().Bold(true).Render("Tell us about your experience, skills, and qualifications:")
	
	// Add the text area view
	textareaView := m.stdinInput.View()
	
	// Add hint
	hint := italicStyle.Render("Press Ctrl+D when finished")
	
	// Compose the view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		instructions,
		"",
		textareaView,
		"",
		hint,
	)
}

// renderConfirmGenerateView generates the confirmation view before generating
func renderConfirmGenerateView(m Model) string {
	// Calculate display width (unused in this view currently)
	_ = getConstrainedWidth(m.width)
	
	// Create a title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(accentColor).
		Padding(1).
		Render(" Ready to Generate Resume ")
	
	// Show info
	info := "Press Enter to confirm and generate your resume"
	if m.sourceContent != "" {
		info = fmt.Sprintf("Source file: %s\n%s", m.sourcePathInput.Value(), info)
	}
	
	// Compose the view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		lipgloss.NewStyle().Bold(true).Render(info),
	)
}

// renderGeneratingView generates the view during resume generation
func renderGeneratingView(m Model) string {
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Calculate display width
	displayWidth := getConstrainedWidth(m.width)
	
	// Create a title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(primaryColor).
		Padding(1).
		Width(displayWidth - 4).
		Align(lipgloss.Center).
		Render("Generating Your Resume")
	
	// Calculate total characters of input
	totalChars := len(m.stdinContent) + len(m.sourceContent)
	
	// Create a spinner with enhanced style
	spinnerStyle := lipgloss.NewStyle().Bold(true).Foreground(accentColor)
	spinnerIcon := spinnerStyle.Render(m.spinner.View())
	
	// Create a progress indicator
	var progressIndicator string
	
	if m.progressStep != "" && m.progressMsg != "" {
		// Show specific progress steps when available
		stepTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(highlightColor).
			Background(accentColor).
			Padding(0, 1).
			Width(displayWidth - 10).
			Align(lipgloss.Center).
			Render("Step: " + m.progressStep)
		
		progressIndicator = lipgloss.JoinVertical(
			lipgloss.Center,
			stepTitle,
			"",
			wrap(m.progressMsg, displayWidth - 10),
		)
		
		// Put it in a nice box
		progressIndicator = secondaryBoxStyle.
			Width(displayWidth - 6).
			Render(progressIndicator)
	} else {
		// Default message when no specific progress is available
		progressIndicator = spinnerIcon + " " + lipgloss.NewStyle().Bold(true).Render("Processing your information...")
	}
	
	// Display input information
	inputInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Processing %d characters of input", totalChars)),
	)
	
	// Show source file info if provided
	if m.sourceContent != "" {
		sourceInfo := "Source file: " + m.sourcePathInput.Value()
		inputInfo = lipgloss.JoinVertical(
			lipgloss.Left,
			inputInfo,
			"",
			wrap(sourceInfo, displayWidth-8),
		)
	}
	
	// Create a styled input info box
	inputInfoBox := primaryBoxStyle.
		Width(displayWidth - 6).
		Render(inputInfo)
	
	// Show estimated time
	estimatedTime := tipStyle.Render(wrap("This may take up to 60 seconds depending on the input size.", displayWidth-8))
	
	// Additional information about the generation process
	processInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		wrap("The Gemini API is analyzing your experience and crafting a professional resume.", displayWidth-8),
		"",
		wrap("You'll be able to review and save the result when it's complete.", displayWidth-8),
	)
	
	// Create a styled process info box
	processInfoBox := accentBoxStyle.
		Width(displayWidth - 6).
		Render(processInfo)
	
	// Compose the complete view with all sections
	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		progressIndicator,
		"",
		inputInfoBox,
		"",
		estimatedTime,
		"",
		processInfoBox,
	)
}

// renderSuccessView generates the success view
func renderSuccessView(m Model) string {
	// Calculate display width (unused in this view currently)
	_ = getConstrainedWidth(m.width)
	
	// Create a title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(successColor).
		Padding(1).
		Render(" Success! ")
	
	// Show output path
	outputInfo := fmt.Sprintf("Resume saved to: %s", m.outputPath)
	
	// Compose the view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		successStyle.Render("✓ Resume successfully generated!"),
		"",
		lipgloss.NewStyle().Bold(true).Render(outputInfo),
		"",
		italicStyle.Render("Press Enter to quit"),
	)
}

// renderErrorView generates the error view
func renderErrorView(m Model) string {
	// Calculate display width (unused in this view currently)
	_ = getConstrainedWidth(m.width)
	
	// Create a title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(errorColor).
		Padding(1).
		Render(" Error ")
	
	// Show error message
	errorBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
		Padding(1).
		Render(errorStyle.Render(m.errorMsg))
	
	// Compose the view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		errorBox,
		"",
		italicStyle.Render("Press Enter to quit"),
	)
}