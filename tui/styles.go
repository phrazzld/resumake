package tui

import (
	"strings"
	
	"github.com/charmbracelet/lipgloss"
)

// Define a consistent color palette with high contrast for both light and dark themes
var (
	// Primary brand colors with high contrast
	primaryColor   = lipgloss.AdaptiveColor{Light: "#0550AE", Dark: "#4C8FFF"} // Blue with good contrast in both modes
	secondaryColor = lipgloss.AdaptiveColor{Light: "#0B6E63", Dark: "#25D1B7"} // Teal with good contrast in both modes
	accentColor    = lipgloss.AdaptiveColor{Light: "#B07C00", Dark: "#FFCC3E"} // Gold with good contrast in both modes
	
	// Semantic colors with high contrast
	successColor   = lipgloss.AdaptiveColor{Light: "#1E6B38", Dark: "#4AE583"} // Green with good contrast in both modes
	errorColor     = lipgloss.AdaptiveColor{Light: "#AE1F3D", Dark: "#FF6B80"} // Red with good contrast in both modes
	
	// Neutral colors for text and backgrounds
	subtleColor    = lipgloss.AdaptiveColor{Light: "#777777", Dark: "#AAAAAA"} // Gray for subtle elements
	textColor      = lipgloss.AdaptiveColor{Light: "#222222", Dark: "#E8E8E8"} // Main text color
	bgAccentColor  = lipgloss.AdaptiveColor{Light: "#E8E8E8", Dark: "#333333"} // Slight contrast from background
	highlightColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"} // Maximum contrast
)

// Base styles to be composed into more complex styles
var (
	// Italic text style
	italicStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Italic(true)
)

// UI element styles
var (
	// Title styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		MarginBottom(1)
	
	// Status styles
	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(successColor)
	
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(errorColor)
	
	// Keyboard hints
	keyboardHintStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(subtleColor)
	
	// Help text styles
	tipStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(secondaryColor)
	
	// (Progress styles are defined inline in views.go)
	
	// Output path style - high contrast for important paths
	pathStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(textColor).
		Background(bgAccentColor).
		Padding(0, 1)
)

// Box styles for consistent containers
var (
	// Primary box - for main content sections
	primaryBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		BorderBackground(bgAccentColor)
	
	// Secondary box - for secondary content
	secondaryBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		BorderBackground(bgAccentColor)
	
	// Accent box - for important content that needs attention
	accentBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2).
		BorderBackground(bgAccentColor)
)

// Utility functions for styled content

// StyledTitle creates a consistently styled title with optional border and alignment
func StyledTitle(title string, withBorder bool, align lipgloss.Position) string {
	style := titleStyle
	
	if withBorder {
		style = style.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 3)
	}
	
	return style.
		AlignHorizontal(align).
		Render(title)
}

// StyledSection creates a box with a title and content
func StyledSection(title string, content string, boxStyle lipgloss.Style) string {
	titleText := lipgloss.NewStyle().Bold(true).Foreground(primaryColor).Render(title)
	return boxStyle.Render(titleText + "\n\n" + content)
}

// LogoText returns a stylized text-based logo for the application
func LogoText() string {
	// Create a high-contrast box for the logo to ensure visibility on any terminal
	logoBox := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Background(primaryColor).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 4).
		MarginBottom(1).
		Align(lipgloss.Center).
		Width(24)
	
	return logoBox.Render("R E S U M A K E")
}

// VersionInfo creates a version info tag for the welcome screen
func VersionInfo(version string) string {
	return lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true).
		Render("v" + version)
}

// KeyboardShortcuts formats a set of keyboard shortcuts consistently
func KeyboardShortcuts(shortcuts map[string]string) string {
	var lines []string
	for key, description := range shortcuts {
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render(key+": ")+description)
	}
	
	return keyboardHintStyle.
		Render(strings.Join(lines, "\n"))
}