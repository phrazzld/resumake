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
// It handles the entire pipeline from raw API response to validated Markdown:
// 1. Validates the response structure for completeness
// 2. Checks for error conditions like safety filters or truncation
// 3. Extracts raw text from the response
// 4. Validates and cleans the Markdown
//
// Parameters:
//   - response: The raw response from the Gemini API
//
// Returns:
//   - string: The processed, validated, and cleaned Markdown content
//   - error: Any error encountered during processing
//
// Example:
//
//	markdownContent, err := output.ProcessResponseContent(apiResponse)
//	if err != nil {
//	    log.Fatalf("Failed to process API response: %v", err)
//	}
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
// It serves as a bridge between the raw text extraction from API responses and
// the Markdown validation and cleaning functionality. This ensures that the text
// is properly formatted as Markdown before being used as output.
//
// Parameters:
//   - responseText: The raw text extracted from the API response
//
// Returns:
//   - string: The validated and cleaned Markdown content
//   - error: Any error encountered during validation or preparation
//
// Example:
//
//	markdown, err := output.ExtractAndValidateMarkdown(rawText)
//	if err != nil {
//	    log.Fatalf("Invalid markdown in response: %v", err)
//	}
func ExtractAndValidateMarkdown(responseText string) (string, error) {
	// Validate the text as Markdown
	if err := ValidateMarkdown(responseText); err != nil {
		return "", fmt.Errorf("invalid markdown content: %w", err)
	}
	
	// Prepare the content for output
	return PrepareForOutput(responseText)
}