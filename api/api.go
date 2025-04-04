package api

import (
	"errors"
	"os"
)

// GetAPIKey retrieves the Gemini API key from environment variables.
// Returns an error if the API key is not set.
func GetAPIKey() (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", errors.New("GEMINI_API_KEY environment variable is required")
	}
	return apiKey, nil
}