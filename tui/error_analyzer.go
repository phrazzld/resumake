package tui

import (
	"strings"
)

// Error categories
const (
	// API-related errors
	categoryAPIAuth       = "API Authentication Error"
	categoryAPIQuota      = "API Quota Error"
	categoryAPINetwork    = "Network Error"
	categoryAPISafety     = "Safety Filter Error"
	categoryAPITruncation = "Content Truncation Error"
	
	// File-related errors
	categoryFileNotFound  = "File Error"
	categoryFileSize      = "File Size Error"
	categoryFilePermission = "File Permission Error"
	
	// Output-related errors
	categoryWritePermission = "Write Permission Error"
	categoryDirError        = "Directory Error"
	
	// Generic error
	categoryGeneric = "Error"
)

// Documentation references
const (
	apiDocRef = "For API issues, visit: https://ai.google.dev/docs/api_errors"
	geminiDocsRef = "Gemini API documentation: https://ai.google.dev/docs"
)

// analyzeError examines the error message and returns:
// 1. A category to help the user understand what went wrong
// 2. Specific troubleshooting hints based on the error type
// 3. Optional documentation reference (if available)
func analyzeError(errorMsg string) (category string, hints []string, docRef string) {
	// Default to generic category
	category = categoryGeneric
	
	// Default to generic hints
	hints = []string{
		"Try running the command again",
		"Check the application logs for more details",
		"Restart the application and try again",
	}
	
	// Now check for specific error patterns, starting with API errors
	
	// API authentication errors
	if containsAny(errorMsg, []string{
		"API authentication error", 
		"UNAUTHENTICATED", 
		"Invalid API key",
		"API key not valid",
	}) {
		category = categoryAPIAuth
		hints = []string{
			"Check your GEMINI_API_KEY environment variable is set correctly",
			"Verify your API key is valid and not expired",
			"Make sure you're using the correct API key format",
		}
		docRef = apiDocRef
		return
	}
	
	// API quota or rate limit errors
	if containsAny(errorMsg, []string{
		"quota or rate limit exceeded",
		"RESOURCE_EXHAUSTED",
		"Quota exceeded",
		"rate limit",
	}) {
		category = categoryAPIQuota
		hints = []string{
			"Wait a few minutes and try again",
			"Check if you've reached your API quota limit for the day",
			"Consider creating a new API key or upgrading your account",
		}
		docRef = apiDocRef
		return
	}
	
	// Network errors
	if containsAny(errorMsg, []string{
		"network error",
		"deadline exceeded",
		"connection",
		"timeout",
	}) {
		category = categoryAPINetwork
		hints = []string{
			"Check your internet connection",
			"Verify you can access the Gemini API (ping ai.google.dev)",
			"If using a proxy or VPN, try disabling it temporarily",
		}
		return
	}
	
	// Safety filter errors
	if containsAny(errorMsg, []string{
		"safety filters",
		"Content was blocked",
		"safety categories flagged",
		"HarmCategory",
	}) {
		category = categoryAPISafety
		hints = []string{
			"Your content was flagged by the AI safety system",
			"Review your input for potentially sensitive or inappropriate content",
			"Try rephrasing any content that might be triggering safety filters",
		}
		docRef = geminiDocsRef
		return
	}
	
	// Content truncation errors
	if containsAny(errorMsg, []string{
		"truncated",
		"maximum token limit",
		"token limit",
		"MaxTokens",
	}) {
		category = categoryAPITruncation
		hints = []string{
			"Your input generated too much output",
			"Try simplifying your input or breaking it into smaller sections",
			"You can still use the partial output that was generated",
		}
		return
	}
	
	// File not found errors
	if containsAny(errorMsg, []string{
		"file does not exist",
		"no such file",
		"could not find file",
	}) {
		category = categoryFileNotFound
		hints = []string{
			"Verify the file path is correct",
			"Check if the file exists in the specified location",
			"Make sure you have permission to read the file",
		}
		return
	}
	
	// File size errors
	if containsAny(errorMsg, []string{
		"file size exceeds",
		"maximum allowed size",
		"file too large",
	}) {
		category = categoryFileSize
		hints = []string{
			"Your file exceeds the 10MB size limit",
			"Try splitting your content into smaller files",
			"Remove unnecessary content to reduce file size",
		}
		return
	}
	
	// File permission errors
	if containsAny(errorMsg, []string{
		"error accessing file",
		"permission denied",
		"cannot read file",
	}) && !strings.Contains(errorMsg, "write") {
		category = categoryFilePermission
		hints = []string{
			"You don't have permission to read the file",
			"Check the file permissions (try 'ls -l' on the file)",
			"Try running the application with appropriate permissions",
		}
		return
	}
	
	// Write permission errors
	if containsAny(errorMsg, []string{
		"error writing output file",
		"failed to write",
		"permission denied",
		"cannot write",
	}) {
		category = categoryWritePermission
		hints = []string{
			"You don't have permission to write to the output location",
			"Try using a different output directory",
			"Run the application with higher privileges if appropriate",
		}
		return
	}
	
	// Directory errors
	if containsAny(errorMsg, []string{
		"directory exists but is not a directory",
		"failed to create directory",
		"failed to check directory",
	}) {
		category = categoryDirError
		hints = []string{
			"There's an issue with the output directory",
			"Make sure the parent directory exists and is writable",
			"Try specifying a different output location",
		}
		return
	}
	
	// Return the defaults for any other error
	return
}

// containsAny checks if the string contains any of the patterns
func containsAny(s string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(strings.ToLower(s), strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}