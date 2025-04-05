package tui

import (
	"fmt"
	"strings"
	
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

// renderSourceFileInputView generates the enhanced source file input view content
func renderSourceFileInputView(m Model) string {
	// Calculate display width
	displayWidth := getConstrainedWidth(m.width)
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Create a centered title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(primaryColor).
		Padding(1).
		Width(displayWidth - 4).
		Align(lipgloss.Center).
		Render("📄 Source File Input")
	
	// Create a description section explaining the purpose
	description := wrap(
		"Provide an existing resume file to enhance. Resumake will use this as a " +
		"starting point to generate an improved version with better formatting and content.",
		displayWidth - 8)
	
	// Build the instructions section with examples and flag indication
	instructionsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Instructions")
	
	// Create instructions content
	instructionsContent := "Enter the path to your existing resume file:"
	
	// Add flag path indication if provided
	if m.flagSourcePath != "" {
		flagInfo := lipgloss.NewStyle().
			Foreground(accentColor).
			Render("• Pre-filled from command line flags: " + m.flagSourcePath)
		instructionsContent += "\n\n" + flagInfo
	}
	
	// Display the input field with focus-aware styling
	inputContent := m.sourcePathInput.View()
	var styledInputView string
	
	// Apply different styling based on focus state
	if m.sourcePathInput.Focused() {
		styledInputView = FocusedStyle(inputContent, displayWidth - 8)
	} else {
		styledInputView = UnfocusedStyle(inputContent, displayWidth - 8)
	}
	
	// Create a helpful tips section
	tipsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Helpful Tips")
	
	tipsContent := "• This step is optional. Press Enter to continue without a source file\n" +
		"• Supported file formats: .txt, .md, .markdown\n" +
		"• Example path: /home/user/documents/my_resume.md or ./resume.txt\n" +
		"• Maximum file size: 10MB\n" +
		"• Using a source file can significantly improve the quality of your generated resume"
	
	// If terminal is narrow, wrap the tips content
	tipsContent = wrap(tipsContent, displayWidth - 12)
	
	// Keyboard shortcuts section
	shortcutsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Keyboard Shortcuts")
	
	shortcutsContent := "• Enter: Continue to next step\n" +
		"• Ctrl+C: Quit application"
	
	// Put instructions, input, and shortcuts in a main content box
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		instructionsTitle,
		"",
		instructionsContent,
		"",
		styledInputView,
		"",
		shortcutsTitle,
		"",
		shortcutsContent,
	)
	
	mainContentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Width(displayWidth - 4).
		Render(mainContent)
	
	// Put tips in a separate tips box
	tipsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Width(displayWidth - 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			tipsTitle,
			"",
			tipsContent,
		))
	
	// Compose the complete view
	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		lipgloss.NewStyle().Width(displayWidth - 8).Render(description),
		"",
		mainContentBox,
		"",
		tipsBox,
	)
}

