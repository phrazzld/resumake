package tui

import (
	"fmt"
	
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
		Width(displayWidth)
	
	// Logo text
	logo := LogoText()
	
	// Use existing styles with higher contrast
	titleText := boldStyle.Copy().
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
	stepsText := boldStyle.Render("How it works:") + "\n\n" +
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
	callToAction := boldStyle.Copy().
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
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Create a title with high contrast
	title := boldStyle.Copy().
		Foreground(highlightColor).
		Background(primaryColor).
		Padding(1).
		Render(" Source File Input ")
	
	// Create instructions with good contrast
	instructions := boldStyle.Render("Enter path to existing resume file (optional):")
	
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
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Create a title with high contrast
	title := boldStyle.Copy().
		Foreground(highlightColor).
		Background(primaryColor).
		Padding(1).
		Render(" Enter Resume Details ")
	
	// Instructions
	instructions := boldStyle.Render("Tell us about your experience, skills, and qualifications:")
	
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
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Create a title with high contrast
	title := boldStyle.Copy().
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
		boldStyle.Render(info),
	)
}

// renderGeneratingView generates the view during generation
func renderGeneratingView(m Model) string {
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Create a title with high contrast
	title := boldStyle.Copy().
		Foreground(highlightColor).
		Background(primaryColor).
		Padding(1).
		Render(" Generating Resume ")
	
	// Show spinner
	spinnerText := m.spinner.View() + " Processing your information..."
	
	// Compose the view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		boldStyle.Render(spinnerText),
	)
}

// renderSuccessView generates the success view
func renderSuccessView(m Model) string {
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Create a title with high contrast
	title := boldStyle.Copy().
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
		boldStyle.Render(outputInfo),
		"",
		italicStyle.Render("Press Enter to quit"),
	)
}

// renderErrorView generates the error view
func renderErrorView(m Model) string {
	// Calculate display width
	displayWidth := m.width
	if displayWidth > 80 {
		displayWidth = 80 // Cap at 80 chars for readability
	}
	if displayWidth < 40 {
		displayWidth = 40 // Minimum width
	}
	
	// Create a title with high contrast
	title := boldStyle.Copy().
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