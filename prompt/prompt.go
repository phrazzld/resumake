package prompt

import (
	"github.com/google/generative-ai-go/genai"
)

// BuildPrompt combines existing resume content and user input into a formatted prompt string.
// It handles various scenarios such as missing resume or input content with appropriate defaults.
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
// This can be used directly with the Gemini API.
func GeneratePromptContent(sourceContent, stdinContent string) *genai.Content {
	promptText := BuildPrompt(sourceContent, stdinContent)
	
	return &genai.Content{
		Parts: []genai.Part{
			genai.Text(promptText),
		},
	}
}