// renderStdinInputView generates the enhanced stdin input view content
func renderStdinInputView(m Model) string {
	// Calculate display width
	displayWidth := getConstrainedWidth(m.width)
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Create a centered title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(primaryColor).
		Padding(1).
		Width(displayWidth - 4).
		Align(lipgloss.Center).
		Render("✏️ Enter Resume Details")
	
	// Create a description section explaining the purpose
	description := wrap(
		"Tell us about your professional background. The more details you provide, "+
		"the better your resume will be. Include your experience, skills, education, and achievements.",
		displayWidth - 8)
	
	// Build the instructions section
	instructionsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Instructions")
	
	instructionsContent := wrap(
		"Type your resume content below. Include as much detail as possible about your "+
		"professional experience and skills. The AI will structure and enhance this information.",
		displayWidth - 16)
	
	// Style for the textarea container with focus-aware styling
	textareaContent := m.stdinInput.View()
	var styledTextareaView string
	
	// Apply different styling based on focus state
	if m.stdinInput.Focused() {
		styledTextareaView = FocusedStyle(textareaContent, displayWidth - 8)
	} else {
		styledTextareaView = UnfocusedStyle(textareaContent, displayWidth - 8)
	}
	
	// Create a suggestions section
	suggestionsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Suggestions:")
	
	suggestionsContent := "• Work Experience: Company names, positions, dates, and key responsibilities\n" +
		"• Skills: Technical, soft, and domain-specific skills\n" +
		"• Education: Degrees, institutions, graduation dates\n" +
		"• Achievements: Awards, certifications, projects\n" +
		"• Use bullet points for better readability\n" +
		"• Highlight metrics and results when possible (e.g., 'increased sales by 20%')"
	
	// If terminal is narrow, wrap the suggestions content
	suggestionsContent = wrap(suggestionsContent, displayWidth - 12)
	
	// Create a formatting examples section
	examplesTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Example Format:")
	
	examplesContent := wrap(
		"Work Experience:\n"+
		"- Senior Software Engineer at XYZ Corp (2019-2023)\n"+
		"- Led a team of 5 developers to deliver a new product feature\n"+
		"- Reduced system latency by 40% through code optimization\n\n"+
		"Skills: JavaScript, React, Node.js, Project Management\n\n"+
		"Education: BS Computer Science, University of Technology (2015)",
		displayWidth - 12)
	
	// Keyboard shortcuts section
	shortcutsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Keyboard Shortcuts")
	
	shortcutsContent := "• Arrow keys: Navigate within the text\n" +
		"• Enter: Add new line\n" +
		"• Tab: Indent text\n" +
		"• Ctrl+D: Finish input and continue\n" +
		"• Ctrl+C: Quit application"
	
	// If terminal is narrow, wrap the shortcuts content
	shortcutsContent = wrap(shortcutsContent, displayWidth - 12)
	
	// Build the main content with input area
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		instructionsTitle,
		"",
		instructionsContent,
		"",
		// Add the text area view with focus-aware styling
		styledTextareaView,
		"",
		shortcutsTitle,
		"",
		shortcutsContent,
	)
	
	// Style the main content box
	mainContentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Width(displayWidth - 4).
		Render(mainContent)
	
	// Create suggestions and examples box
	tipsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		suggestionsTitle,
		"",
		suggestionsContent,
		"",
		examplesTitle,
		"",
		examplesContent,
	)
	
	// Style the tips box
	tipsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Width(displayWidth - 4).
		Render(tipsContent)
	
	// Compose the complete view
	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		lipgloss.NewStyle().Width(displayWidth - 8).Render(description),
		"",
		mainContentBox,
		"",
		tipsBox,
	)
}

// renderConfirmGenerateView generates the confirmation view before generating
func renderConfirmGenerateView(m Model) string {
	// Calculate display width
	displayWidth := getConstrainedWidth(m.width)
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Create a centered title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(accentColor).
		Padding(1).
		Width(displayWidth - 4).
		Align(lipgloss.Center).
		Render("🚀 Ready to Generate Resume")
	
	// Create a summary section
	summaryTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Summary of Input")
	
	// Build summary content
	var summaryContent strings.Builder
	
	// Add source file info if provided
	if m.sourceContent != "" {
		sourceInfo := fmt.Sprintf("📄 Source file: %s", m.sourcePathInput.Value())
		summaryContent.WriteString(wrap(sourceInfo, displayWidth - 16) + "\n\n")
	}
	
	// Add input content summary (truncated)
	inputLength := len(m.stdinContent)
	if inputLength > 0 {
		var contentPreview string
		if inputLength > 100 {
			contentPreview = m.stdinContent[:97] + "..."
		} else {
			contentPreview = m.stdinContent
		}
		
		contentInfo := fmt.Sprintf("✏️ Input: %d characters\n\n", inputLength)
		summaryContent.WriteString(contentInfo)
		summaryContent.WriteString(wrap("Preview: "+contentPreview, displayWidth - 16))
	}
	
	// Add output path info if provided via flags
	if m.flagOutputPath != "" {
		outputInfo := fmt.Sprintf("\n\n📁 Output path: %s", m.flagOutputPath)
		summaryContent.WriteString(wrap(outputInfo, displayWidth - 16))
	}
	
	// Build the summary box
	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Width(displayWidth - 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			summaryTitle,
			"",
			summaryContent.String(),
		))
	
	// Add confirmatation instruction
	instruction := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Render("Press Enter to confirm and generate your resume")
	
	// Add hint about ESC
	hint := italicStyle.Render("Press ESC to go back and edit your input")
	
	// Compose the complete view
	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		summaryBox,
		"",
		instruction,
		"",
		hint,
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

