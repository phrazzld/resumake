package tui

import (
	"strings"
	"testing"
)

func TestEnhancedSuccessView(t *testing.T) {
	// Create a model with success information
	model := Model{
		state:         stateResultSuccess,
		outputPath:    "/tmp/resume_out.md",
		resultMessage: "2500", // Content length
		stdinContent:  "Sample content for resume",
		sourceContent: "Sample source content",
		width:         80,
		height:        24,
	}
	
	// Get the rendered view
	successView := renderSuccessView(model)
	
	// TEST 1: Check for essential elements
	essentialElements := []string{
		"Success",                   // Title element
		"/tmp/resume_out.md",        // Output path
		"2500 characters",           // Content length with label
		"Enter",                     // Quit instruction
	}
	
	for _, element := range essentialElements {
		if !strings.Contains(successView, element) {
			t.Errorf("Success view should contain '%s'", element)
		}
	}
	
	// TEST 2: Check for next steps section
	if !strings.Contains(successView, "Next Steps") {
		t.Errorf("Success view should contain a 'Next Steps' section")
	}
	
	// TEST 3: Check for formatting options
	formatOptions := []string{"Markdown", "PDF", "DOCX", "HTML"}
	formatOptionsFound := 0
	
	for _, format := range formatOptions {
		if strings.Contains(successView, format) {
			formatOptionsFound++
		}
	}
	
	if formatOptionsFound < 2 {
		t.Errorf("Success view should mention at least 2 formatting options")
	}
	
	// TEST 4: Check for celebratory elements
	celebratorySymbols := []string{"✓", "✅", "🎉", "🚀", "💯", "congratulations", "Congratulations"}
	celebratoryFound := false
	
	for _, symbol := range celebratorySymbols {
		if strings.Contains(successView, symbol) {
			celebratoryFound = true
			break
		}
	}
	
	if !celebratoryFound {
		t.Errorf("Success view should include celebratory elements")
	}
	
	// TEST 5: Check for input information
	if model.sourceContent != "" && !strings.Contains(strings.ToLower(successView), "source file") {
		t.Errorf("Success view should mention source file when one was used")
	}
	
	// TEST 6: Check for proper structure with empty lines between sections
	lines := strings.Split(successView, "\n")
	emptyLineCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			emptyLineCount++
		}
	}
	if emptyLineCount < 3 {
		t.Errorf("Success view should have proper spacing between sections (at least 3 empty lines)")
	}
	
	// TEST 7: Test with different output paths
	// Test with absolute path
	modelAbsPath := model
	modelAbsPath.outputPath = "/absolute/path/to/resume_out.md"
	viewAbsPath := renderSuccessView(modelAbsPath)
	if !strings.Contains(viewAbsPath, "/absolute/path/to/resume_out.md") {
		t.Errorf("Success view should display absolute paths correctly")
	}
	
	// Test with relative path
	modelRelPath := model
	modelRelPath.outputPath = "./relative/path/resume_out.md"
	viewRelPath := renderSuccessView(modelRelPath)
	if !strings.Contains(viewRelPath, "./relative/path/resume_out.md") {
		t.Errorf("Success view should display relative paths correctly")
	}
}