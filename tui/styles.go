package tui

import (
	"strings"
	
	"github.com/charmbracelet/lipgloss"
)

// Define a consistent color palette with semantic meaning
var (
	// Primary brand colors
	primaryColor   = lipgloss.Color("#2C8CFF") // Vibrant blue - main brand color
	secondaryColor = lipgloss.Color("#15B097") // Teal - secondary brand color
	accentColor    = lipgloss.Color("#F2C94C") // Gold - attention-grabbing accent
	
	// Semantic colors
	successColor   = lipgloss.Color("#27AE60") // Green - success states
	errorColor     = lipgloss.Color("#EB5757") // Red - error states
	warningColor   = lipgloss.Color("#F2994A") // Orange - warning states
	
	// Neutral colors
	subtleColor    = lipgloss.Color("#BDBDBD") // Light gray - subtle text, hints
	neutralColor   = lipgloss.Color("#F2F2F2") // Off-white - neutral backgrounds
	darkColor      = lipgloss.Color("#333333") // Dark gray - text on light backgrounds
	bgColor        = lipgloss.Color("#121212") // Near black - dark backgrounds
)

// Base styles to be composed into more complex styles
var (
	// Base text styles
	baseStyle = lipgloss.NewStyle().
		Foreground(neutralColor)
	
	boldStyle = baseStyle.Copy().
		Bold(true)
	
	italicStyle = baseStyle.Copy().
		Italic(true)
)

// UI element styles
var (
	// Title styles
	titleStyle = boldStyle.Copy().
		Foreground(primaryColor).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		MarginBottom(1)
	
	subtitleStyle = italicStyle.Copy().
		Foreground(secondaryColor).
		MarginBottom(1)
	
	// Status styles
	successStyle = boldStyle.Copy().
		Foreground(successColor)
	
	errorStyle = boldStyle.Copy().
		Foreground(errorColor)
	
	warningStyle = boldStyle.Copy().
		Foreground(warningColor)
	
	infoStyle = boldStyle.Copy().
		Foreground(accentColor)
	
	// Keyboard hints
	keyboardHintStyle = italicStyle.Copy().
		Foreground(subtleColor)
	
	// Input styles
	inputLabelStyle = boldStyle.Copy().
		Foreground(primaryColor)
	
	inputFocusStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(0, 1)
	
	flagValueStyle = boldStyle.Copy().
		Foreground(accentColor)
	
	// Help text styles
	tipStyle = italicStyle.Copy().
		Foreground(secondaryColor)
	
	exampleStyle = baseStyle.Copy().
		Foreground(primaryColor)
	
	// Progress styles
	progressStyle = boldStyle.Copy().
		Foreground(accentColor)
	
	stepStyle = boldStyle.Copy().
		Foreground(primaryColor)
	
	completedStyle = boldStyle.Copy().
		Foreground(successColor)
	
	// Output path style
	pathStyle = boldStyle.Copy().
		Foreground(neutralColor).
		Background(darkColor).
		Padding(0, 1)
	
	// Error styles
	errorTitleStyle = boldStyle.Copy().
		Foreground(neutralColor).
		Background(errorColor).
		Padding(0, 1).
		MarginBottom(1)
	
	errorMsgStyle = boldStyle.Copy().
		Foreground(errorColor)
	
	troubleshootStyle = boldStyle.Copy().
		Foreground(primaryColor)
)

// Box styles for consistent containers
var (
	// Primary box - for main content sections
	primaryBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2)
	
	// Secondary box - for secondary content
	secondaryBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2)
	
	// Accent box - for important content that needs attention
	accentBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2)
	
	// Error box - for error messages
	errorBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
		Padding(1, 2)
	
	// Success box - for success messages
	successBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(successColor).
		Padding(1, 2)
)

// Utility functions for styled content

// StyledTitle creates a consistently styled title with optional border
func StyledTitle(title string, withBorder bool) string {
	if withBorder {
		return titleStyle.Copy().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 3).
			Render(title)
	}
	return titleStyle.Render(title)
}

// StyledSection creates a box with a title and content
func StyledSection(title string, content string, boxStyle lipgloss.Style) string {
	titleText := boldStyle.Copy().Foreground(primaryColor).Render(title)
	return boxStyle.Render(titleText + "\n\n" + content)
}

// KeyboardShortcuts formats a set of keyboard shortcuts consistently
func KeyboardShortcuts(shortcuts map[string]string) string {
	var lines []string
	for key, description := range shortcuts {
		lines = append(lines, boldStyle.Render(key+": ")+description)
	}
	
	return keyboardHintStyle.Copy().
		Bold(false).
		Render(strings.Join(lines, "\n"))
}