// renderSuccessView generates the enhanced success view with celebratory elements
func renderSuccessView(m Model) string {
	// Calculate display width
	displayWidth := getConstrainedWidth(m.width)
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Create a celebratory title with high contrast
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(successColor).
		Padding(1).
		Width(displayWidth - 4).
		Align(lipgloss.Center).
		Render("🎉 Success! 🎉")
	
	// Create a celebratory message
	celebrationMsg := lipgloss.NewStyle().
		Bold(true).
		Foreground(successColor).
		Align(lipgloss.Center).
		Width(displayWidth - 4).
		Render("✅ Your professional resume has been successfully generated!")
	
	// Create a stats section
	// Parse the content length
	contentLength := "Unknown"
	if m.resultMessage != "" {
		contentLength = m.resultMessage + " characters"
	}
	
	// Calculate input stats
	sourceFileInfo := ""
	if m.sourceContent != "" {
		sourceFile := m.sourcePathInput.Value()
		sourceFileInfo = fmt.Sprintf("📄 Source file: %s\n\n", sourceFile)
	}
	
	// Build statistics section
	statsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("📊 Resume Stats")
	
	statsContent := fmt.Sprintf("%s📏 Size: %s\n\n⏱️ Generated in seconds", sourceFileInfo, contentLength)
	
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(successColor).
		Padding(1, 2).
		Width(displayWidth - 10).
		Render(statsTitle + "\n\n" + statsContent)
	
	// Output path with clear formatting and highlighting
	outputPathTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("📂 Output Location")
	
	pathText := fmt.Sprintf("Your resume is saved at:\n\n%s", 
		lipgloss.NewStyle().
			Background(bgAccentColor).
			Padding(0, 1).
			Render(m.outputPath))
	
	outputPathBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2).
		Width(displayWidth - 10).
		Render(outputPathTitle + "\n\n" + pathText)
	
	// Next steps guidance
	nextStepsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("🚀 Next Steps")
	
	nextStepsContent := "1. Your resume is in Markdown format (.md)\n\n" +
		"2. You can convert it to other formats:\n" +
		"   • PDF: Use a markdown editor or online converter\n" +
		"   • DOCX: Import to Word or Google Docs\n" +
		"   • HTML: Use a markdown to HTML converter\n\n" +
		"3. Review and customize before sending to employers"
	
	nextStepsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Width(displayWidth - 10).
		Render(nextStepsTitle + "\n\n" + wrap(nextStepsContent, displayWidth - 20))
	
	// Exit instructions
	exitInstructions := italicStyle.Render("Press Enter to quit or run again")
	
	// Compose the view with all sections
	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		celebrationMsg,
		"",
		statsBox,
		"",
		outputPathBox,
		"",
		nextStepsBox,
		"",
		exitInstructions,
	)
}

// renderErrorView generates the error view with contextual troubleshooting
func renderErrorView(m Model) string {
	// Calculate display width
	displayWidth := getConstrainedWidth(m.width)
	
	// Use the shared wrapText utility for consistent text wrapping
	wrap := func(text string, width int) string {
		return wrapText(text, width)
	}
	
	// Analyze the error to determine the category and troubleshooting hints
	category, hints, docRef := analyzeError(m.errorMsg)
	
	// Create a title with high contrast that includes the error category
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(errorColor).
		Padding(1).
		Render(" Error: " + category + " ")
	
	// Show error message with consistent wrapping
	errorBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
		Padding(1, 2).
		Width(displayWidth - 4).
		Render(errorStyle.Render(wrap(m.errorMsg, displayWidth - 10)))
	
	// Create a troubleshooting box with hints
	troubleshootingTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Render("Troubleshooting")
	
	// Build the hints section
	var hintsContent strings.Builder
	for i, hint := range hints {
		if i > 0 {
			hintsContent.WriteString("\n\n")
		}
		hintsContent.WriteString("• " + hint)
	}
	
	// Add doc reference if available
	if docRef != "" {
		hintsContent.WriteString("\n\n" + italicStyle.Render(docRef))
	}
	
	troubleshootingBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Width(displayWidth - 4).
		Render(troubleshootingTitle + "\n\n" + hintsContent.String())
	
	// Compose the view with all sections
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		errorBox,
		"",
		troubleshootingBox,
		"",
		italicStyle.Render("Press Enter to quit"),
	)
}