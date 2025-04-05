package tui

import (
	"strings"
)

// wrapText wraps text at the specified width to ensure it fits in the terminal
// It handles word wrapping, respecting word boundaries where possible
func wrapText(text string, width int) string {
	if width <= 0 {
		width = 80 // Default to 80 if we don't have width info
	}
	
	// Simple word wrapping
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	
	var lines []string
	currentLine := ""
	
	for _, word := range words {
		// Handle words longer than the width by breaking them
		if len(word) > width {
			// If we have content on the current line, add it to lines and start fresh
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = ""
			}
			
			// Split the long word into chunks
			for len(word) > 0 {
				if len(word) <= width {
					// Last piece fits on its own line
					lines = append(lines, word)
					word = ""
				} else {
					// Take a width-sized chunk and continue
					lines = append(lines, word[:width])
					word = word[width:]
				}
			}
		} else if len(currentLine)+len(word)+1 > width && currentLine != "" {
			// Word would exceed line width, start a new line
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			// Add word to current line with space if needed
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}
	
	// Add the last line if it has content
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return strings.Join(lines, "\n")
}