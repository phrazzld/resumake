package tui

import (
	"strings"
	"testing"
)

func TestErrorAnalyzer(t *testing.T) {
	testCases := []struct {
		name                string
		errorMsg            string
		expectedCategory    string
		expectedHints       []string
		shouldContainDocRef bool
	}{
		{
			name:             "API authentication error",
			errorMsg:         "error executing API request: API authentication error: UNAUTHENTICATED: Invalid API key",
			expectedCategory: "API Authentication Error",
			expectedHints: []string{
				"Check your GEMINI_API_KEY environment variable is set correctly",
				"Verify your API key is valid and not expired",
				"Make sure you're using the correct API key format",
			},
			shouldContainDocRef: true,
		},
		{
			name:             "API quota exceeded",
			errorMsg:         "error executing API request: API quota or rate limit exceeded: RESOURCE_EXHAUSTED: Quota exceeded",
			expectedCategory: "API Quota Error",
			expectedHints: []string{
				"Wait a few minutes and try again",
				"Check if you've reached your API quota limit for the day",
				"Consider creating a new API key or upgrading your account",
			},
			shouldContainDocRef: true,
		},
		{
			name:             "Network error",
			errorMsg:         "error executing API request: network error while contacting API: deadline exceeded",
			expectedCategory: "Network Error",
			expectedHints: []string{
				"Check your internet connection",
				"Verify you can access the Gemini API (ping ai.google.dev)",
				"If using a proxy or VPN, try disabling it temporarily",
			},
			shouldContainDocRef: false,
		},
		{
			name:             "File not found error",
			errorMsg:         "failed to read source file: file does not exist: /path/to/nonexistent.md",
			expectedCategory: "File Error",
			expectedHints: []string{
				"Verify the file path is correct",
				"Check if the file exists in the specified location",
				"Make sure you have permission to read the file",
			},
			shouldContainDocRef: false,
		},
		{
			name:             "File size error",
			errorMsg:         "failed to read source file: file size exceeds the maximum allowed size of 10485760 bytes: /path/to/large.md",
			expectedCategory: "File Size Error",
			expectedHints: []string{
				"Your file exceeds the 10MB size limit",
				"Try splitting your content into smaller files",
				"Remove unnecessary content to reduce file size",
			},
			shouldContainDocRef: false,
		},
		{
			name:             "Write permission error",
			errorMsg:         "error writing output file: failed to write output: failed to write to file: permission denied",
			expectedCategory: "Write Permission Error",
			expectedHints: []string{
				"You don't have permission to write to the output location",
				"Try using a different output directory",
				"Run the application with higher privileges if appropriate",
			},
			shouldContainDocRef: false,
		},
		{
			name:             "Content truncation error",
			errorMsg:         "error processing API response: response was truncated because it reached maximum token limit",
			expectedCategory: "Content Truncation Error",
			expectedHints: []string{
				"Your input generated too much output",
				"Try simplifying your input or breaking it into smaller sections",
				"You can still use the partial output that was generated",
			},
			shouldContainDocRef: false,
		},
		{
			name:             "Safety filter error",
			errorMsg:         "error processing API response: Content was blocked due to safety filters",
			expectedCategory: "Safety Filter Error",
			expectedHints: []string{
				"Your content was flagged by the AI safety system",
				"Review your input for potentially sensitive or inappropriate content",
				"Try rephrasing any content that might be triggering safety filters",
			},
			shouldContainDocRef: true,
		},
		{
			name:             "Generic unrecognized error",
			errorMsg:         "an unknown error occurred: something went wrong",
			expectedCategory: "Error",
			expectedHints: []string{
				"Try running the command again",
				"Check the application logs for more details",
				"Restart the application and try again",
			},
			shouldContainDocRef: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			category, hints, docRef := analyzeError(tc.errorMsg)
			
			// Check category
			if category != tc.expectedCategory {
				t.Errorf("Expected category '%s', got '%s'", tc.expectedCategory, category)
			}
			
			// Check if all expected hints are present
			for _, expectedHint := range tc.expectedHints {
				found := false
				for _, actualHint := range hints {
					if actualHint == expectedHint {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected hint '%s' not found in hints: %v", expectedHint, hints)
				}
			}
			
			// Check doc reference
			if tc.shouldContainDocRef && docRef == "" {
				t.Errorf("Expected documentation reference but got none")
			} else if !tc.shouldContainDocRef && docRef != "" {
				t.Errorf("Did not expect documentation reference but got: %s", docRef)
			}
		})
	}
}

func TestEnhancedErrorView(t *testing.T) {
	// Test that the error view includes troubleshooting tips
	errorMsg := "error executing API request: API authentication error: UNAUTHENTICATED: Invalid API key"
	model := Model{errorMsg: errorMsg, width: 100, height: 40}
	
	errorView := renderErrorView(model)
	
	// The view should contain the error message
	if !containsString(errorView, "API authentication error") {
		t.Errorf("Error view should contain the error message")
	}
	
	// The view should contain a troubleshooting section
	if !containsString(errorView, "Troubleshooting") {
		t.Errorf("Error view should contain a troubleshooting section")
	}
	
	// The view should contain at least one hint from our expected list
	expectedHints := []string{
		"Check your GEMINI_API_KEY environment variable",
		"Verify your API key is valid",
		"API key format",
	}
	
	foundHint := false
	for _, hint := range expectedHints {
		if containsString(errorView, hint) {
			foundHint = true
			break
		}
	}
	
	if !foundHint {
		t.Errorf("Error view should contain at least one troubleshooting hint")
	}
}

// Helper function for string checking
func containsString(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}