package api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

// ModelInterface defines the minimal interface needed for executing API requests
type ModelInterface interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
	SetMaxOutputTokens(tokens int32)
	SetTemperature(temp float32)
}

// ExecuteRequest sends the provided content to the Gemini API and returns the response.
// It requires a valid model, content, and a context for the request.
func ExecuteRequest(ctx context.Context, model ModelInterface, content *genai.Content) (*genai.GenerateContentResponse, error) {
	// Input validation
	if model == nil {
		return nil, errors.New("model cannot be nil")
	}
	if content == nil {
		return nil, errors.New("content cannot be nil")
	}

	// Set generation parameters
	model.SetMaxOutputTokens(8192)
	model.SetTemperature(0.7) // Balanced between creativity and determinism

	// Make the API request
	fmt.Println("Sending request to Gemini API...")
	response, err := model.GenerateContent(ctx, content.Parts...)
	if err != nil {
		// Parse the error to provide more detailed information
		return nil, handleAPIError(err)
	}

	// Check if response is valid
	if response == nil {
		return nil, errors.New("received nil response from API")
	}

	return response, nil
}

// handleAPIError parses API errors and returns more user-friendly messages
// with potential solutions when possible.
func handleAPIError(err error) error {
	errorMsg := err.Error()
	
	// Handle quota exceeded errors
	if strings.Contains(errorMsg, "RESOURCE_EXHAUSTED") || 
	   strings.Contains(errorMsg, "Quota exceeded") ||
	   strings.Contains(errorMsg, "rate limit") {
		return fmt.Errorf("API quota or rate limit exceeded: %w. "+
			"Please wait a few minutes and retry, or check your quota management settings", err)
	}
	
	// Handle authentication errors
	if strings.Contains(errorMsg, "UNAUTHENTICATED") || 
	   strings.Contains(errorMsg, "API key") ||
	   strings.Contains(errorMsg, "authentication") {
		return fmt.Errorf("API authentication error: %w. "+
			"Please verify your GEMINI_API_KEY environment variable is correct and valid", err)
	}
	
	// Handle network/timeout errors
	if strings.Contains(errorMsg, "deadline exceeded") ||
	   strings.Contains(errorMsg, "connection") ||
	   strings.Contains(errorMsg, "network") {
		return fmt.Errorf("network error while contacting API: %w. "+
			"Please check your internet connection and try again", err)
	}
	
	// Handle invalid request errors
	if strings.Contains(errorMsg, "INVALID_ARGUMENT") {
		return fmt.Errorf("invalid request to API: %w. "+
			"Please check the format of your prompt", err)
	}
	
	// Default case for unrecognized errors
	return fmt.Errorf("error generating content: %w", err)
}

// ProcessResponse extracts and processes the text from the API response.
// Returns the generated text and any error that occurred.
func ProcessResponse(response *genai.GenerateContentResponse) (string, error) {
	// Input validation
	if response == nil {
		return "", errors.New("response cannot be nil")
	}

	// Check for valid response content
	if len(response.Candidates) == 0 {
		return "", errors.New("no candidates in response")
	}

	candidate := response.Candidates[0]
	
	// Handle specific error conditions
	if candidate.FinishReason != genai.FinishReasonStop && candidate.FinishReason != genai.FinishReasonUnspecified {
		switch candidate.FinishReason {
		case genai.FinishReasonSafety:
			return handleSafetyError(candidate)
		case genai.FinishReasonMaxTokens:
			return "", fmt.Errorf("response was truncated because it reached maximum token limit; try simplifying your input")
		case genai.FinishReasonRecitation:
			return "", fmt.Errorf("response was filtered due to content repetition; try adding more variation to your input")
		default:
			return "", fmt.Errorf("generation did not complete successfully: %s", candidate.FinishReason)
		}
	}

	// Check for content in the candidate
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", errors.New("no content in response")
	}

	// Extract the generated text
	return ParseGeneratedContent(candidate.Content)
}

// handleSafetyError processes safety-related errors and provides detailed information
// about which safety policies were triggered and how to address them.
func handleSafetyError(candidate *genai.Candidate) (string, error) {
	// Start with a base error message
	errMsg := "Content was blocked due to safety filters"
	
	// Add details about specific safety categories if available
	if len(candidate.SafetyRatings) > 0 {
		errMsg += ". Safety categories flagged:"
		
		for i, rating := range candidate.SafetyRatings {
			if rating.Probability >= genai.HarmProbabilityHigh {
				if i > 0 {
					errMsg += ","
				}
				errMsg += fmt.Sprintf(" %s (probability: %s)", 
					formatHarmCategory(rating.Category), 
					formatHarmProbability(rating.Probability))
			}
		}
	}
	
	// Add guidance on how to address the issue
	errMsg += ". Consider reviewing your input for potentially sensitive or inappropriate content."
	
	return "", errors.New(errMsg)
}

// formatHarmCategory converts a HarmCategory to a human-readable string
func formatHarmCategory(category genai.HarmCategory) string {
	switch category {
	case genai.HarmCategoryHarassment:
		return "Harassment"
	case genai.HarmCategoryHateSpeech:
		return "Hate Speech"
	case genai.HarmCategoryDangerous:
		return "Dangerous Content"
	case genai.HarmCategorySexuallyExplicit:
		return "Sexually Explicit Content"
	default:
		return fmt.Sprintf("Category %d", category)
	}
}

// formatHarmProbability converts a HarmProbability to a human-readable string
func formatHarmProbability(probability genai.HarmProbability) string {
	switch probability {
	case genai.HarmProbabilityHigh:
		return "High"
	case genai.HarmProbabilityMedium:
		return "Medium"
	case genai.HarmProbabilityLow:
		return "Low"
	case genai.HarmProbabilityNegligible:
		return "Negligible"
	case genai.HarmProbabilityUnspecified:
		return "Unspecified"
	default:
		return fmt.Sprintf("%d", probability)
	}
}

// TryRecoverPartialContent attempts to extract usable content from a truncated response.
// It adds a warning annotation but allows the user to see the partial content.
func TryRecoverPartialContent(response *genai.GenerateContentResponse) (string, error) {
	// Input validation
	if response == nil {
		return "", errors.New("response cannot be nil")
	}

	// Check for valid response content
	if len(response.Candidates) == 0 {
		return "", errors.New("no candidates in response")
	}

	candidate := response.Candidates[0]
	
	// We can only recover from MaxTokens truncation
	if candidate.FinishReason != genai.FinishReasonMaxTokens {
		return "", fmt.Errorf("can only recover partial content from token limit truncation, not %s", candidate.FinishReason)
	}
	
	// Check for content in the candidate
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", errors.New("no content in response")
	}
	
	// Extract whatever content exists
	rawText, err := ParseGeneratedContent(candidate.Content)
	if err != nil {
		return "", fmt.Errorf("failed to extract partial content: %w", err)
	}
	
	// Append a note about truncation
	warning := "\n\n---\n\n**Note: This content was truncated due to reaching the maximum token limit. The resume may be incomplete.**"
	
	return rawText + warning, nil
}

// ParseGeneratedContent extracts text from the content parts of a response.
// Returns the concatenated text from all text parts.
func ParseGeneratedContent(content *genai.Content) (string, error) {
	if content == nil {
		return "", errors.New("content cannot be nil")
	}

	var result string
	hasText := false

	// Iterate through parts and concatenate all text
	for _, part := range content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			result += string(textPart)
			hasText = true
		}
	}

	if !hasText {
		return "", errors.New("no text content found in response")
	}

	return result, nil
}