package api

import (
	"context"
	"errors"
	"fmt"

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
		return nil, fmt.Errorf("error generating content: %w", err)
	}

	// Check if response is valid
	if response == nil {
		return nil, errors.New("received nil response from API")
	}

	return response, nil
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
	
	// Check for generation errors in the first candidate
	if candidate.FinishReason != genai.FinishReasonStop && candidate.FinishReason != genai.FinishReasonUnspecified {
		return "", fmt.Errorf("generation did not complete successfully: %s", candidate.FinishReason)
	}

	// Check for content in the candidate
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", errors.New("no content in response")
	}

	// Extract the generated text
	return ParseGeneratedContent(candidate.Content)
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