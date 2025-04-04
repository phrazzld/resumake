package api

import (
	"context"
	"os"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	// Save original environment to restore after tests
	originalAPIKey := os.Getenv("GEMINI_API_KEY")
	defer os.Setenv("GEMINI_API_KEY", originalAPIKey)

	// Test case 1: API key is set correctly
	t.Run("API key is set", func(t *testing.T) {
		expected := "test-api-key-123"
		os.Setenv("GEMINI_API_KEY", expected)
		
		actual, err := GetAPIKey()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if actual != expected {
			t.Errorf("Expected API key %q, got %q", expected, actual)
		}
	})

	// Test case 2: API key is missing
	t.Run("API key is missing", func(t *testing.T) {
		os.Unsetenv("GEMINI_API_KEY")
		
		_, err := GetAPIKey()
		if err == nil {
			t.Error("Expected error when API key is missing, got nil")
		}
	})
	
	// Test case 3: API key is empty
	t.Run("API key is empty", func(t *testing.T) {
		os.Setenv("GEMINI_API_KEY", "")
		
		_, err := GetAPIKey()
		if err == nil {
			t.Error("Expected error when API key is empty, got nil")
		}
	})
}

func TestInitializeClient(t *testing.T) {
	// Save original environment to restore after tests
	originalAPIKey := os.Getenv("GEMINI_API_KEY")
	defer os.Setenv("GEMINI_API_KEY", originalAPIKey)

	// Test case 1: Successfully initialize client with valid API key
	t.Run("Successfully initialize client", func(t *testing.T) {
		// Set a dummy API key (won't make real API calls in test)
		apiKey := "test-api-key-123"
		
		// Initialize client with the API key
		ctx := context.Background()
		client, model, err := InitializeClient(ctx, apiKey)
		
		// Assert client is not nil and no error is returned
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if client == nil {
			t.Error("Expected client to be initialized, got nil")
		}
		
		// Check if model is initialized with default model
		if model == nil {
			t.Error("Expected model to be initialized, got nil")
		}
	})

	// Test case 2: Fail to initialize client with invalid API key
	t.Run("Fail to initialize client with empty API key", func(t *testing.T) {
		// Set an invalid API key
		apiKey := ""
		
		// Attempt to initialize client
		ctx := context.Background()
		_, _, err := InitializeClient(ctx, apiKey)
		
		// Assert an error is returned
		if err == nil {
			t.Error("Expected error when API key is empty, got nil")
		}
	})

	// Test case 3: Initialize client with specific model
	t.Run("Initialize client with specific model", func(t *testing.T) {
		// Set a dummy API key
		apiKey := "test-api-key-123"
		
		// Initialize client with a specific model name
		customModel := "custom-model-name"
		ctx := context.Background()
		client, model, err := InitializeClientWithModel(ctx, apiKey, customModel)
		
		// Assert client and model are properly initialized
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if client == nil {
			t.Error("Expected client to be initialized, got nil")
		}
		
		if model == nil {
			t.Error("Expected model to be initialized, got nil")
		}
	})
}