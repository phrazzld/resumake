// Package api provides functionality for interacting with the Gemini API.
//
// It handles API key retrieval, client initialization, request execution,
// response processing, and error handling. The package abstracts away the
// complexities of direct API interaction while providing detailed error
// information and recovery mechanisms for common failure scenarios.
package api

import (
	"context"
	"errors"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// DefaultModelName is the identifier for the specific Gemini model version used.
// This model is optimized for resume generation with strong text formatting capabilities.
const DefaultModelName = "gemini-2.5-pro-exp-03-25"

// SystemInstructions defines the system instructions for the resume generation model.
// These instructions guide the model to generate professional resumes in Markdown format
// based on the user's input, without fabricating information not present in the input.
const SystemInstructions = `You are an expert resume writing assistant. Your goal is to synthesize the provided existing resume information (if any) and the raw stream-of-consciousness input into a single, coherent, professional resume formatted strictly in Markdown.

Prioritize clarity, conciseness, and professional language. Structure the output logically with clear headings (e.g., Summary, Experience, Projects, Skills, Education). Infer structure and dates where possible, but do not fabricate information not present in the inputs. Focus on elevating the user's actual experience. Ensure the final output is only Markdown content.`

// GetAPIKey retrieves the Gemini API key from the GEMINI_API_KEY environment variable.
// This key is required for authenticating with the Gemini API.
//
// Returns:
//   - string: The API key if found
//   - error: An error if the environment variable is not set or is empty
//
// Example:
//
//	apiKey, err := api.GetAPIKey()
//	if err != nil {
//	    log.Fatal("API key not found, please set GEMINI_API_KEY environment variable")
//	}
func GetAPIKey() (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", errors.New("GEMINI_API_KEY environment variable is required")
	}
	return apiKey, nil
}

// InitializeClient creates a new Gemini client with the provided API key 
// and initializes it with the default model (DefaultModelName).
// It also configures the model with system instructions for resume generation.
//
// Parameters:
//   - ctx: The context for API requests
//   - apiKey: The Gemini API key for authentication
//
// Returns:
//   - *genai.Client: The initialized API client
//   - *genai.GenerativeModel: The configured model instance
//   - error: Any error that occurred during initialization
//
// Example:
//
//	ctx := context.Background()
//	client, model, err := api.InitializeClient(ctx, apiKey)
//	if err != nil {
//	    log.Fatalf("Failed to initialize API client: %v", err)
//	}
//	defer client.Close()
func InitializeClient(ctx context.Context, apiKey string) (*genai.Client, *genai.GenerativeModel, error) {
	return InitializeClientWithModel(ctx, apiKey, DefaultModelName)
}

// InitializeClientWithModel creates a new Gemini client with the provided API key
// and initializes it with the specified model name. This function allows using
// a different model than the default one, which can be useful for testing or
// when newer models become available.
//
// Parameters:
//   - ctx: The context for API requests
//   - apiKey: The Gemini API key for authentication
//   - modelName: The specific Gemini model identifier to use
//
// Returns:
//   - *genai.Client: The initialized API client
//   - *genai.GenerativeModel: The configured model instance
//   - error: Any error that occurred during initialization
//
// Example:
//
//	ctx := context.Background()
//	client, model, err := api.InitializeClientWithModel(ctx, apiKey, "gemini-3.0-pro")
//	if err != nil {
//	    log.Fatalf("Failed to initialize API client: %v", err)
//	}
//	defer client.Close()
func InitializeClientWithModel(ctx context.Context, apiKey string, modelName string) (*genai.Client, *genai.GenerativeModel, error) {
	// Validate API key
	if apiKey == "" {
		return nil, nil, errors.New("API key cannot be empty")
	}

	// Initialize client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, nil, err
	}

	// Get model
	model := client.GenerativeModel(modelName)
	if model == nil {
		// If client creation succeeded but model is nil, close the client to avoid resource leaks
		client.Close()
		return nil, nil, errors.New("failed to initialize model: " + modelName)
	}

	// Configure model with system instructions
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(SystemInstructions),
		},
	}

	return client, model, nil
}

