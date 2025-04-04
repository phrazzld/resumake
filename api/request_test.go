package api

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/generative-ai-go/genai"
)

// MockModel is a test double that satisfies the GenerativeModel interface needed for testing
type MockGenerativeModel struct {
	generateContentFunc func(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
	callCount           int
}

// GenerateContent implements the necessary method for the test interface
func (m *MockGenerativeModel) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
	m.callCount++
	if m.generateContentFunc != nil {
		return m.generateContentFunc(ctx, parts...)
	}
	return nil, errors.New("not implemented")
}

// SetMaxOutputTokens is a mock implementation
func (m *MockGenerativeModel) SetMaxOutputTokens(tokens int32) {}

// SetTemperature is a mock implementation
func (m *MockGenerativeModel) SetTemperature(temp float32) {}

// ExecuteRequestInterface is a minimal interface for our mock to implement
type ExecuteRequestInterface interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
	SetMaxOutputTokens(tokens int32)
	SetTemperature(temp float32)
}

func TestExecuteRequest(t *testing.T) {
	// This test will verify the basic functionality of executing an API request
	t.Run("Successfully execute API request", func(t *testing.T) {
		// Setup will use a mock model so no actual API calls are made
		// But we need to assert that the function correctly calls the model's GenerateContent
		// And returns the response
		ctx := context.Background()
		mockModel := &MockGenerativeModel{
			generateContentFunc: func(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
				return &genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{
							Content: &genai.Content{
								Parts: []genai.Part{
									genai.Text("Generated resume content"),
								},
							},
							FinishReason: genai.FinishReasonStop,
						},
					},
				}, nil
			},
		}
		
		content := &genai.Content{
			Parts: []genai.Part{
				genai.Text("Test prompt"),
			},
		}
		
		// Cast to genai.GenerativeModel
		var model ExecuteRequestInterface = mockModel
		resp, err := ExecuteRequest(ctx, model, content)
		if err != nil {
			t.Errorf("ExecuteRequest() returned error: %v", err)
		}
		if resp == nil {
			t.Error("ExecuteRequest() returned nil response")
		}
		if mockModel.callCount != 1 {
			t.Errorf("Expected GenerateContent to be called 1 time, was called %d times", mockModel.callCount)
		}
	})

	// Test error cases
	t.Run("Handle nil model", func(t *testing.T) {
		// Test that the function properly handles a nil model input
		// Should return an error without panicking
		ctx := context.Background()
		content := &genai.Content{
			Parts: []genai.Part{
				genai.Text("Test prompt"),
			},
		}
		
		_, err := ExecuteRequest(ctx, nil, content)
		if err == nil {
			t.Error("ExecuteRequest() with nil model should return error, got nil")
		}
	})

	t.Run("Handle nil content", func(t *testing.T) {
		// Test that the function properly handles a nil content input
		// Should return an error without panicking
		ctx := context.Background()
		mockModel := &MockGenerativeModel{}
		
		var model ExecuteRequestInterface = mockModel
		_, err := ExecuteRequest(ctx, model, nil)
		if err == nil {
			t.Error("ExecuteRequest() with nil content should return error, got nil")
		}
	})

	t.Run("Handle API error", func(t *testing.T) {
		// Test that the function properly handles an error from the API
		// Should propagate the error to the caller
		ctx := context.Background()
		expectedErr := errors.New("API error")
		mockModel := &MockGenerativeModel{
			generateContentFunc: func(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
				return nil, expectedErr
			},
		}
		
		content := &genai.Content{
			Parts: []genai.Part{
				genai.Text("Test prompt"),
			},
		}
		
		var model ExecuteRequestInterface = mockModel
		_, err := ExecuteRequest(ctx, model, content)
		if err == nil {
			t.Error("ExecuteRequest() should propagate API errors, got nil")
		}
	})
}

func TestProcessResponse(t *testing.T) {
	// This test will verify that we can extract the generated text from the response
	t.Run("Extract text from valid response", func(t *testing.T) {
		// Setup a valid response with known text content
		expectedText := "Generated resume content"
		response := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{
					Content: &genai.Content{
						Parts: []genai.Part{
							genai.Text(expectedText),
						},
					},
					FinishReason: genai.FinishReasonStop,
				},
			},
		}
		
		text, err := ProcessResponse(response)
		if err != nil {
			t.Errorf("ProcessResponse() returned error: %v", err)
		}
		if text != expectedText {
			t.Errorf("ProcessResponse() = %q, want %q", text, expectedText)
		}
	})

	t.Run("Handle empty response", func(t *testing.T) {
		// Test how the function handles an empty response
		// Should return an error or empty string depending on design
		response := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{},
		}
		
		_, err := ProcessResponse(response)
		if err == nil {
			t.Error("ProcessResponse() with empty candidates should return error, got nil")
		}
	})

	t.Run("Handle nil response", func(t *testing.T) {
		// Test that the function properly handles a nil response
		// Should return an error without panicking
		_, err := ProcessResponse(nil)
		if err == nil {
			t.Error("ProcessResponse() with nil response should return error, got nil")
		}
	})

	t.Run("Handle response with no text", func(t *testing.T) {
		// Test how the function handles a response with no text parts
		// Should return an error or empty string depending on design
		response := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{
					Content: &genai.Content{
						Parts: []genai.Part{},
					},
					FinishReason: genai.FinishReasonStop,
				},
			},
		}
		
		_, err := ProcessResponse(response)
		if err == nil {
			t.Error("ProcessResponse() with no text parts should return error, got nil")
		}
	})
}

