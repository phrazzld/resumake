// Package prompt provides functionality for constructing API prompts.
//
// It handles the formatting and combining of different input sources
// (existing resume content and stdin input) into properly structured
// prompts suitable for sending to the Gemini API for resume generation.
package prompt

import (
	"github.com/google/generative-ai-go/genai"
)

// BuildPrompt combines existing resume content and user input into a formatted prompt string.
// It creates a clearly sectioned text that distinguishes between the existing resume 
// and the new user input, handling cases where either might be empty with appropriate 
// placeholder text.
//
// Parameters:
//   - sourceContent: Content from an existing resume file (can be empty)
//   - stdinContent: User input from stdin (can be empty)
//
// Returns:
//   - string: A formatted prompt string suitable for the Gemini API
//
// Example:
//
//	promptText := prompt.BuildPrompt(resumeContent, userInput)
//	fmt.Println("Generated prompt with length:", len(promptText))
func BuildPrompt(sourceContent, stdinContent string) string {
	var formattedPrompt string

	// Format existing resume section
	formattedPrompt = "EXISTING RESUME:\n"
	if sourceContent == "" {
		formattedPrompt += "(No existing resume provided)"
	} else {
		formattedPrompt += sourceContent
	}

	// Format user input section
	formattedPrompt += "\n\nUSER INPUT:\n"
	if stdinContent == "" {
		formattedPrompt += "(No additional input provided)"
	} else {
		formattedPrompt += stdinContent
	}

	return formattedPrompt
}

// GeneratePromptContent creates a genai.Content object from the source content and stdin input.
// This function builds on BuildPrompt but returns a structured Content object that
// can be used directly with the Gemini API's GenerateContent method.
//
// Parameters:
//   - sourceContent: Content from an existing resume file (can be empty)
//   - stdinContent: User input from stdin (can be empty)
//
// Returns:
//   - *genai.Content: A content object ready for sending to the Gemini API
//
// Example:
//
//	content := prompt.GeneratePromptContent(resumeContent, userInput)
//	response, err := model.GenerateContent(ctx, content.Parts...)
func GeneratePromptContent(sourceContent, stdinContent string) *genai.Content {
	promptText := BuildPrompt(sourceContent, stdinContent)
	
	return &genai.Content{
		Parts: []genai.Part{
			genai.Text(promptText),
		},
	}
}