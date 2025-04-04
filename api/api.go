package api

import (
	"context"
	"errors"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Default model to use for Gemini API
const DefaultModelName = "gemini-2.5-pro-exp-03-25"

// SystemInstructions defines the system instructions for the resume generation model
const SystemInstructions = `You are an expert resume writing assistant. Your goal is to synthesize the provided existing resume information (if any) and the raw stream-of-consciousness input into a single, coherent, professional resume formatted strictly in Markdown.

Prioritize clarity, conciseness, and professional language. Structure the output logically with clear headings (e.g., Summary, Experience, Projects, Skills, Education). Infer structure and dates where possible, but do not fabricate information not present in the inputs. Focus on elevating the user's actual experience. Ensure the final output is only Markdown content.`

// GetAPIKey retrieves the Gemini API key from environment variables.
// Returns an error if the API key is not set.
func GetAPIKey() (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", errors.New("GEMINI_API_KEY environment variable is required")
	}
	return apiKey, nil
}

// InitializeClient creates a new Gemini client with the provided API key 
// and initializes it with the default model.
// Returns the client, model, and any error that occurred.
func InitializeClient(ctx context.Context, apiKey string) (*genai.Client, *genai.GenerativeModel, error) {
	return InitializeClientWithModel(ctx, apiKey, DefaultModelName)
}

// InitializeClientWithModel creates a new Gemini client with the provided API key
// and initializes it with the specified model name.
// Returns the client, model, and any error that occurred.
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