func TestParseGeneratedContent(t *testing.T) {
	// This test will verify that we can parse the generated content from the API response
	t.Run("Parse markdown from response", func(t *testing.T) {
		// Setup a response with markdown content
		markdownContent := "# Resume\n\n## Skills\n\n- Go\n- Python"
		content := &genai.Content{
			Parts: []genai.Part{
				genai.Text(markdownContent),
			},
		}
		
		text, err := ParseGeneratedContent(content)
		if err != nil {
			t.Errorf("ParseGeneratedContent() returned error: %v", err)
		}
		if text != markdownContent {
			t.Errorf("ParseGeneratedContent() = %q, want %q", text, markdownContent)
		}
	})

	t.Run("Handle multiple content parts", func(t *testing.T) {
		// Test how the function handles responses with multiple content parts
		// Should concatenate or select the appropriate parts
		part1 := "First part"
		part2 := "Second part"
		expected := part1 + part2
		
		content := &genai.Content{
			Parts: []genai.Part{
				genai.Text(part1),
				genai.Text(part2),
			},
		}
		
		text, err := ParseGeneratedContent(content)
		if err != nil {
			t.Errorf("ParseGeneratedContent() returned error: %v", err)
		}
		if text != expected {
			t.Errorf("ParseGeneratedContent() = %q, want %q", text, expected)
		}
	})
	
	t.Run("Handle nil content", func(t *testing.T) {
		_, err := ParseGeneratedContent(nil)
		if err == nil {
			t.Error("ParseGeneratedContent() with nil content should return error, got nil")
		}
	})
	
	t.Run("Handle content with no text parts", func(t *testing.T) {
		content := &genai.Content{
			Parts: []genai.Part{},
		}
		
		_, err := ParseGeneratedContent(content)
		if err == nil {
			t.Error("ParseGeneratedContent() with no text parts should return error, got nil")
		}
	})
}

// TestHandleSafetyRatings verifies that responses with safety rating issues 
// are properly detected and reported
func TestHandleSafetyRatings(t *testing.T) {
	// Setup a response with safety rating issues
	response := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text("Some content that was flagged"),
					},
				},
				FinishReason: genai.FinishReasonSafety,
				SafetyRatings: []*genai.SafetyRating{
					{
						Category:    genai.HarmCategoryHarassment,
						Probability: genai.HarmProbabilityHigh,
					},
				},
			},
		},
	}
	
	// Process the response
	_, err := ProcessResponse(response)
	
	// Verify error handling
	if err == nil {
		t.Error("ProcessResponse() should return error for safety-blocked content")
	}
	
	// Verify error message content
	errorMsg := err.Error()
	if !stringContainsAny(errorMsg, []string{"safety", "content policy", "blocked", "filtered"}) {
		t.Errorf("Error message should mention safety issues: %s", errorMsg)
	}
	
	// Verify that the error message contains information about the specific safety category
	if !stringContainsAny(errorMsg, []string{"harassment", "Harassment"}) {
		t.Errorf("Error message should mention specific safety categories: %s", errorMsg)
	}
}

// TestHandleRateLimitErrors verifies that rate limit or quota errors 
// are properly detected and reported
func TestHandleRateLimitErrors(t *testing.T) {
	// Setup a mock model that returns a quota exceeded error
	ctx := context.Background()
	mockModel := &MockGenerativeModel{
		generateContentFunc: func(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
			return nil, errors.New("RESOURCE_EXHAUSTED: Quota exceeded for quota metric 'GenerationRequests'")
		},
	}
	
	content := &genai.Content{
		Parts: []genai.Part{
			genai.Text("Test prompt"),
		},
	}
	
	// Cast to the interface
	var model ExecuteRequestInterface = mockModel
	
	// Attempt to execute the request
	_, err := ExecuteRequest(ctx, model, content)
	
	// Verify error handling
	if err == nil {
		t.Error("ExecuteRequest() should return error for quota exceeded")
	}
	
	// Verify error message content
	errorMsg := err.Error()
	if !stringContainsAny(errorMsg, []string{"API quota", "rate limit", "Quota exceeded"}) {
		t.Errorf("Error message should indicate quota or rate limit issues: %s", errorMsg)
	}
	
	// Verify that the error message includes advice
	if !stringContainsAny(errorMsg, []string{"retry", "wait", "quota management"}) {
		t.Errorf("Error message should include advice on addressing quota issues: %s", errorMsg)
	}
}

// TestRecoverFromPartialResponse verifies that we can extract usable content
// from truncated responses when possible
func TestRecoverFromPartialResponse(t *testing.T) {
	// Setup a response with a truncated completion
	partialContent := "# Resume\n\n## Skills\n\n- Go\n- Python\n\n## Experience\n\n- Software Engineer at"
	response := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(partialContent),
					},
				},
				FinishReason: genai.FinishReasonMaxTokens,
			},
		},
	}
	
	// Process the response
	content, err := TryRecoverPartialContent(response)
	
	// Verify we get usable content with a warning
	if err != nil {
		t.Errorf("TryRecoverPartialContent() should not return an error, got: %v", err)
	}
	
	// Should return the partial content
	if content == "" {
		t.Error("TryRecoverPartialContent() should return partial content even if truncated")
	}
	
	// Content should have a warning about being truncated
	if !strings.Contains(content, partialContent) {
		t.Error("Recovered content should contain the original partial content")
	}
	
	if !strings.Contains(content, "Note: This content was truncated") {
		t.Error("Recovered content should include a warning about truncation")
	}
}

// Helper function to check if a string contains any of the provided substrings
func stringContainsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}