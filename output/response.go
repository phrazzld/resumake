package output

import (
	"errors"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// FinishReasonMessages maps finish reasons to descriptive messages
var FinishReasonMessages = map[genai.FinishReason]string{
	genai.FinishReasonStop:       "normal completion",
	genai.FinishReasonMaxTokens:  "maximum tokens reached",
	genai.FinishReasonSafety:     "content filtered due to safety settings",
	genai.FinishReasonRecitation: "content filtered due to repetition",
	genai.FinishReasonOther:      "unknown reason",
}

// ProcessResponseContent processes a Gemini API response and extracts valid Markdown content.
// It validates the response, extracts text, and ensures it's properly formatted Markdown.
// Returns the processed Markdown content and any error that occurred.
func ProcessResponseContent(response *genai.GenerateContentResponse) (string, error) {
	// Input validation
	if response == nil {
		return "", errors.New("response cannot be nil")
	}

	// Check for valid response content
	if len(response.Candidates) == 0 {
		return "", errors.New("no candidates in response")
	}

	candidate := response.Candidates[0]
	
	// Check for generation errors in the first candidate
	if candidate.FinishReason != genai.FinishReasonStop && 
	   candidate.FinishReason != genai.FinishReasonUnspecified {
		// Get a descriptive message for the finish reason
		reason := "unknown reason"
		if msg, ok := FinishReasonMessages[candidate.FinishReason]; ok {
			reason = msg
		}
		return "", fmt.Errorf("generation did not complete successfully: %s", reason)
	}

	// Check for content in the candidate
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", errors.New("no content in response")
	}

	// Extract the raw text from the response
	var rawText string
	hasText := false

	// Iterate through parts and concatenate all text
	for _, part := range candidate.Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			rawText += string(textPart)
			hasText = true
		}
	}

	if !hasText {
		return "", errors.New("no text content found in response")
	}

	// Process the extracted text
	return ExtractAndValidateMarkdown(rawText)
}

// ExtractAndValidateMarkdown extracts and validates Markdown content from raw text.
// It ensures the text contains valid Markdown syntax and formatting.
// Returns the validated Markdown content and any error that occurred.
func ExtractAndValidateMarkdown(responseText string) (string, error) {
	// Validate the text as Markdown
	if err := ValidateMarkdown(responseText); err != nil {
		return "", fmt.Errorf("invalid markdown content: %w", err)
	}
	
	// Prepare the content for output
	return PrepareForOutput(responseText)
}